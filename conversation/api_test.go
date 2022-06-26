package conversation

import (
	messagebird "github.com/messagebird/go-rest-api/v7"
	"github.com/messagebird/go-rest-api/v7/internal/mbtest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestRequestSandboxEnabled(t *testing.T) {
	data := struct{}{}
	method := http.MethodGet
	reqPath := "qwerty"

	client := mbtest.MockClient().(*mbtest.ClientMock)
	client.On("IsFeatureEnabled", messagebird.FeatureConversationsAPIWhatsAppSandbox).Return(true)
	client.On("Request", data, method, whatsappSandboxAPIRoot+"/"+reqPath, data).Return(nil)

	err := request(client, data, method, reqPath, data)

	assert.NoError(t, err)
}

func TestRequestSandboxDisabled(t *testing.T) {
	data := struct{}{}
	method := http.MethodGet
	reqPath := "qwerty"

	client := mbtest.MockClient().(*mbtest.ClientMock)
	client.On("IsFeatureEnabled", messagebird.FeatureConversationsAPIWhatsAppSandbox).Return(false)
	client.On("Request", data, method, apiRoot+"/"+reqPath, data).Return(nil)

	err := request(client, data, method, reqPath, data)

	assert.NoError(t, err)
}
