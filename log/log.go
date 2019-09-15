package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

var defaultLogger = New(os.Stderr)

var _ Logger = (*logrusLogger)(nil)

type Logger interface {
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	SetOutput(w io.Writer)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

func NewDebugHook() logrus.Hook {
	return debugHook{}
}

type debugHook struct{}

func (h debugHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h debugHook) Fire(entry *logrus.Entry) error {
	if pc, file, line, ok := runtime.Caller(8); ok {
		funcName := runtime.FuncForPC(pc).Name()

		entry.Data["source"] = fmt.Sprintf("%s:%v:%s", path.Base(file), line, path.Base(funcName))
	}

	return nil
}

func New(w io.Writer) Logger {
	return newLogrusLogger(w)
}

// NewStandardOutLogger returns a new logger
// using the os standard output.
func NewStandardOutLogger() Logger {
	return New(os.Stdout)
}

func newLogrusLogger(w io.Writer) *logrusLogger {
	l := &logrusLogger{
		l: logrus.New(),
	}
	l.l.SetOutput(w)
	return l
}

type logrusLogger struct {
	out io.Writer
	l   *logrus.Logger
}

func (l *logrusLogger) init() {
	if l.out == nil {
		l.out = os.Stdout
	}

	l.l = logrus.New()
	l.l.SetOutput(l.out)
}

func (l *logrusLogger) SetOutput(w io.Writer) {
	l.l.SetOutput(w)
}

func (l *logrusLogger) Infof(format string, v ...interface{}) {
	l.l.Printf(format, v...)
}

func (l *logrusLogger) Errorf(format string, v ...interface{}) {
	l.l.Errorf(format, v...)
}
