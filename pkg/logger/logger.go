package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

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
		panic("cannot initialize sugar logger, error : " + err.Error())
	}

	zapLog = logger
}

func InfoErr(message string, err error) {
	zapLog.Sugar().Infow(message, zap.Error(err))
}

func DebugErr(message string, err error) {
	zapLog.Sugar().Debugw(message, zap.Error(err))
}

func ErrorErr(message string, err error) {
	zapLog.Sugar().Errorw(message, zap.Error(err))
}

func FatalErr(message string, err error) {
	zapLog.Sugar().Fatalw(message, zap.Error(err))
}

func Sync() {
	_ = zapLog.Sync()
}
