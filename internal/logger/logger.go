package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()

	// Настройка формата логов
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Создание директории для логов если её нет
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Log.Fatal("Не удалось создать директорию для логов:", err)
	}

	// Создание файла лога с текущей датой
	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Fatal("Не удалось открыть файл лога:", err)
	}

	// Настройка вывода логов в файл и консоль
	Log.SetOutput(file)
	Log.SetOutput(os.Stdout)

	// Установка уровня логирования
	Log.SetLevel(logrus.InfoLevel)
}

// Обертки для удобного логирования
func Info(args ...interface{}) {
	Log.Info(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
