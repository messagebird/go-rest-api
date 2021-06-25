package conversation

import (
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/v7/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	message, err := CreateMessage(client, "convid", &MessageCreateRequest{
		ChannelID: "chid",
		Content: &MessageContent{
			Text: "Hello world",
		},
		Type: MessageTypeText,
	})
	assert.NoError(t, err)
	assert.Equal(t, "mesid", message.ID)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/conversations/convid/messages")
	mbtest.AssertTestdata(t, "messageCreateRequest.json", mbtest.Request.Body)
}

func TestListMessages(t *testing.T) {
	t.Run("limit_offset", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "messageListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		messageList, err := ListMessages(client, "convid", &ListOptions{Limit: 20, Offset: 2})
		assert.NoError(t, err)

		assert.Equal(t, 2, messageList.Offset)

		assert.Equal(t, 20, messageList.Limit)

		assert.Equal(t, "mesid", messageList.Items[0].ID)

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations/convid/messages")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "limit=20&offset=2", query)
	})

	t.Run("all", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "allMessageListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		messageList, err := ListMessages(client, "convid", nil)
		assert.NoError(t, err)

		assert.Equal(t, 10, messageList.Limit)

		assert.Equal(t, 0, messageList.Offset)

		assert.Equal(t, "mesid", messageList.Items[0].ID)

		assert.Len(t, messageList.Items, 2)

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations/convid/messages")

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
