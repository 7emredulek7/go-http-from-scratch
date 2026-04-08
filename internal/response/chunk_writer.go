package response

import (
	"fmt"
	"httpserver/internal/headers"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {

	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
