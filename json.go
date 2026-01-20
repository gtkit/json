//go:build !jsoniter && !go_json && !(sonic && (linux || windows || darwin))

package json

import (
	"encoding/json"
	"io"
	"log"
)

// Package indicates what library is being used for JSON encoding.
const Package = "encoding/json"

func init() {
	API = jsonApi{}
}

type jsonApi struct{}

func (j jsonApi) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonApi) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j jsonApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (j jsonApi) NewEncoder(writer io.Writer) Encoder {
	return json.NewEncoder(writer)
}

func (j jsonApi) NewDecoder(reader io.Reader) Decoder {
	return json.NewDecoder(reader)
}

func CheckJSON() {
	log.Println("standard json package is used for JSON")
}

func SupportPrivateFields() {}
