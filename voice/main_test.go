package voice

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/messagebird/go-rest-api"
)

func testClient(status int, body []byte) (*messagebird.Client, func()) {
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
	mbClient := &messagebird.Client{
		AccessKey:  "test_gshuPaZoeEG6ovbc8M79w0QyM",
		HTTPClient: &http.Client{Transport: transport},
	}
	if testing.Verbose() {
		mbClient.DebugLog = log.New(os.Stdout, "DEBUG", log.Lshortfile)
	}
	return mbClient, func() { mbServer.Close() }
}
