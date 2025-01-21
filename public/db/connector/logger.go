package connector

import (
	"context"
	"errors"
	"financia/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

func init() {
	lumberjacklogger := &lumberjack.Logger{
		Filename:   "./runlog/log-rotate.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	defer lumberjacklogger.Close()

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.DateTime))
	} // 设置时间格式
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	var logOutput zapcore.WriteSyncer
	if config.Configs.Logger.Type == "file" {
		logOutput = zapcore.AddSync(lumberjacklogger)
	} else {
		logOutput = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(
		fileEncoder,    //编码设置
		logOutput,      //输出到文件
		zap.DebugLevel, //日志等级
	)

	log := zap.New(core, zap.AddCaller())
	defer log.Sync()
	zap.ReplaceGlobals(log)
}

type CustomLogger struct {
	logLevel                  logger.LogLevel
	ignoreRecordNotFoundError bool
	slowThreshold             time.Duration
}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *CustomLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel >= logger.Info {
		zap.S().Debugf(msg, args...)
	}
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel >= logger.Warn {
		zap.S().Infof(msg, args...)
	}
}

func (l *CustomLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel >= logger.Error {
		zap.S().Errorf(msg, args...)
	}
}

func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if l.ignoreRecordNotFoundError && errors.Is(err, gorm.ErrRecordNotFound) {
		// Ignore record not found errors
		return
	}

	switch {
	case err != nil && l.logLevel >= logger.Error:
		zap.S().Errorf("err: %v | elapsed: %v | sql: %s | rows: %d", err, elapsed, sql, rows)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		zap.S().Infof("elapsed: %v | slow sql: %s | rows: %d", elapsed, sql, rows)
	case l.logLevel >= logger.Info:
		zap.S().Debugf("elapsed: %v | sql: %s | rows: %d", elapsed, sql, rows)
	}
}
