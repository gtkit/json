//go:build jsoniter

package json

import (
	"log/slog"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	Marshal = json.Marshal

	Unmarshal = json.Unmarshal

	MarshalIndent = json.MarshalIndent

	NewDecoder = json.NewDecoder

	NewEncoder = json.NewEncoder
)

func CheckJSON() {
	slog.Info("jsoniter is used for JSON")
}

func SupportPrivateFields() {
	// Enable support for private fields
	extra.SupportPrivateFields()
}
