package headers

import (
	"bytes"
	"errors"
	"maps"
	"strings"
)

type Headers map[string]string

var SEPERATOR = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Get(fieldName string) (fieldValue string, ok bool) {
	fieldValue, ok = h[strings.ToLower(fieldName)]
	return fieldValue, ok
}

func (h Headers) GetAll() map[string]string {
	headers := make(map[string]string)
	maps.Copy(headers, h)
	return headers
}

func (h Headers) Set(fieldName, fieldValue string) error {
	existingFieldValue, ok := h.Get(fieldName)
	if ok {
		h[strings.ToLower(fieldName)] = existingFieldValue + ", " + fieldValue
		return nil
	}

	for _, c := range fieldName {
		if !isValidFieldNameChar(c) {
			return errors.New("invalid header field name: contains illegal character")
		}
	}
	h[strings.ToLower(fieldName)] = fieldValue
	return nil
}

func isValidFieldNameChar(c rune) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		strings.ContainsRune("!#$%&'*+-.^_`|~", c)
}

func parseHeader(fieldLine []byte) (fieldName, fieldValue string, err error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return fieldName, fieldValue, errors.New("invalid header format")
	}

	fieldName = string(parts[0])
	fieldValue = string(bytes.TrimSpace(parts[1]))
	if bytes.HasSuffix([]byte(fieldName), []byte(" ")) || bytes.HasPrefix([]byte(fieldName), []byte(" ")) {
		return fieldName, fieldValue, errors.New("invalid header format: field name cannot have leading or trailing spaces")
	}
	return fieldName, fieldValue, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	for {
		idx := bytes.Index(data, SEPERATOR)
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			break
		}

		fieldName, fieldValue, err := parseHeader(data[:idx])
		if err != nil {
			return n, done, err
		}

		err = h.Set(fieldName, fieldValue)
		if err != nil {
			return n, done, err
		}
		n += idx + len(SEPERATOR)
		data = data[idx+len(SEPERATOR):]
	}

	return n, done, nil
}
