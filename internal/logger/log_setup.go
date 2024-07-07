package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger
var once sync.Once

type LoggerFactory interface {
	SetupLogger()
	SyncLogger()
}

type loggerFactory struct{}

func NewLoggerFactory() LoggerFactory {
	return &loggerFactory{}
}

func (lf *loggerFactory) SetupLogger() {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		logger := zap.Must(config.Build())
		Log = logger.Sugar()
	})
}

func (lf *loggerFactory) SyncLogger() {
	if Log != nil {
		_ = Log.Sync()
	}
}
