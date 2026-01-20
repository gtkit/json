//go:build jsoniter

package json

import (
	"io"
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
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

func CheckJSON() {
	log.Println("jsoniter is used for JSON")
}

func SupportPrivateFields() {
	// Enable support for private fields
	extra.SupportPrivateFields() // 支持非导出字段
}

func RegisterFuzzyDecoders() {
	extra.RegisterFuzzyDecoders() // 开启 php 兼容模式
}
