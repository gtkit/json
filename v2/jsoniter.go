//go:build jsoniter

package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

// Package indicates the JSON library in use.
const Package = "github.com/json-iterator/go"

func init() {
	API = jsoniterAPI{}
}

var jsoniterJSON = jsoniter.ConfigCompatibleWithStandardLibrary

type jsoniterAPI struct{}

func (jsoniterAPI) Marshal(v any) ([]byte, error) {
	return jsoniterJSON.Marshal(v)
}

func (jsoniterAPI) Unmarshal(data []byte, v any) error {
	return jsoniterJSON.Unmarshal(data, v)
}

func (jsoniterAPI) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return jsoniterJSON.MarshalIndent(v, prefix, indent)
}

func (jsoniterAPI) MarshalToString(v any) (string, error) {
	return jsoniterJSON.MarshalToString(v)
}

func (jsoniterAPI) NewEncoder(writer io.Writer) Encoder {
	return jsoniterJSON.NewEncoder(writer)
}

func (jsoniterAPI) NewDecoder(reader io.Reader) Decoder {
	return jsoniterJSON.NewDecoder(reader)
}

func (jsoniterAPI) Valid(data []byte) bool {
	return jsoniterJSON.Valid(data)
}

// SupportPrivateFields enables encoding/decoding of unexported struct fields.
// This is a jsoniter-specific feature; other backends do not support this.
// Call this early in main() before any concurrent access.
func SupportPrivateFields() {
	extra.SupportPrivateFields()
}

// RegisterFuzzyDecoders enables PHP-compatible fuzzy decoding:
// string "123" can decode to int, "true" can decode to bool, etc.
// This is a jsoniter-specific feature.
// Call this early in main() before any concurrent access.
func RegisterFuzzyDecoders() {
	extra.RegisterFuzzyDecoders()
}
