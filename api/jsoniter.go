//go:build jsoniter

package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

// Package indicates what library is being used for JSON encoding.
const Package = "github.com/json-iterator/go"

func init() {
	API = jsoniterApi{}
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type jsoniterApi struct{}

func (j jsoniterApi) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsoniterApi) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j jsoniterApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (j jsoniterApi) NewEncoder(writer io.Writer) Encoder {
	return json.NewEncoder(writer)
}

func (j jsoniterApi) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

func (j jsoniterApi) SupportPrivateFields() {
	// sonic does not support private fields
}

func (j jsoniterApi) RegisterFuzzyDecoders() {}
