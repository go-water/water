package water

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	opts := []zap.Option{zap.ErrorOutput(os.Stdout)}

	if viper.GetBool("zap.addCaller") {
		opts = []zap.Option{zap.AddCaller(), zap.AddCallerSkip(2)}
	}

	if viper.GetBool("zap.color") {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if viper.GetBool("zap.fullCaller") {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	if viper.GetString("zap.encoder") == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	newCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(zapcore.AddSync(os.Stdout)),
		zap.NewAtomicLevelAt(convertLogLevel(viper.GetString("zap.level"))))
	l := zap.New(newCore, opts...)
	return l
}

func convertLogLevel(levelStr string) (level zapcore.Level) {
	switch levelStr {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	return
}
