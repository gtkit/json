//go:build go_json

package json

import (
	"io"
	"unsafe"

	gojson "github.com/goccy/go-json"
)

// Package indicates the JSON library in use.
const Package = "github.com/goccy/go-json"

func init() {
	API = goJSONAPI{}
}

type goJSONAPI struct{}

func (goJSONAPI) Marshal(v any) ([]byte, error) {
	return gojson.Marshal(v)
}

func (goJSONAPI) Unmarshal(data []byte, v any) error {
	return gojson.Unmarshal(data, v)
}

func (goJSONAPI) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return gojson.MarshalIndent(v, prefix, indent)
}

func (goJSONAPI) MarshalToString(v any) (string, error) {
	b, err := gojson.Marshal(v)
	if err != nil {
		return "", err
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func (goJSONAPI) NewEncoder(writer io.Writer) Encoder {
	return gojson.NewEncoder(writer)
}

func (goJSONAPI) NewDecoder(reader io.Reader) Decoder {
	return gojson.NewDecoder(reader)
}

func (goJSONAPI) Valid(data []byte) bool {
	return gojson.Valid(data)
}
