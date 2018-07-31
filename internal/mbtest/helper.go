package mbtest

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const testdataDir = "testdata"

// Testdata returns a file's bytes based on the path relative to the testdata
// directory. It fails the test if the testdata file can not be read.
func Testdata(t *testing.T, relativePath string) []byte {
	path := filepath.Join(testdataDir, relativePath)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("%s", err)
	}

	return b
}

// AssertTestdata gets testdata and asserts it equals actual.
func AssertTestdata(t *testing.T, relativePath string, actual []byte) {
	expected := Testdata(t, relativePath)

	if !bytes.Equal(expected, actual) {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

// AssertEndpointCalled fails the test if the last request was not made to the
// provided endpoint (e.g. combination of HTTP method and path).
func AssertEndpointCalled(t *testing.T, method, path string) {
	if Request.Method != method {
		t.Fatalf("expected %s, got %s", method, Request.Method)
	}

	if escapedPath := Request.URL.EscapedPath(); escapedPath != path {
		t.Fatalf("expected %s, got %s", path, escapedPath)
	}
}
