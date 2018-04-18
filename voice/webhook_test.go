package voice

import (
	"testing"
)

func TestCreateWebhook(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	const url = "https://example.com/voice-webhook"
	newWh, err := CreateWebHook(mbClient, url, "token")
	if err != nil {
		t.Fatal(err)
	}
	if newWh.URL != url {
		t.Fatalf("Unexpected URL: %q", newWh.URL)
	}
}

func TestWebhookList(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	for _, c := range []string{"foo", "bar", "baz"} {
		if _, err := CreateWebHook(mbClient, "https://example/com/"+c, "token"); err != nil {
			t.Fatal(err)
		}
	}

	i := 0
	for cf := range Webhooks(mbClient).Stream() {
		if err, ok := cf.(error); ok {
			t.Fatal(err)
		}
		i++
	}
	if i == 0 {
		t.Fatal("no webhooks were fetched")
	}
}
