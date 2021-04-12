package log

import (
	"github.com/gemalto/flume"
)

// SetDebugMode sets debug mode on to get more verbose logs
func SetDebugMode(v bool) error {
	if err := flume.Configure(flume.Config{
		DefaultLevel: flume.DebugLevel,
		Encoding:     "ltsv",
	}); err != nil {
		return err
	}
	if !v {
		if err := flume.Configure(flume.Config{
			DefaultLevel: flume.InfoLevel,
			Encoding:     "ltsv",
		}); err != nil {
			return err
		}
	}
	return nil
}

// Info Writes info level logs
func Info(space string, msg string, args ...interface{}) {
	log := flume.New(space)
	log.Info(msg, args)
}

// Debug Writes debug level logs
func Debug(space string, msg string, args ...interface{}) {
	log := flume.New(space)
	log.Debug(msg, args)
}
