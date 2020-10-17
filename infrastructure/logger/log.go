package logger

import (
	"log"
	"os"
)

func NewLog() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
}

func NewLogFake() *log.Logger {
	return log.New(fakeWriter{}, "", 0)
}

type fakeWriter struct{}

func (f fakeWriter) Write(_ []byte) (n int, err error) { return 0, nil }
