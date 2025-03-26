package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds the configuration for the logger.
type Config struct {
	Level       string   // Debug, Info, Warn, Error, Fatal
	Encoding    string   // json or console
	OutputPaths []string // list of URLs or file paths
}

// DefaultConfig returns the default logger configuration.
func DefaultConfig() Config {
	return Config{
		Level:       "debug",
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
	}
}

var zapLog *zap.Logger

// Initialize sets up the logger with the given configuration.
func Initialize(cfg Config) error {
	var level zapcore.Level

	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}

	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Encoding:    cfg.Encoding,
		OutputPaths: cfg.OutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			TimeKey:        "time",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			CallerKey:      "source",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			StacktraceKey:  "stack-trace",
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}

	logger, err := logConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	zapLog = logger

	return nil
}

// init initializes the logger with default configuration.
func init() {
	if err := Initialize(DefaultConfig()); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// Info logs a message at info level.
func Info(msg string, fields ...zap.Field) {
	zapLog.Info(msg, fields...)
}

// Debug logs a message at debug level.
func Debug(msg string, fields ...zap.Field) {
	zapLog.Debug(msg, fields...)
}

// Warn logs a message at warn level.
func Warn(msg string, fields ...zap.Field) {
	zapLog.Warn(msg, fields...)
}

// Error logs a message at error level.
func Error(msg string, fields ...zap.Field) {
	zapLog.Error(msg, fields...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1).
func Fatal(msg string, fields ...zap.Field) {
	zapLog.Fatal(msg, fields...)
}

// With creates a child logger and adds structured context to it.
func With(fields ...zap.Field) *zap.Logger {
	return zapLog.With(fields...)
}

// Sync flushes any buffered log entries.
func Sync() error {
	return zapLog.Sync()
}

// GetLogger returns the underlying zap logger instance.
func GetLogger() *zap.Logger {
	return zapLog
}
