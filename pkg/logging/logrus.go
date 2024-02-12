package logging

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var entry *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func New() *Logger {
	return &Logger{
		entry,
	}
}

type writeHook struct {
	w   []io.Writer
	lev []logrus.Level
}

func (hook *writeHook) Fire(e *logrus.Entry) error {
	line, err := e.String()
	if err != nil {
		return err
	}
	for _, ww := range hook.w {
		if _, err := ww.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}

func (hook *writeHook) Levels() []logrus.Level {
	return hook.lev
}

func init() {
	logger := logrus.New()
	// настройка строки вывода журналирования
	logger.ReportCaller = false
	logger.Formatter = &logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}

	// файл журналирования
	file, err := os.OpenFile("logger.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	logger.SetOutput(io.Discard)

	logrus.SetLevel(logrus.TraceLevel)

	logger.AddHook(&writeHook{
		w:   []io.Writer{file},
		lev: logrus.AllLevels,
	})

	entry = logrus.NewEntry(logger)
}
