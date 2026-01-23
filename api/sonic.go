//go:build sonic && (linux || windows || darwin)

package json

import (
	"io"

	"github.com/bytedance/sonic"
)

// Package indicates what library is being used for JSON encoding.
const Package = "github.com/bytedance/sonic"

func init() {
	API = sonicApi{}
}

var json = sonic.ConfigStd

type sonicApi struct{}

func (j sonicApi) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j sonicApi) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j sonicApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (j sonicApi) NewEncoder(writer io.Writer) Encoder {
	return json.NewEncoder(writer)
}

func (j sonicApi) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

func (j sonicApi) SupportPrivateFields() {
	// sonic does not support private fields
}

func (j sonicApi) RegisterFuzzyDecoders() {}
