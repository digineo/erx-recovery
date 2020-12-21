package main

import (
	"bytes"
	"os"
)

type logWriter struct {
	buffer bytes.Buffer
	file   *os.File
}

func (w *logWriter) setFile(path string) error {
	file, err := os.Create(path)
	if err == nil {
		w.file = file
		w.buffer.WriteTo(file)
	}

	return err
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	if w.file == nil {
		return w.buffer.Write(p)
	}

	return w.file.Write(p)
}

func (w *logWriter) Close() error {
	if w.file == nil {
		return nil
	}

	return w.file.Close()
}
