package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.JSONFormatter{})

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	Log.SetLevel(parsedLevel)
}

func withFields(kv ...any) *logrus.Entry {
	fields := logrus.Fields{}
	for i := 0; i < len(kv); i += 2 {
		if i+1 < len(kv) {
			key, ok := kv[i].(string)
			if ok {
				fields[key] = kv[i+1]
			}
		}
	}
	return Log.WithFields(fields)
}

func Fatal(msg string, kv ...any) {
	withFields(kv...).Fatal(msg)
}

func Error(msg string, kv ...any) {
	withFields(kv...).Error(msg)
}

func Info(msg string, kv ...any) {
	withFields(kv...).Info(msg)
}

func Warn(msg string, kv ...any) {
	withFields(kv...).Warn(msg)
}

func Debug(msg string, kv ...any) {
	withFields(kv...).Debug(msg)
}
