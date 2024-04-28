package json

import (
	"log/slog"
	"time"
)

func TrackTime(pre time.Time) time.Duration {
	elapsed := time.Since(pre)
	slog.Info("elapsed time: ", elapsed)
	return elapsed
}
