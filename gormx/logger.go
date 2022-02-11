package gormx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/waytohome/lightning/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	level2ModeMap = map[string]logger.LogLevel{
		"debug": logger.Info,
		"info":  logger.Info,
		"warn":  logger.Warn,
		"error": logger.Error,
	}
)

type LoggerConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

type Logger struct {
	logger.Config

	logger *logx.Logger
}

func WrapLogger(l *logx.Logger, conf LoggerConfig) logger.Interface {
	logLevel, ok := level2ModeMap[l.GetLevel()]
	if !ok {
		logLevel = logger.Silent
	}
	wl := &Logger{
		Config: logger.Config{
			SlowThreshold:             conf.SlowThreshold,
			Colorful:                  false,
			IgnoreRecordNotFoundError: conf.IgnoreRecordNotFoundError,
			LogLevel:                  logLevel,
		},
		logger: l,
	}
	return wl
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	panic("specifying the log level is not allowed, because log level is follow logx.Logger")
}

func (l *Logger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Info {
		l.logger.Debug(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *Logger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Warn {
		l.logger.Warn(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *Logger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= logger.Error {
		l.logger.Error(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *Logger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.Config.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	var fields []logx.Field
	fields = append(fields, logx.Float64("cost/ms", float64(elapsed.Nanoseconds())/1e6))
	switch {
	case l.Config.LogLevel >= logger.Error && err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		fields = append(fields, logx.String("err", err.Error()), logx.String("sql", sql), logx.Int64("rows", rows), logx.String("sql_location", utils.FileWithLineNum()))
		logx.Error("gorm trace log, err happened when sql execute", fields...)
	case l.Config.LogLevel >= logger.Warn && elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		fields = append(fields, logx.Int64("threshold/ms", l.SlowThreshold.Milliseconds()), logx.String("sql", sql), logx.Int64("rows", rows), logx.String("sql_location", utils.FileWithLineNum()))
		logx.Warn("gorm trace log, slow sql detected", fields...)
	case l.Config.LogLevel == logger.Info:
		sql, rows := fc()
		fields = append(fields, logx.String("sql", sql), logx.Int64("rows", rows), logx.String("sql_location", utils.FileWithLineNum()))
		l.logger.Debug("gorm trace log", fields...)
	default:
	}
}
