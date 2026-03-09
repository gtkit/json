//go:build sonic && (linux || windows || darwin)

package json

import (
	"log"

	"github.com/bytedance/sonic"
)

// 使用类型别名，避免与标准库 json 包的类型冲突.
type Encoder = sonic.Encoder
type Decoder = sonic.Decoder
type RawMessage = sonic.NoCopyRawMessage

var (
	json = sonic.ConfigStd

	Marshal = json.Marshal

	Unmarshal = json.Unmarshal

	MarshalToString = json.MarshalToString

	UnmarshalFromString = json.UnmarshalFromString

	MarshalIndent = json.MarshalIndent

	NewDecoder = json.NewDecoder

	NewEncoder = json.NewEncoder

	Valid = json.Valid
)

func CheckJSON() {
	log.Println("sonic is used for JSON")
}

func SupportPrivateFields() {
	// sonic does not support private fields
}

func RegisterFuzzyDecoders() {}

func SetFastest() {
	json = sonic.ConfigFastest
}
