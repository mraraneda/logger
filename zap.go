package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapConsoleConfiguration struct {
	Enable     bool
	JSONFormat bool
	Level      string
}

type ZapFileConfiguration struct {
	Enable     bool
	JSONFormat bool
	Level      string
	Path       string
	MaxSize    int // MB
	Compress   bool
	MaxAge     int // Days
	MaxBackups int // Maximum number of files
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func newZapLogger(config Configuration) (Logger, error) {
	cores := []zapcore.Core{}
	consoleConfig, fileConfig := getZapConfig(config)
	if consoleConfig.Enable {
		level := getZapLevel(consoleConfig.Level)
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(consoleConfig.JSONFormat), writer, level)
		cores = append(cores, core)
	}
	if fileConfig.Enable {
		level := getZapLevel(fileConfig.Level)
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileConfig.Path,
			MaxSize:    fileConfig.MaxSize,
			Compress:   fileConfig.Compress,
			MaxAge:     fileConfig.MaxAge,
			MaxBackups: fileConfig.MaxBackups,
		})
		core := zapcore.NewCore(getEncoder(fileConfig.JSONFormat), writer, level)
		cores = append(cores, core)
	}
	combinedCore := zapcore.NewTee(cores...)
	logger := zap.New(combinedCore,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	).Sugar()
	return &zapLogger{
		sugaredLogger: logger,
	}, nil
}

func getZapConfig(config Configuration) (ZapConsoleConfiguration, ZapFileConfiguration) {
	var consoleConfig ZapConsoleConfiguration
	var fileConfig ZapFileConfiguration
	if config, ok := config[ZapConsoleConfig]; ok {
		consoleConfig = config.(ZapConsoleConfiguration)
	}
	if config, ok := config[ZapFileConfig]; ok {
		fileConfig = config.(ZapFileConfiguration)
	}
	return consoleConfig, fileConfig
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) GetLogger() interface{} {
	return l.sugaredLogger
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}
