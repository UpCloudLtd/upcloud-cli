package log

import (
	"github.com/gemalto/flume"
)

var (
	log flume.Logger
)

func init() {
	log = flume.New("upctl")
}

// SetDebugMode sets debug mode on to get more verbose logs
func SetDebugMode(v bool) error {
	if err := flume.Configure(flume.Config{
		DefaultLevel: flume.DebugLevel,
		Encoding:     "term-color",
	}); err != nil {
		return err
	}
	if !v {
		if err := flume.Configure(flume.Config{
			DefaultLevel: flume.InfoLevel,
			Encoding:     "term-color",
		}); err != nil {
			return err
		}
	}
	return nil
}

// Info Writes info level logs
func Info(nameSpace string, msg string, args ...interface{}) {
	log.With("namespace: ", nameSpace).Info(msg, args)
}

// Debug Writes debug level logs
func Debug(nameSpace string, msg string, args ...interface{}) {
	log.With("namespace: ", nameSpace).Debug(msg, args)
}
