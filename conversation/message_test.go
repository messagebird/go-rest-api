package conversation

import (
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
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
	if err != nil {
		t.Fatalf("unexpected error creating Message: %s", err)
	}

	if message.ID != "mesid" {
		t.Fatalf("got %s, expected mesid", message.ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/conversations/convid/messages")
	mbtest.AssertTestdata(t, "messageCreateRequest.json", mbtest.Request.Body)
}

func TestListMessages(t *testing.T) {
	t.Run("limit_offset", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "messageListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		messageList, err := ListMessages(client, "convid", &ListOptions{Limit: 20, Offset: 2})
		if err != nil {
			t.Fatalf("unexpected error listing Messages: %s", err)
		}

		if messageList.Offset != 2 {
			t.Fatalf("got %d, expected 2", messageList.Offset)
		}

		if messageList.Limit != 20 {
			t.Fatalf("got %d, expected 20", messageList.Limit)
		}

		if messageList.Items[0].ID != "mesid" {
			t.Fatalf("got %s, expected mesid", messageList.Items[0].ID)
		}

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations/convid/messages")

		if query := mbtest.Request.URL.RawQuery; query != "limit=20&offset=2" {
			t.Fatalf("got %s, expected limit=10&offset=0", query)
		}
	})

	t.Run("all", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "allMessageListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		messageList, err := ListMessages(client, "convid", nil)
		if err != nil {
			t.Fatalf("unexpected error listing Messages: %s", err)
		}

		if messageList.Limit != 10 {
			t.Fatalf("got %d, expected 10", messageList.Limit)
		}

		if messageList.Offset != 0 {
			t.Fatalf("got %d, expected 0", messageList.Offset)
		}

		if messageList.Items[0].ID != "mesid" {
			t.Fatalf("got %s, expected mesid", messageList.Items[0].ID)
		}

		if len(messageList.Items) != 2 {
			t.Fatalf("got %d, expected 2", len(messageList.Items))
		}

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations/convid/messages")

		if query := mbtest.Request.URL.RawQuery; query != "" {
			t.Fatalf("got %s, expected empty", query)
		}
	})
}

func TestReadMessage(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	message, err := ReadMessage(client, "mesid")
	if err != nil {
		t.Fatalf("unexpected error creating Message: %s", err)
	}

	if message.Content.Text != "Hello world" {
		t.Fatalf("got %s, expected Hello world", message.Content.Text)
	}

	if message.Direction != MessageDirectionReceived {
		t.Fatalf("got %s, expected received", message.Direction)
	}

	if message.Status != MessageStatusFailed {
		t.Fatalf("got %s, expected failed", message.Status)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/messages/mesid")
}
