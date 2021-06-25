package mbtest

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"testing"

	messagebird "github.com/messagebird/go-rest-api/v7"
)

// testWriter can be used to have the client write to the tests's error log.
type testWriter struct {
	t *testing.T
}

// Write logs the provided buffer to the current test's error log.
func (w testWriter) Write(p []byte) (int, error) {
	w.t.Logf("%s", p)

	return len(p), nil
}

// Client initializes a new MessageBird client that uses the
func Client(t *testing.T) *messagebird.Client {
	return client(t, "")
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
	client := messagebird.New(accessKey)
	client.HTTPClient.Transport = transport
	client.DebugLog = testLogger(t)

	return client
}

// testLogger creates a new logger that writes to the test's output.
func testLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t: t}, "", 0)
}
