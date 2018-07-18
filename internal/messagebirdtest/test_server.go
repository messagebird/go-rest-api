package messagebirdtest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var server *httptest.Server

var responseBody []byte
var status int

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

// WillReturn sets the body (r) and status (s) for the test server to respond with.
func WillReturn(b []byte, s int) {
	responseBody = b
	status = s
}

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
