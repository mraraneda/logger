package logger

import "errors"

// A global variable so that log functions can be directly accessed
var log Logger

// Fields type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Available logger level
const (
	// Debug has verbose message
	Debug = "debug"
	// Info is default log level
	Info = "info"
	// Warn is for logging messages about possible issues
	Warn = "warn"
	// Error is for logging errors
	Error = "error"
	// Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	Fatal = "fatal"
)

const (
	LogrusConsoleConfig = "LOGRUS_CONSOLE_CONFIG"
	LogrusFileConfig    = "LOGRUS_FILE_CONFIG"
	ZapConsoleConfig    = "ZAP_CONSOLE_CONFIG"
	ZapFileConfig       = "ZAP_FILE_CONFIG"
)

// Available logger instance
const (
	//InstanceZapLogger will be used to create Zap instance for the logger
	InstanceZapLogger int = iota
	//InstanceLogrusLogger will be used to create Logrus instance for the logger
	InstanceLogrusLogger
)

// Available error message
var (
	errInvalidLoggerInstance = errors.New("invalid logger instance")
)

// Logger is the list of all available method for logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
	GetLogger() interface{}
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration map[string]interface{}

// NewLogger returns an instance of logger provided
func NewLogger(config Configuration, loggerInstance int) error {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(config)
		if err != nil {
			return err
		}
		log = logger
		return nil

	case InstanceLogrusLogger:
		logger, err := newLogrusLogger(config)
		if err != nil {
			return err
		}
		log = logger
		return nil

	default:
		return errInvalidLoggerInstance
	}
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func GetLogger() interface{} {
	return log.GetLogger()
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}
