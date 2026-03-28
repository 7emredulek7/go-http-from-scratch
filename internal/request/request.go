package request

import (
	"fmt"
	"httpserver/internal/headers"
	"httpserver/internal/requestline"
	"io"
)

var SEPERATOR = []byte("\r\n")

const BUFF_SIZE = 8

type ParserState int

const (
	INIT            ParserState = 0
	PARSING_HEADERS ParserState = 1
	DONE            ParserState = 2
)

type Request struct {
	RequestLine requestline.RequestLine
	Headers     headers.Headers
	State       ParserState
}

func newRequest() *Request {
	return &Request{
		State:   INIT,
		Headers: headers.NewHeaders(),
	}
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case INIT:
		requestLine, n, err := requestline.Parse(data)
		if err != nil {
			return 0, err
		}

		if n == 0 { // need more data
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = PARSING_HEADERS
		return n + len(SEPERATOR), nil
	case PARSING_HEADERS:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = DONE
		}
		return n, nil
	case DONE:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFF_SIZE)
	readToIndex := 0
	request := newRequest()

	for request.State != DONE {
		if readToIndex >= len(buf) {
			growBuf := make([]byte, len(buf)*2)
			copy(growBuf, buf)
			buf = growBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			request.State = DONE
			break
		}

		readToIndex += n
		nParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[nParsed:])
		readToIndex -= nParsed
	}
	return request, nil
}
