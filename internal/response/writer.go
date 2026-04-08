package response

import (
	"httpserver/internal/headers"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	b := []byte{}
	for key, value := range headers.GetAll() {
		b = append(b, []byte(key+": "+value+"\r\n")...)
	}

	b = append(b, []byte("\r\n")...)
	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}
