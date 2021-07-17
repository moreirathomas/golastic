package logger_test

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/moreirathomas/golastic/pkg/logger"
)

// Tests

func TestFileWriter(t *testing.T) {
	const filename = "hello.txt"

	w := logger.DefaultFile(filename).Writer()
	testWriteString(t, w, "hello\n")
	testWriteString(t, w, "world\n")
	defer mustRemoveFile(filename)

	b, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("expected to read file %s, got error: %s", filename, err)
	}

	if exp, got := "hello\nworld\n", string(b); got != exp {
		t.Errorf("expected body `%s`, got `%s`", exp, got)
	}
}

func testWriteString(t *testing.T, w io.Writer, s string) {
	t.Helper()
	if _, err := w.Write([]byte(s)); err != nil {
		t.Errorf("expected to write %s, got error: %s", s, err)
	}
}

// Helpers

func mustRemoveFile(filename string) {
	if err := os.Remove(filename); err != nil {
		log.Panicf("failed to remove file %s", filename)
	}
}
