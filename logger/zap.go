package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	Level      string `yaml:"level"`
	Encoding   string `yaml:"encoding"`
	CallFull   bool   `yaml:"call_full"`
	Filename   string `yaml:"file_name"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	LocalTime  bool   `yaml:"local_time"`
	Compress   bool   `yaml:"compress"`
	Color      bool   `yaml:"color"`
	AddCaller  bool   `yaml:"add_caller"`
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
	}
	return
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

	if conf.Filename == "" {
		newCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(zapcore.AddSync(os.Stdout)),
			zap.NewAtomicLevelAt(convertLogLevel(conf.Level)))
		l := zap.New(newCore, opts...)
		return l
	}

	newCore := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.NewAtomicLevelAt(convertLogLevel(conf.Level)))
	l := zap.New(newCore, opts...)
	return l
}
