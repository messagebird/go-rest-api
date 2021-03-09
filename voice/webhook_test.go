package voice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWebhook(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	const url = "https://example.com/voice-webhook"
	newWh, err := CreateWebHook(mbClient, url, "token")
	assert.NoError(t, err)
	assert.Equal(t, url, newWh.URL)
}

func TestWebhookList(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	for _, c := range []string{"foo", "bar", "baz"} {
		_, err := CreateWebHook(mbClient, "https://example/com/"+c, "token")
		assert.NoError(t, err)
	}

	i := 0
	for cf := range Webhooks(mbClient).Stream() {
		err, ok := cf.(error)
		assert.NoError(t, err)
		assert.True(t, ok)
		i++
	}
	assert.Equal(t, 0, i)
}
