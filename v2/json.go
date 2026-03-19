//go:build !jsoniter && !go_json && !(sonic && (linux || windows || darwin))

package json

// Package indicates the JSON library in use.
const Package = "encoding/json"

func init() {
	// stdFallback is defined in api.go, reuse it as the default backend.
	API = stdFallback{}
}
