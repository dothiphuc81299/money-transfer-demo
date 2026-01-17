package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(name string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.EncoderConfig.LevelKey = "log_level"
	config.EncoderConfig.StacktraceKey = zapcore.OmitKey

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	log = log.Named(name)
	return log, nil
}
