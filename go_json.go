//go:build go_json

package json

import (
	"log"

	"github.com/goccy/go-json"
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
	log.Println("go-json is used for JSON")
}

func SupportPrivateFields() {
	// go-json does not support private fields
}
