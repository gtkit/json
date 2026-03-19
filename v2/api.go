package json

import (
	stdjson "encoding/json"
	"io"
	"unsafe"
)

// API is the active JSON codec, set via init() based on build tags.
// If no backend init() runs, the package-level init below guarantees
// API is never nil by falling back to encoding/json.
//
// Callers should prefer the top-level functions (Marshal, Unmarshal, etc.)
// which delegate to API internally.
var API Core

func init() {
	// This init runs after all backend-specific init() functions (file-level init
	// order is undefined across files, but all init()s complete before main).
	// If a backend already set API, this is a no-op.
	// If API is still nil (unexpected build tag combination, or the default backend
	// file was excluded), fall back to encoding/json to prevent nil-panic.
	if API == nil {
		API = stdFallback{}
	}
}

// Core defines the full capability set of a JSON codec.
type Core interface {
	// Marshal returns the JSON encoding of v.
	Marshal(v any) ([]byte, error)

	// Unmarshal parses the JSON-encoded data and stores the result in v.
	Unmarshal(data []byte, v any) error

	// MarshalIndent is like Marshal but applies Indent to format the output.
	MarshalIndent(v any, prefix, indent string) ([]byte, error)

	// MarshalToString returns the JSON encoding of v as a string.
	MarshalToString(v any) (string, error)

	// NewEncoder returns a new Encoder that writes to w.
	NewEncoder(writer io.Writer) Encoder

	// NewDecoder returns a new Decoder that reads from r.
	NewDecoder(reader io.Reader) Decoder

	// Valid reports whether data is a valid JSON encoding.
	Valid(data []byte) bool
}

// Encoder writes JSON values to an output stream.
type Encoder interface {
	// SetEscapeHTML specifies whether problematic HTML characters
	// should be escaped inside JSON quoted strings.
	SetEscapeHTML(on bool)

	// Encode writes the JSON encoding of v to the stream,
	// followed by a newline character.
	Encode(v any) error
}

// Decoder reads and decodes JSON values from an input stream.
type Decoder interface {
	// UseNumber causes the Decoder to unmarshal a number into an any as a
	// Number instead of as a float64.
	UseNumber()

	// DisallowUnknownFields causes the Decoder to return an error when the destination
	// is a struct and the input contains object keys which do not match any
	// non-ignored, exported fields in the destination.
	DisallowUnknownFields()

	// Decode reads the next JSON-encoded value from its
	// input and stores it in the value pointed to by v.
	Decode(v any) error
}

// stdFallback is the encoding/json implementation.
// Serves as both the default backend (json.go) and the safety-net fallback.
type stdFallback struct{}

func (stdFallback) Marshal(v any) ([]byte, error) {
	return stdjson.Marshal(v)
}

func (stdFallback) Unmarshal(data []byte, v any) error {
	return stdjson.Unmarshal(data, v)
}

func (stdFallback) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return stdjson.MarshalIndent(v, prefix, indent)
}

func (stdFallback) MarshalToString(v any) (string, error) {
	b, err := stdjson.Marshal(v)
	if err != nil {
		return "", err
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func (stdFallback) NewEncoder(w io.Writer) Encoder {
	return stdjson.NewEncoder(w)
}

func (stdFallback) NewDecoder(r io.Reader) Decoder {
	return stdjson.NewDecoder(r)
}

func (stdFallback) Valid(data []byte) bool {
	return stdjson.Valid(data)
}

// Top-level convenience functions that delegate to the active API.
// These provide the familiar json.Marshal / json.Unmarshal calling convention.

// Marshal returns the JSON encoding of v.
func Marshal(v any) ([]byte, error) {
	return API.Marshal(v)
}

// Unmarshal parses the JSON-encoded data and stores the result in v.
func Unmarshal(data []byte, v any) error {
	return API.Unmarshal(data, v)
}

// MarshalIndent is like Marshal but applies Indent to format the output.
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return API.MarshalIndent(v, prefix, indent)
}

// MarshalToString returns the JSON encoding of v as a string.
func MarshalToString(v any) (string, error) {
	return API.MarshalToString(v)
}

// NewEncoder returns a new Encoder that writes to w.
func NewEncoder(writer io.Writer) Encoder {
	return API.NewEncoder(writer)
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder(reader io.Reader) Decoder {
	return API.NewDecoder(reader)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
	return API.Valid(data)
}
