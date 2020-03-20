package log

type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	ErrorLevel
	DisabledLevel
)
