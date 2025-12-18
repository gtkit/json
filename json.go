//go:build !jsoniter && !go_json && !(sonic && avx && (linux || windows || darwin) && amd64)

package json

import (
	"encoding/json"
	"log"
)

type Encoder = json.Encoder

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
