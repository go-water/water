package water

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	Level     zapcore.Level `yaml:"level"`
	Encoding  string        `yaml:"encoding"`
	CallFull  bool          `yaml:"call_full"`
	Color     bool          `yaml:"color"`
	AddCaller bool          `yaml:"add_caller"`
}

func (conf *Config) NewLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	opts := []zap.Option{zap.ErrorOutput(os.Stdout)}

	if conf.AddCaller {
		opts = []zap.Option{zap.AddCaller(), zap.AddCallerSkip(2)}
	}

	if conf.Color {
		// 彩色显示，纯控制台输出可以彩色输出
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if conf.CallFull {
		// 全路径地址，通常不需要
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	if conf.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	newCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(zapcore.AddSync(os.Stdout)),
		zap.NewAtomicLevelAt(conf.Level))
	l := zap.New(newCore, opts...)
	return l
}
