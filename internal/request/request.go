package request

import (
	"fmt"
	"httpserver/internal/headers"
	"httpserver/internal/requestline"
	"io"
	"strconv"
)

var SEPERATOR = []byte("\r\n")

const BUFF_SIZE = 8

type ParserState int

const (
	INIT            ParserState = 0
	PARSING_HEADERS ParserState = 1
	PARSING_BODY    ParserState = 2
	DONE            ParserState = 3
)

type Request struct {
	RequestLine requestline.RequestLine
	Headers     headers.Headers
	Body        []byte
	State       ParserState
}

func newRequest() *Request {
	return &Request{
		State:   INIT,
		Headers: headers.NewHeaders(),
		Body:    []byte(""),
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
			r.State = PARSING_BODY
			return n + len(SEPERATOR), nil
		}
		return n, nil
	case PARSING_BODY:
		contentLength, ok := r.Headers.Get("Content-Length")

		if !ok || contentLength == "0" {
			r.State = DONE
			return 0, nil
		}
		length, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, fmt.Errorf("invalid content length")
		}

		readBody := min(length-len(r.Body), len(data))
		r.Body = append(r.Body, data[:readBody]...)

		if len(r.Body) > length {
			return 0, fmt.Errorf("body length exceeds content length")
		}

		if len(r.Body) == length {
			r.State = DONE
		}

		return readBody, nil
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
			if request.State == PARSING_BODY {
				contentLength, ok := request.Headers.Get("Content-Length")
				if ok && contentLength != "0" {
					length, _ := strconv.Atoi(contentLength)
					if len(request.Body) < length {
						return nil, fmt.Errorf("body length shorter than content length")
					}
				}
			}
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
