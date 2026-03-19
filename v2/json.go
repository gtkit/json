//go:build !jsoniter && !go_json && !(sonic && (linux || windows || darwin))

package json

import (
	"encoding/json"
	"io"
	"unsafe"
)

// Package indicates the JSON library in use.
const Package = "encoding/json"

func init() {
	API = stdAPI{}
}

type stdAPI struct{}

func (stdAPI) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (stdAPI) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (stdAPI) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (stdAPI) MarshalToString(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func (stdAPI) NewEncoder(writer io.Writer) Encoder {
	return json.NewEncoder(writer)
}

func (stdAPI) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

func (stdAPI) Valid(data []byte) bool {
	return json.Valid(data)
}
