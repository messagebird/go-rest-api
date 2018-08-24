package conversation

import (
	"fmt"
	"net/http"

	messagebird "github.com/messagebird/go-rest-api"
)

type MessageCreateRequest struct {
	ChannelID string          `json:"channelid"`
	Content   *MessageContent `json:"content"`
	Type      MessageType     `json:"type"`
}

// CreateMessage sends a new message to the specified conversation. To create a
// new conversation and send an initial message, use conversation.Start().
func CreateMessage(c *messagebird.Client, conversationID string, req *MessageCreateRequest) (*Message, error) {
	uri := fmt.Sprintf("%s/%s/%s", path, conversationID, messagesPath)

	message := &Message{}
	if err := request(c, message, http.MethodPost, uri, req); err != nil {
		return nil, err
	}

	return message, nil
}

// ListMessages gets a collection of messages from a conversation. Pagination
// can be set in the options.
func ListMessages(c *messagebird.Client, conversationID string, options *ListOptions) (*MessageList, error) {
	uri := fmt.Sprintf("%s/%s/%s", path, conversationID, messagesPath)

	messageList := &MessageList{}
	if err := request(c, messageList, http.MethodGet, uri, options); err != nil {
		return nil, err
	}

	return messageList, nil
}

// ReadMessage gets a single message based on its ID.
func ReadMessage(c *messagebird.Client, messageID string) (*Message, error) {
	message := &Message{}
	if err := request(c, message, http.MethodGet, messagesPath+"/"+messageID, nil); err != nil {
		return nil, err
	}

	return message, nil
}
