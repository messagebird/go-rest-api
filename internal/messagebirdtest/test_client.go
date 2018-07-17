package messagebirdtest

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"testing"

	messagebird "github.com/messagebird/go-rest-api"
)

// testWriter can be used to have the client write to the tests's error log.
type testWriter struct {
	t *testing.T
}

// Write logs the provided buffer to the current test's error log.
func (w testWriter) Write(p []byte) (int, error) {
	w.t.Log(p)

	return len(p), nil
}

// Client initializes a new MessageBird client that uses the
func Client(t *testing.T) *messagebird.Client {
	return client(t, "test_gshuPaZoeEG6ovbc8M79w0QyM")
}

func client(t *testing.T, accessKey string) *messagebird.Client {
	transport := &http.Transport{
		DialTLS: func(network, _ string) (net.Conn, error) {
			addr := server.Listener.Addr().String()
			return tls.Dial(network, addr, &tls.Config{
				InsecureSkipVerify: true,
			})
		},
	}

	return &messagebird.Client{
		AccessKey: accessKey,
		HTTPClient: &http.Client{
			Transport: transport,
		},
		DebugLog: log.New(testWriter{t: t}, "", 0),
	}
}

func testLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t: t}, "", 0)
}
