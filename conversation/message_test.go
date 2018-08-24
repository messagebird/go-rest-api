package conversation

import (
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/internal/mbtest"
)

func TestCreateMessage(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationMessageObject.json", http.StatusCreated)
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
	mbtest.AssertTestdata(t, "conversationCreateMessageRequest.json", mbtest.Request.Body)
}

func TestListMessages(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationListMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	messageList, err := ListMessages(client, "convid", DefaultListOptions)
	if err != nil {
		t.Fatalf("unexpected error listing Messages: %s", err)
	}

	if messageList.Limit != 10 {
		t.Fatalf("got %d, expected 10", messageList.Limit)
	}

	if messageList.Items[0].ID != "mesid" {
		t.Fatalf("got %s, expected mesid", messageList.Items[0].ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations/convid/messages")
}

func TestReadMessage(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationMessageObject.json", http.StatusOK)
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
