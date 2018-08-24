package conversation

import (
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/internal/mbtest"
)

func TestCreateWebhook(t *testing.T) {
	mbtest.WillReturnTestdata(t, "webhookObject.json", http.StatusOK)
	client := mbtest.Client(t)

	webhook, err := CreateWebhook(client, &WebhookCreateRequest{
		ChannelID: "chid",
		Events: []WebhookEvent{
			WebhookEventConversationCreated,
			WebhookEventMessageUpdated,
		},
		URL: "https://example.com/webhooks",
	})
	if err != nil {
		t.Fatalf("unexpected error creating Webhook: %s", err)
	}

	if webhook.ID != "whid" {
		t.Fatalf("got %s, expected whid", webhook.ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/webhooks")
	mbtest.AssertTestdata(t, "webhookCreateRequest.json", mbtest.Request.Body)
}

func TestDeleteWebhook(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	if err := DeleteWebhook(client, "whid"); err != nil {
		t.Fatalf("unexpected error deleting Webhook: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/webhooks/whid")
}

func TestListWebhooks(t *testing.T) {
	mbtest.WillReturnTestdata(t, "webhookListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	webhookList, err := ListWebhooks(client, DefaultListOptions)
	if err != nil {
		t.Fatalf("unexpected error listing Webhooks: %s", err)
	}

	if webhookList.TotalCount != 1 {
		t.Fatalf("got %d, expected 1", webhookList.TotalCount)
	}

	if webhookList.Items[0].Events[0] != WebhookEventMessageCreated {
		t.Fatalf("got %s expected message.created", webhookList.Items[0].Events[0])
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/webhooks")
}

func TestReadWebhook(t *testing.T) {
	mbtest.WillReturnTestdata(t, "webhookObject.json", http.StatusOK)
	client := mbtest.Client(t)

	webhook, err := ReadWebhook(client, "whid")
	if err != nil {
		t.Fatalf("unexpected error reading Webhook: %s", err)
	}

	if webhook.ChannelID != "chid" {
		t.Fatalf("got %s, expected chid", webhook.ChannelID)
	}

	if count := len(webhook.Events); count != 2 {
		t.Fatalf("got %d events, expected 2", count)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/webhooks/whid")
}
