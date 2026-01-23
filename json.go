//go:build !jsoniter && !go_json && !(sonic && (linux || windows || darwin))

package json

import (
	"encoding/json"
	"log"
)

type Encoder = json.Encoder
type Decoder = json.Decoder

var (
	Marshal = json.Marshal

	Unmarshal = json.Unmarshal

	MarshalIndent = json.MarshalIndent

	NewDecoder = json.NewDecoder

	NewEncoder = json.NewEncoder
)

func CheckJSON() {
	log.Println("standard json package is used for JSON")
}

func SupportPrivateFields() {
	// standard json package does not support private fields
}

func RegisterFuzzyDecoders() {}
