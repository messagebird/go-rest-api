package messagebirdtest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var server *httptest.Server

var responseBody []byte
var status int = 200

// EnableServer starts a fake server, runs the test and closes the server.
func EnableServer(m *testing.M) {
	initAndStartServer()
	exitCode := m.Run()
	closeServer()

	os.Exit(exitCode)
}

func initAndStartServer() {
	server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// status and responseBody are defined in returns.go.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if _, err := w.Write(responseBody); err != nil {
			panic(err.Error())
		}
	}))
}

func closeServer() {
	server.Close()
}

// WillReturnTestdata sets the status (s) for the test server to respond with.
// Additionally it reads the bytes from the relativePath file and returns that
// for requests. It fails the test if the file can not be read. The path is
// relative to the testdata directory (the go tool ignores directories named
// "testdata" in test packages: https://golang.org/cmd/go/#hdr-Test_packages).
func WillReturnTestdata(t *testing.T, relativePath string, s int) {
	responseBody = testdata(t, relativePath)
	status = s
}

// WillReturnAccessKeyError sets the response body and status for requests to
// indicate the request is not allowed due to an incorrect access key.
func WillReturnAccessKeyError() {
	responseBody = []byte(`
		{
			"errors": [
				{
					"code":2,
					"description":"Request not allowed (incorrect access_key)",
					"parameter":"access_key"
				}
			]
		}`)
	status = 401
}
