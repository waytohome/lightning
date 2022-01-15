package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type config struct {
	encoderConfig *zapcore.EncoderConfig
	logger        *lumberjack.Logger
	level         *zap.AtomicLevel

	dumpLogFile bool // 是否需要输出日志文件
}

type Option func(conf *config)

// WithLogLevelOption 指定日志输出级别
func WithLogLevelOption(level string) Option {
	return func(conf *config) {
		if lv, ok := levelMap[level]; ok {
			conf.level.SetLevel(lv)
		}
	}
}

func WithDumpLogFileOption(need bool) Option {
	return func(conf *config) {
		conf.dumpLogFile = need
	}
}

// WithLogOutputOption 指定日志文件名
func WithLogOutputOption(path string) Option {
	return func(conf *config) {
		conf.logger.Filename = path
	}
}

// WithLogMaxSizeOption 指定日志文件大小
func WithLogMaxSizeOption(size int) Option {
	return func(conf *config) {
		conf.logger.MaxSize = size
	}
}

// WithLogMaxBackupsOption 保存日志文件数量
func WithLogMaxBackupsOption(count int) Option {
	return func(conf *config) {
		conf.logger.MaxBackups = count
	}
}

// WithLogMaxAgeOption 日志文件保存天数
func WithLogMaxAgeOption(days int) Option {
	return func(conf *config) {
		conf.logger.MaxAge = days
	}
}
