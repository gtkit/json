//go:build jsonv2

package json

import (
	"bytes"
	stdjson "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"errors"
	"io"
)

// Package indicates the JSON library in use.
const Package = "encoding/json/v2"

func init() {
	API = jsonV2API{}
}

// jsonV2API is the encoding/json/v2 implementation.
//
// It keeps the native v2 defaults except for the three cases listed on
// marshalOpts and unmarshalOpts, so the observable behavior differs from the
// default backend: HTML characters are not escaped, duplicate object names are
// rejected, object names match case-sensitively, nil slices and maps encode as
// [] and {} instead of null, and invalid UTF-8 is an error.
type jsonV2API struct{}

// marshalOpts and unmarshalOpts hold the three places where the native v2
// default is overridden. Each one is a case where v2 alone would make this the
// only backend of the five that behaves differently, which would break the
// premise that a -tags switch needs no changes to calling code.
//
//   - Deterministic keeps map members sorted by key. Without it v2 emits them
//     in Go's randomized map iteration order, so the same value encodes to
//     different bytes from one call to the next and anything that signs, caches
//     or snapshots the output breaks intermittently.
//   - FormatDurationAsNano encodes a time.Duration as a nanosecond count.
//     Without it v2 reports "no default representation" and marshaling any
//     struct holding a Duration fails outright, in both directions.
//   - OmitEmptyWithLegacySemantics omits a field whose Go value is false, 0 or
//     a nil pointer. v2 alone only omits values that encode as null, "", {} or
//     [], so zero numbers and false would start appearing in output that has
//     always left them out.
//
// The three come from encoding/json rather than encoding/json/v2 because that
// is where the toolchain exposes them; api.go already depends on that package
// for Number and the Compact, Indent and HTMLEscape helpers.
var (
	marshalOpts = jsonv2.JoinOptions(
		jsonv2.Deterministic(true),
		stdjson.FormatDurationAsNano(true),
		stdjson.OmitEmptyWithLegacySemantics(true),
	)

	// FormatDurationAsNano affects decoding too; the other two are
	// marshal-only and are left out here.
	unmarshalOpts = stdjson.FormatDurationAsNano(true)
)

func (jsonV2API) Marshal(v any) ([]byte, error) {
	return jsonv2.Marshal(v, marshalOpts)
}

func (jsonV2API) Unmarshal(data []byte, v any) error {
	return jsonv2.Unmarshal(data, v, unmarshalOpts)
}

func (jsonV2API) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	b, err := jsonv2.Marshal(v, marshalOpts)
	if err != nil {
		return nil, err
	}
	return indentJSON(b, prefix, indent)
}

func (jsonV2API) MarshalToString(v any) (string, error) {
	b, err := jsonv2.Marshal(v, marshalOpts)
	if err != nil {
		return "", err
	}
	return marshalBytesToString(b), nil
}

func (jsonV2API) NewEncoder(writer io.Writer) Encoder {
	return &jsonV2Encoder{w: writer}
}

func (jsonV2API) NewDecoder(reader io.Reader) Decoder {
	return &jsonV2Decoder{
		dec:  jsontext.NewDecoder(reader),
		opts: []jsonv2.Options{unmarshalOpts},
	}
}

func (jsonV2API) Valid(data []byte) bool {
	return jsontext.Value(data).IsValid()
}

// indentJSON applies prefix and indent to already encoded JSON, mirroring how
// encoding/json builds MarshalIndent on top of Marshal.
//
// jsontext.WithIndent and jsontext.WithIndentPrefix are deliberately avoided:
// they panic on any argument containing something other than spaces and tabs,
// while every other backend accepts an arbitrary prefix.
//
// Callers decide whether to indent at all: MarshalIndent always indents, so
// MarshalIndent(v, "", "") produces multiline output with no indentation just
// as encoding/json does, whereas Encoder skips this entirely when neither
// prefix nor indent is set, which is what makes SetIndent("", "") disable
// indentation.
func indentJSON(b []byte, prefix, indent string) ([]byte, error) {
	var buf bytes.Buffer
	if err := stdjson.Indent(&buf, b, prefix, indent); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// unmarshalAnyAsNumber reproduces Decoder.UseNumber: a JSON number decoded
// into an any becomes a Number instead of a float64.
//
// encoding/json/v2 exposes no option for this, so every attempt to unmarshal
// into an any is intercepted. Elements of []any and map[string]any are also
// any, which is what makes the interception apply at every nesting level.
var unmarshalAnyAsNumber = jsonv2.WithUnmarshalers(
	jsonv2.UnmarshalFromFunc(func(dec *jsontext.Decoder, val *any) error {
		if dec.PeekKind() != jsontext.KindNumber {
			// Fall back to the default behavior for every other kind.
			return errors.ErrUnsupported
		}
		raw, err := dec.ReadValue()
		if err != nil {
			return err
		}
		*val = Number(raw)
		return nil
	}),
)

// jsonV2Encoder adapts encoding/json/v2 to the Encoder interface, whose
// SetEscapeHTML and SetIndent may be called after construction and take effect
// on every subsequent Encode.
//
// It encodes into a buffer and writes once, the same way encoding/json does,
// instead of driving a jsontext.Encoder whose options are fixed at construction.
type jsonV2Encoder struct {
	w            io.Writer
	escapeHTML   bool
	indentPrefix string
	indentValue  string
}

// SetEscapeHTML specifies whether problematic HTML characters
// should be escaped inside JSON quoted strings.
// This backend defaults to false, matching its own Marshal.
func (e *jsonV2Encoder) SetEscapeHTML(on bool) {
	e.escapeHTML = on
}

// SetIndent instructs the encoder to format each subsequent encoded value
// as if indented by the package-level function Indent.
// Calling SetIndent("", "") disables indentation.
func (e *jsonV2Encoder) SetIndent(prefix, indent string) {
	e.indentPrefix = prefix
	e.indentValue = indent
}

// Encode writes the JSON encoding of v to the stream,
// followed by a newline character.
func (e *jsonV2Encoder) Encode(v any) error {
	b, err := jsonv2.Marshal(v, marshalOpts, jsontext.EscapeForHTML(e.escapeHTML))
	if err != nil {
		return err
	}
	if e.indentPrefix != "" || e.indentValue != "" {
		if b, err = indentJSON(b, e.indentPrefix, e.indentValue); err != nil {
			return err
		}
	}
	_, err = e.w.Write(append(b, '\n'))
	return err
}

// jsonV2Decoder adapts encoding/json/v2 to the Decoder interface.
//
// opts starts with unmarshalOpts; UseNumber and DisallowUnknownFields append
// to it and every entry is passed to each UnmarshalDecode call, where they take
// precedence over the options the jsontext.Decoder was built with. That is what
// lets them be set after an earlier Decode and still apply to the ones that
// follow.
type jsonV2Decoder struct {
	dec           *jsontext.Decoder
	opts          []jsonv2.Options
	useNumber     bool
	rejectUnknown bool
}

// UseNumber causes the Decoder to unmarshal a number into an any as a
// Number instead of as a float64.
func (d *jsonV2Decoder) UseNumber() {
	if d.useNumber {
		return
	}
	d.useNumber = true
	d.opts = append(d.opts, unmarshalAnyAsNumber)
}

// DisallowUnknownFields causes the Decoder to return an error when the destination
// is a struct and the input contains object keys which do not match any
// non-ignored, exported fields in the destination.
func (d *jsonV2Decoder) DisallowUnknownFields() {
	if d.rejectUnknown {
		return
	}
	d.rejectUnknown = true
	d.opts = append(d.opts, jsonv2.RejectUnknownMembers(true))
}

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
func (d *jsonV2Decoder) Decode(v any) error {
	return jsonv2.UnmarshalDecode(d.dec, v, d.opts...)
}

// Buffered returns a reader of the data remaining in the Decoder's buffer.
func (d *jsonV2Decoder) Buffered() io.Reader {
	return bytes.NewReader(d.dec.UnreadBuffer())
}

// More reports whether there is another element in the current array or object
// being parsed.
func (d *jsonV2Decoder) More() bool {
	switch d.dec.PeekKind() {
	case jsontext.KindInvalid, jsontext.KindEndObject, jsontext.KindEndArray:
		return false
	default:
		return true
	}
}
