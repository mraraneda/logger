package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogrusConsoleConfiguration struct {
	Enable     bool
	JSONFormat bool
	Level      string
}

type LogrusFileConfiguration struct {
	Enable     bool
	JSONFormat bool
	Level      string
	Path       string
	MaxSize    int // MB
	Compress   bool
	MaxAge     int // Days
	MaxBackups int
}

type logrusLogEntry struct {
	entry *logrus.Entry
}

type logrusLogger struct {
	logger *logrus.Logger
}

func newLogrusLogger(config Configuration) (Logger, error) {
	consoleConfig, fileConfig := getLogrusConfig(config)
	level, err := getLogLevel(consoleConfig.Level, fileConfig.Level)
	if err != nil {
		return nil, err
	}
	lLogger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: getFormatter(consoleConfig.JSONFormat),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}
	log := &logrusLogger{
		logger: lLogger,
	}
	log.setOutput(consoleConfig.Enable, fileConfig)
	return log, nil
}

func getLogrusConfig(config Configuration) (LogrusConsoleConfiguration, LogrusFileConfiguration) {
	var consoleConfig LogrusConsoleConfiguration
	var fileConfig LogrusFileConfiguration
	if config, ok := config[LogrusConsoleConfig]; ok {
		consoleConfig = config.(LogrusConsoleConfiguration)
	}
	if config, ok := config[LogrusFileConfig]; ok {
		fileConfig = config.(LogrusFileConfiguration)
	}
	return consoleConfig, fileConfig
}

func getLogLevel(consoleLevel, filelevel string) (logrus.Level, error) {
	logLevel := consoleLevel
	if logLevel == "" {
		logLevel = filelevel
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return 0, err
	}
	return level, nil
}

func getFormatter(isJSON bool) logrus.Formatter {
	if isJSON {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}
}

func (l *logrusLogger) setOutput(enableConsole bool, fileConfig LogrusFileConfiguration) {

	fileHandler := &lumberjack.Logger{
		Filename:   fileConfig.Path,
		MaxSize:    fileConfig.MaxSize,
		Compress:   fileConfig.Compress,
		MaxAge:     fileConfig.MaxAge,
		MaxBackups: fileConfig.MaxBackups,
	}

	if enableConsole && fileConfig.Enable {
		l.logger.SetOutput(io.MultiWriter(l.logger.Out, fileHandler))
	} else {
		if fileConfig.Enable {
			l.logger.SetOutput(fileHandler)
			l.logger.SetFormatter(getFormatter(fileConfig.JSONFormat))
		}
	}

}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogger) GetLogger() interface{} {
	return l.logger
}
func (l *logrusLogEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusLogEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) Panicf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogEntry) GetLogger() interface{} {
	return l.entry
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}
