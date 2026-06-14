package handlers

// Logger интерфейс для логирования
type Logger interface {
	Printf(format string, v ...interface{})
}
