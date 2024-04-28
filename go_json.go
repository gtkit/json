//go:build go_json

package json

import (
	"log/slog"
)

var (
	Marshal = json.Marshal

	Unmarshal = json.Unmarshal

	MarshalIndent = json.MarshalIndent

	NewDecoder = json.NewDecoder

	NewEncoder = json.NewEncoder
)

func CheckJSON() {
	slog.Info("go-json is used for JSON")
}
