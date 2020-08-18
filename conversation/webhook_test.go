package conversation

import (
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
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
	t.Run("limit_offset", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "webhookListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		webhookList, err := ListWebhooks(client, &ListOptions{Limit: 20, Offset: 2})
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

		if query := mbtest.Request.URL.RawQuery; query != "limit=20&offset=2" {
			t.Fatalf("got %s, expected limit=20&offset=2", query)
		}
	})

	t.Run("all", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "allWebhookListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		webhookList, err := ListWebhooks(client, nil)
		if err != nil {
			t.Fatalf("unexpected error listing Webhooks: %s", err)
		}

		if webhookList.Limit != 10 {
			t.Fatalf("got %d, expected 10", webhookList.Limit)
		}

		if webhookList.Offset != 0 {
			t.Fatalf("got %d, expected 0", webhookList.Offset)
		}

		if webhookList.TotalCount != 2 {
			t.Fatalf("got %d, expected 2", webhookList.TotalCount)
		}

		if webhookList.Items[0].Events[0] != WebhookEventMessageCreated {
			t.Fatalf("got %s expected message.created", webhookList.Items[0].Events[0])
		}

		if len(webhookList.Items) != 2 {
			t.Fatalf("got %d, expected 2", len(webhookList.Items))
		}

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/webhooks")

		if query := mbtest.Request.URL.RawQuery; query != "" {
			t.Fatalf("got %s, expected empty", query)
		}
	})
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

func TestUpdateWebhook(t *testing.T) {
	mbtest.WillReturnTestdata(t, "webhookUpdatedObject.json", http.StatusOK)
	client := mbtest.Client(t)

	webhookUpdateRequest := &WebhookUpdateRequest{
		Events: []WebhookEvent{
			WebhookEventConversationUpdated,
		},
		URL:    "https://example.com/mynewwebhookurl",
		Status: WebhookStatusDisabled,
	}

	webhook, err := UpdateWebhook(client, "whid", webhookUpdateRequest)
	if err != nil {
		t.Fatalf("unexpected error updating Webhook: %s", err)
	}

	if webhook.URL != "https://example.com/mynewwebhookurl" {
		t.Fatalf("Expected https://example.com/mynewwebhookurl, got %s", webhook.URL)
	}

	if webhook.UpdatedDatetime == nil {
		t.Fatalf("Expected the UpdatedDatetime value to be added, but was nil")
	}

	if webhook.Status != WebhookStatusDisabled {
		t.Fatalf("Expected status to be disabled, was %s", webhook.Status)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPatch, "/v1/webhooks/whid")
	mbtest.AssertTestdata(t, "webhookUpdateRequest.json", mbtest.Request.Body)
}
