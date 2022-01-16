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
	levelMap = map[logger.LogLevel]string{
		logger.Info:   "debug",
		logger.Warn:   "warn",
		logger.Error:  "error",
		logger.Silent: "silent",
	}
)

type Logger struct {
	logger.Config

	logger *logx.Logger
	level  logger.LogLevel
}

func WrapLogger(l *logx.Logger, conf logger.Config) logger.Interface {
	wl := &Logger{
		Config: conf,
		logger: l,
		level:  conf.LogLevel,
	}
	// 同步日志级别
	wl.LogMode(conf.LogLevel)
	return wl
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	if lv, ok := levelMap[level]; ok && lv != "silent" {
		l.logger.SetLevel(lv)
	}
	l.level = level
	return l
}

func (l *Logger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		l.logger.Info(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *Logger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		l.logger.Info(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *Logger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		l.logger.Info(fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *Logger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	var fields []logx.Field
	fields = append(fields, logx.Float64("cost/ms", float64(elapsed.Nanoseconds())/1e6))
	switch {
	case l.level >= logger.Error && err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		fields = append(fields, logx.String("err", err.Error()), logx.String("sql", sql), logx.Int64("rows", rows), logx.String("sql_location", utils.FileWithLineNum()))
		logx.Error("gorm trace log, err happened when sql execute", fields...)
	case l.level >= logger.Warn && elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		fields = append(fields, logx.Int64("threshold/ms", l.SlowThreshold.Milliseconds()), logx.String("sql", sql), logx.Int64("rows", rows), logx.String("sql_location", utils.FileWithLineNum()))
		logx.Warn("gorm trace log, slow sql detected", fields...)
	case l.level == logger.Info:
		sql, rows := fc()
		fields = append(fields, logx.String("sql", sql), logx.Int64("rows", rows), logx.String("sql_location", utils.FileWithLineNum()))
		l.logger.Info("gorm trace log", fields...)
	default:
	}
}
