package messagebirdtest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type resetFunc func()

type request struct {
	Body   []byte
	Method string
	URL    *url.URL
}

// Request contains the lastly received http.Request by the fake server.
var Request request

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
		Request = request{
			Method: r.Method,
			URL:    r.URL,
		}

		var err error

		Request.Body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}

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

// setRequest sets some basic fields from the http.Request in our global Request
// struct.
func setRequest(r *http.Request) error {
	Request := request{
		Method: r.Method,
		URL:    r.URL,
	}

	var err error

	// Reading from the request body is fine, as it's not used elsewhere.
	// Server always returns fake data/testdata.
	Request.Body, err = ioutil.ReadAll(r.Body)

	return err
}

func WillReturn(b []byte, s int) {
	responseBody = b
	status = s
}

// WillReturnTestdata sets the status (s) for the test server to respond with.
// Additionally it reads the bytes from the relativePath file and returns that
// for requests. It fails the test if the file can not be read. The path is
// relative to the testdata directory (the go tool ignores directories named
// "testdata" in test packages: https://golang.org/cmd/go/#hdr-Test_packages).
func WillReturnTestdata(t *testing.T, relativePath string, s int) {
	WillReturn(Testdata(t, relativePath), s)
}

// WillReturnAccessKeyError sets the response body and status for requests to
// indicate the request is not allowed due to an incorrect access key.
func WillReturnAccessKeyError() {
	WillReturn([]byte(`
		{
			"errors": [
				{
					"code":2,
					"description":"Request not allowed (incorrect access_key)",
					"parameter":"access_key"
				}
			]
		}
	`), 401)
}
