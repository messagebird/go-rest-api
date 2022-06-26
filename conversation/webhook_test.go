package conversation

import (
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/v7/internal/mbtest"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, "whid", webhook.ID)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/webhooks")
	mbtest.AssertTestdata(t, "webhookCreateRequest.json", mbtest.Request.Body)
}

func TestDeleteWebhook(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	err := DeleteWebhook(client, "whid")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/webhooks/whid")
}

func TestListWebhooks(t *testing.T) {
	t.Run("limit_offset", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "webhookListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		webhookList, err := ListWebhooks(client, &PaginationRequest{Limit: 20, Offset: 2})
		assert.NoError(t, err)

		assert.Equal(t, 1, webhookList.TotalCount)

		assert.Equal(t, WebhookEventMessageCreated, webhookList.Items[0].Events[0])

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/webhooks")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "limit=20&offset=2", query)
	})

	t.Run("all", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "allWebhookListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		webhookList, err := ListWebhooks(client, nil)
		assert.NoError(t, err)

		assert.Equal(t, 10, webhookList.Limit)

		assert.Equal(t, 0, webhookList.Offset)

		assert.Equal(t, 2, webhookList.TotalCount)

		assert.Equal(t, WebhookEventMessageCreated, webhookList.Items[0].Events[0])

		assert.Len(t, webhookList.Items, 2)

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/webhooks")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "", query)
	})
}

func TestReadWebhook(t *testing.T) {
	mbtest.WillReturnTestdata(t, "webhookObject.json", http.StatusOK)
	client := mbtest.Client(t)

	webhook, err := ReadWebhook(client, "whid")
	assert.NoError(t, err)
	assert.Equal(t, "chid", webhook.ChannelID)

	count := len(webhook.Events)
	assert.Equal(t, 2, count)

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
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/mynewwebhookurl", webhook.URL)

	assert.NotNil(t, webhook.UpdatedDatetime)
	assert.Equal(t, WebhookStatusDisabled, webhook.Status)

	mbtest.AssertEndpointCalled(t, http.MethodPatch, "/v1/webhooks/whid")
	mbtest.AssertTestdata(t, "webhookUpdateRequest.json", mbtest.Request.Body)
}
