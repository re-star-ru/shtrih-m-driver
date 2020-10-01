package logger

type Logger interface {
	Debug(args ...interface{})
	Fatal(args ...interface{})
}
