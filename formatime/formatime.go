package formatime

import (
	"go.uber.org/zap/zapcore"
	"time"
)

var RFC3339MS = "2006-01-02T15:04:05.000Z07:00"

// out -> 2022-04-06T23:17:31.385+08:00
func RFC3339TimeEncoderKibana(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(RFC3339MS))
}
