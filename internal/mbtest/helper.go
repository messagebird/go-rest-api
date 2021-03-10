package mbtest

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testdataDir = "testdata"

// Testdata returns a file's bytes based on the path relative to the testdata
// directory. It fails the test if the testdata file can not be read.
func Testdata(t *testing.T, relativePath string) []byte {
	path := filepath.Join(testdataDir, relativePath)

	b, err := ioutil.ReadFile(path)
	assert.NoError(t, err)

	return b
}

// AssertTestdata gets testdata and asserts it equals actual. We start by
// slicing off all leading and trailing white space, as defined by Unicode.
func AssertTestdata(t *testing.T, relativePath string, actual []byte) {
	expected := bytes.TrimSpace(Testdata(t, relativePath))
	actual = bytes.TrimSpace(actual)

	assert.Truef(t, bytes.Equal(expected, actual), "expected %s, got %s", expected, actual)
}

// AssertEndpointCalled fails the test if the last request was not made to the
// provided endpoint (e.g. combination of HTTP method and path).
func AssertEndpointCalled(t *testing.T, method, path string) {
	assert.Equal(t, method, Request.Method)

	escapedPath := Request.URL.EscapedPath()
	assert.Equal(t, path, escapedPath)
}
