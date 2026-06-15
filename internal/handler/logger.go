package handler

type Logger interface {
	Printf(format string, v ...any)
}
