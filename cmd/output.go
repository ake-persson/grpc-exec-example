package cmd

import (
	"bufio"
	"bytes"
	"io"
	"time"
)

type Reader struct {
	reader    io.Reader
	BytesRead int
}

func newReader(r io.Reader) *Reader {
	return &Reader{reader: r}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	r.BytesRead += n
	return n, err
}

type Message struct {
	Received time.Time
	Line     int
	Stderr   bool
	Message  string
}

type Writer struct {
	line   int
	stderr bool
	ch     chan Message
}

func NewWriter() (*Writer, *Writer) {
	ch := make(chan Message, 1024)
	return &Writer{line: 0, stderr: false, ch: ch}, &Writer{line: 0, stderr: true, ch: ch}
}

func (w *Writer) Chan() <-chan Message {
	return w.ch
}

func (w *Writer) Write(p []byte) (int, error) {
	reader := newReader(bytes.NewReader(p))
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		w.ch <- Message{
			Received: time.Now(),
			Line:     w.line,
			Stderr:   w.stderr,
			Message:  scanner.Text(),
		}
		w.line++
	}
	if scanner.Err() != io.EOF && scanner.Err() != nil {
		return reader.BytesRead, scanner.Err()
	}
	return reader.BytesRead, nil
}

func (w *Writer) Close() error {
	close(w.ch)
	return nil
}
