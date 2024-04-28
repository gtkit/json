//go:build sonic && avx && (linux || windows || darwin) && amd64

package json

import (
	"log/slog"

	"github.com/bytedance/sonic"
)

var (
	json = sonic.ConfigStd

	Marshal = json.Marshal

	Unmarshal = json.Unmarshal

	MarshalIndent = json.MarshalIndent

	NewDecoder = json.NewDecoder

	NewEncoder = json.NewEncoder
)

func CheckJSON() {
	slog.Info("sonic is used for JSON")
}
