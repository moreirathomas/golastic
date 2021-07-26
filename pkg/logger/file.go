package logger

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	defaultPerm = fs.ModePerm
)

type fileLogger struct {
	filename string
}

func (fl fileLogger) Write(b []byte) (int, error) {
	f, err := openFileAll(fl.filename, defaultFlag, defaultPerm)
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

// openFileAll, like os.OpenFile, opens a file creating it if necessary.
// Unlike os.OpenFile, it also creates the missing parent directories
// using os.MkdirAll.
func openFileAll(filename string, flag int, perm fs.FileMode) (*os.File, error) {
	f, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		// os.ErrNotExist is returned if a parent directory is missing
		if errors.Is(err, os.ErrNotExist) {
			// create missing directories
			if e := os.MkdirAll(filepath.Dir(filename), fs.ModePerm); e != nil {
				return nil, e
			}
			// reopen file
			return os.OpenFile(filename, flag, perm)
		}
		return nil, err
	}
	return f, nil
}
