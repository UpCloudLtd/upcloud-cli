package log

import (
	"github.com/gemalto/flume"
)

// SetDebugMode sets debug mode on to get more verbose logs
func SetDebugMode(v bool) error {
	if err := flume.Configure(flume.Config{
		AddCaller:    &v,
		DefaultLevel: flume.DebugLevel,
		Encoding:     "term-color",
	}); err != nil {
		return err
	}
	if !v {
		if err := flume.Configure(flume.Config{
			AddCaller:    &v,
			DefaultLevel: flume.InfoLevel,
			Encoding:     "term-color",
		}); err != nil {
			return err
		}
	}
	return nil
}

// Info Writes info level logs
func Info(log flume.Logger, msg string, args ...interface{}) {
	log.Info(msg, args)
}

// Debug Writes debug level logs
func Debug(log flume.Logger, msg string, args ...interface{}) {
	log.Debug(msg, args)
}
