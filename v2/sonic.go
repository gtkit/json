//go:build sonic && (linux || windows || darwin)

package json

import (
	"io"

	"github.com/bytedance/sonic"
)

// Package indicates the JSON library in use.
const Package = "github.com/bytedance/sonic"

func init() {
	API = sonicAPI{}
}

// sonicJSON holds the sonic config. Default is ConfigStd for maximum compatibility.
// Call SetFastest() to switch to ConfigFastest for higher throughput.
var sonicJSON sonic.API = sonic.ConfigStd

type sonicAPI struct{}

func (sonicAPI) Marshal(v any) ([]byte, error) {
	return sonicJSON.Marshal(v)
}

func (sonicAPI) Unmarshal(data []byte, v any) error {
	return sonicJSON.Unmarshal(data, v)
}

func (sonicAPI) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return sonicJSON.MarshalIndent(v, prefix, indent)
}

func (sonicAPI) MarshalToString(v any) (string, error) {
	return sonicJSON.MarshalToString(v)
}

func (sonicAPI) NewEncoder(writer io.Writer) Encoder {
	return sonicJSON.NewEncoder(writer)
}

func (sonicAPI) NewDecoder(reader io.Reader) Decoder {
	return sonicJSON.NewDecoder(reader)
}

func (sonicAPI) Valid(data []byte) bool {
	return sonicJSON.Valid(data)
}

// SetFastest switches sonic to ConfigFastest mode for maximum throughput.
// This disables some standard library compatibility features:
//   - sorting map keys is disabled
//   - HTML escaping is disabled
//   - key pretouch validation is disabled
//
// Call this early in main() before any concurrent access.
func SetFastest() {
	sonicJSON = sonic.ConfigFastest
}
