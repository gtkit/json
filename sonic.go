//go:build sonic && avx && (linux || windows || darwin) && amd64

package json

import (
	"log"

	"github.com/bytedance/sonic"
)

type Encoder = sonic.Encoder

var (
	json = sonic.ConfigStd

	Marshal = json.Marshal

	Unmarshal = json.Unmarshal

	MarshalIndent = json.MarshalIndent

	NewDecoder = json.NewDecoder

	NewEncoder = json.NewEncoder
)

func CheckJSON() {
	log.Println("sonic is used for JSON")
}

func SupportPrivateFields() {
	// sonic does not support private fields
}
