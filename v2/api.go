package json

import "io"

// API is the active JSON codec, set via init() based on build tags.
// Callers should prefer the top-level functions (Marshal, Unmarshal, etc.)
// which delegate to API internally.
var API Core

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
