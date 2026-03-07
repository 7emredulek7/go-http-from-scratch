package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

var SEPERATOR = []byte("\r\n")

const BUFF_SIZE = 8

type ParserState int

const (
	INIT ParserState = 0
	DONE ParserState = 1
)

type Request struct {
	RequestLine RequestLine
	State       ParserState
}

func newRequest() *Request {
	return &Request{
		State: INIT,
	}
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func newRequestLine(requestLineParts []string) (*RequestLine, error) {
	method := requestLineParts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := requestLineParts[1]

	httpParts := strings.Split(requestLineParts[2], "/")
	if len(httpParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", requestLineParts)
	}

	httpPart := httpParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := httpParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpParts[1],
	}, nil

}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	endOfRequestLineIdx := bytes.Index(b, SEPERATOR)
	if endOfRequestLineIdx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(b[:endOfRequestLineIdx])
	requestLineParts := strings.Split(requestLineText, " ")
	requestLine, err := newRequestLine(requestLineParts)

	if err != nil {
		return nil, 0, err
	}
	return requestLine, endOfRequestLineIdx, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case INIT:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 { // need more data
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = DONE
		return n, nil
	case DONE:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFF_SIZE, BUFF_SIZE)
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
