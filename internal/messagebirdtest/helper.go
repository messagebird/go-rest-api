package messagebirdtest

import (
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
