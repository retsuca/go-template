package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func init() {
	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05,000"))
			},
			CallerKey:     "source",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			StacktraceKey: "stack-trace",
		},
	}.Build()

	if err != nil {
		panic(fmt.Sprintf("cannot initialize sugar logger, error : %s", err.Error()))
	}

	zap.ReplaceGlobals(logger)
}
