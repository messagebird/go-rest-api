package mbtest

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	messagebird "github.com/messagebird/go-rest-api/v9"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) EnableFeatures(feature messagebird.Feature) {
}
func (c *ClientMock) DisableFeatures(feature messagebird.Feature) {
}
func (c *ClientMock) IsFeatureEnabled(feature messagebird.Feature) bool {
	return false
}
func (c *ClientMock) Request(v interface{}, method, path string, data interface{}) error {
	return nil
}

func (c *ClientMock) RequestV2(v interface{}, method, path string, data interface{}) (*http.Response, error) {
	return nil, nil
}

// MockClient initializes a new mock of MessageBird client
func MockClient() messagebird.Client {
	return &ClientMock{}
}

// Client initializes a new MessageBird client that uses the
func Client(t *testing.T) *messagebird.DefaultClient {
	return newClient(t, "")
}

func newClient(t *testing.T, accessKey string) *messagebird.DefaultClient {
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
	client.DebugLog = newTestLogger(t)

	return client
}

// testWriter can be used to have the client write to the tests's error log.
type testWriter struct {
	t *testing.T
}

// Write logs the provided buffer to the current test's error log.
func (w testWriter) Write(p []byte) (int, error) {
	w.t.Logf("%s", p)

	return len(p), nil
}

// testLogger creates a new logger that writes to the test's output.
func newTestLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t: t}, "", 0)
}
