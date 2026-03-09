//go:build go_json

package json

import (
	"log"

	"github.com/goccy/go-json"
)

type Encoder = json.Encoder
type Decoder = json.Decoder
type RawMessage = json.RawMessage

var (
	Marshal = json.Marshal

	Unmarshal = json.Unmarshal

	MarshalIndent = json.MarshalIndent

	NewDecoder = json.NewDecoder

	NewEncoder = json.NewEncoder
)

func CheckJSON() {
	log.Println("go-json is used for JSON")
}

func SupportPrivateFields() {
	// go-json does not support private fields
}

func RegisterFuzzyDecoders() {}

func SetFastest() {}
