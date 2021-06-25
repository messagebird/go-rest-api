package voice

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	messagebird "github.com/messagebird/go-rest-api/v7"
)

func testRequest(status int, body []byte) (*messagebird.Client, func()) {
	mbServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(body)
	}))
	addr := mbServer.Listener.Addr().String()
	transport := &http.Transport{
		DialTLS: func(netw, _ string) (net.Conn, error) {
			return tls.Dial(netw, addr, &tls.Config{InsecureSkipVerify: true})
		},
	}
	mbClient := messagebird.New("test_gshuPaZoeEG6ovbc8M79w0QyM")
	mbClient.HTTPClient = &http.Client{Transport: transport}

	if testing.Verbose() {
		mbClient.DebugLog = log.New(os.Stdout, "DEBUG", log.Lshortfile)
	}
	return mbClient, func() { mbServer.Close() }
}

func testClient(t *testing.T) (*messagebird.Client, bool) {
	key, ok := os.LookupEnv("MB_TEST_KEY")
	if !ok {
		return nil, false
	}
	client := &messagebird.Client{
		AccessKey:  key,
		HTTPClient: &http.Client{},
		DebugLog:   log.New(testWriter{T: t}, "", 0),
	}
	return client, true
}

type testWriter struct {
	T *testing.T
}

func (tw testWriter) Write(b []byte) (int, error) {
	tw.T.Logf("%s", b)
	return len(b), nil
}
