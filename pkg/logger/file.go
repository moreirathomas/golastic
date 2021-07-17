package logger

import (
	"fmt"
	"log"
	"os"
)

type fileLogger struct {
	filename string
}

func (fl fileLogger) Write(b []byte) (int, error) {
	f, err := os.OpenFile(fl.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return 0, fmt.Errorf("write error: could not open file %s: %w", fl.filename, err)
	}
	return f.Write(b)
}

func File(filename, prefix string, flag int) *log.Logger {
	return log.New(fileLogger{filename: filename}, prefix, flag)
}

func DefaultFile(filename string) *log.Logger {
	return File(filename, log.Default().Prefix(), log.Default().Flags())
}
