package conversation

import (
	messagebird "github.com/messagebird/go-rest-api/v8"
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/v8/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageSendResponse.json", http.StatusAccepted)
	client := mbtest.Client(t)

	message, err := SendMessage(client, &SendMessageRequest{
		To:   "+31624971134",
		From: "MessageBird",
		Type: MessageTypeText,
		Content: &MessageContent{
			Text: "Hello world",
		},
		ReportUrl: "https://myreport.site",
		Source:    map[string]interface{}{"name": "Valera"},
	})
	assert.NoError(t, err)
	t.Log(message)
	assert.Equal(t, "2e15efafec384e1c82e9842075e87beb", message.ID)
	assert.Equal(t, MessageStatusAccepted, message.Status)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/send")
}

func TestListMessages(t *testing.T) {
	t.Run("limit_offset", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "messageListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		messageList, err := ListMessages(client, &ListMessagesRequest{Ids: "5f3437fdb8444583aea093a047ac014b,4abc37fdb8444583aea093a047ac014c"})
		assert.NoError(t, err)

		assert.Equal(t, 2, messageList.Offset)

		assert.Equal(t, 20, messageList.Limit)

		assert.Equal(t, "mesid", messageList.Items[0].ID)

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/messages")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "ids=5f3437fdb8444583aea093a047ac014b%2C4abc37fdb8444583aea093a047ac014c", query)
	})

	t.Run("all", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "allMessageListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		messageList, err := ListMessages(client, nil)
		assert.NoError(t, err)

		assert.Equal(t, 10, messageList.Limit)

		assert.Equal(t, 0, messageList.Offset)

		assert.Equal(t, "mesid", messageList.Items[0].ID)

		assert.Len(t, messageList.Items, 2)

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/messages")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "", query)
	})
}

func TestListConversationMessages(t *testing.T) {
	conversationId := "5f3437fdb8444583aea093a047ac014b"

	t.Run("limit_offset", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "messageListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		messageList, err := ListConversationMessages(
			client,
			conversationId,
			&ListConversationMessagesRequest{messagebird.CommonPaginationRequest{Limit: 20, Offset: 2}, "sms,whatsapp,facebook"},
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, messageList.Offset)
		assert.Equal(t, 20, messageList.Limit)
		assert.Equal(t, "mesid", messageList.Items[0].ID)
		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations/"+conversationId+"/messages")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "excludePlatforms=sms%2Cwhatsapp%2Cfacebook&limit=20&offset=2", query)
	})

	t.Run("all", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "allMessageListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		messageList, err := ListConversationMessages(client, "5f3437fdb8444583aea093a047ac014b", nil)
		assert.NoError(t, err)
		assert.Equal(t, 0, messageList.Offset)
		assert.Equal(t, 10, messageList.Limit)
		assert.Equal(t, "mesid", messageList.Items[0].ID)
		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations/"+conversationId+"/messages")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "", query)
	})
}

func TestReadMessage(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	message, err := ReadMessage(client, "mesid")
	assert.NoError(t, err)

	assert.Equal(t, "Hello world", message.Content.Text)
	assert.Equal(t, MessageDirectionReceived, message.Direction)
	assert.Equal(t, MessageStatusFailed, message.Status)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/messages/mesid")
}
