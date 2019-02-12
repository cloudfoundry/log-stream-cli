package command

import (
	"io"
	"sync"
)

type DedupeWriter struct {
	w         io.Writer
	mu        sync.Mutex
	lastWrite []byte
}

func (s *DedupeWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if string(p) == string(s.lastWrite) {
		return 0, nil
	}

	s.lastWrite = p
	return s.w.Write(p)
}

func NewDedupeWriter(writer io.Writer) *DedupeWriter {
	return &DedupeWriter{
		w: writer,
	}
}
