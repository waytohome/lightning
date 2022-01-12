package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

var (
	LevelDebug = zap.DebugLevel
	LevelInfo  = zap.InfoLevel
	LevelWarn  = zap.WarnLevel
	LevelError = zap.ErrorLevel
	LevelFatal = zap.FatalLevel
	LevelPanic = zap.PanicLevel

	levelMap = map[string]Level{
		"trace": LevelDebug,
		"debug": LevelDebug,
		"info":  LevelInfo,
		"warn":  LevelWarn,
		"error": LevelError,
		"fatal": LevelFatal,
		"panic": LevelPanic,
	}
)
