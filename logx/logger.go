package logx

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultLogOutput     = "./logs/log.log"
	DefaultLogMaxSize    = 5
	DefaultLogMaxBackups = 30
	DefaultLogMaxAge     = 30
	DefaultLogCompress   = true
)

var (
	Def      = NewLogger()
	Debug    = Def.Debug
	Info     = Def.Info
	Warn     = Def.Warn
	Error    = Def.Error
	SetLevel = Def.SetLevel
	GetLevel = Def.GetLevel
)

type Logger struct {
	*zap.Logger
	level *zap.AtomicLevel
}

func NewLogger(opts ...Option) *Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.CallerKey = "pos"
	encoderConfig.StacktraceKey = "stack"

	lumberjackLogger := lumberjack.Logger{
		Filename:   DefaultLogOutput,
		MaxSize:    DefaultLogMaxSize,
		MaxBackups: DefaultLogMaxBackups,
		MaxAge:     DefaultLogMaxAge,
		Compress:   DefaultLogCompress,
	}

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.WarnLevel)

	conf := config{
		encoderConfig: &encoderConfig,
		logger:        &lumberjackLogger,
		level:         &level,
		dumpLogFile:   false,
	}
	for _, opt := range opts {
		opt(&conf)
	}

	encoder := zapcore.NewJSONEncoder(*conf.encoderConfig)
	var writers []zapcore.WriteSyncer
	writers = append(writers, zapcore.AddSync(os.Stdout))
	if conf.dumpLogFile {
		writers = append(writers, zapcore.AddSync(conf.logger))
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(writers...)
	core := zapcore.NewCore(encoder, writeSyncer, conf.level)

	var zopts []zap.Option
	zopts = append(zopts, zap.AddCaller())
	zopts = append(zopts, zap.ErrorOutput(writeSyncer))
	zopts = append(zopts, zap.AddStacktrace(zap.ErrorLevel))
	l := zap.New(core, zopts...)
	return &Logger{l, &level}
}

func (l *Logger) SetLevel(level string) {
	if lv, ok := levelMap[level]; ok {
		l.level.SetLevel(lv)
	}
}

func (l *Logger) GetLevel() Level {
	return l.level.Level()
}
