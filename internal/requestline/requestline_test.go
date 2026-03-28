package requestline

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	requestLine, n, err := Parse([]byte("GET /coffee HTTP/1.1\r\nHost: example.com\r\n"))
	require.NoError(t, err)
	require.NotNil(t, requestLine)
	assert.Equal(t, 20, n)
	assert.Equal(t, "GET", requestLine.Method)
	assert.Equal(t, "/coffee", requestLine.RequestTarget)
	assert.Equal(t, "1.1", requestLine.HttpVersion)
}

func TestParseIncomplete(t *testing.T) {
	requestLine, n, err := Parse([]byte("GET /coffee HTTP/1.1"))
	require.NoError(t, err)
	assert.Nil(t, requestLine)
	assert.Equal(t, 0, n)
}

func TestParseInvalidMethod(t *testing.T) {
	requestLine, n, err := Parse([]byte("Get /coffee HTTP/1.1\r\n"))
	require.Error(t, err)
	assert.Nil(t, requestLine)
	assert.Equal(t, 0, n)
}
