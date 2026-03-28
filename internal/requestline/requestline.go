package requestline

import (
	"bytes"
	"fmt"
	"strings"
)

var SEPERATOR = []byte("\r\n")

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func NewRequestLine(requestLineParts []string) (*RequestLine, error) {
	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("malformed start-line: %v", requestLineParts)
	}

	method := requestLineParts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := requestLineParts[1]

	httpParts := strings.Split(requestLineParts[2], "/")
	if len(httpParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %v", requestLineParts)
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
		HttpVersion:   version,
	}, nil
}

func Parse(data []byte) (*RequestLine, int, error) {
	endOfRequestLineIdx := bytes.Index(data, SEPERATOR)
	if endOfRequestLineIdx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:endOfRequestLineIdx])
	requestLineParts := strings.Split(requestLineText, " ")
	requestLine, err := NewRequestLine(requestLineParts)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, endOfRequestLineIdx, nil
}
