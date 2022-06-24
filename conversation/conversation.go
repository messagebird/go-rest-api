package conversation

import (
	"fmt"
	"net/http"

	messagebird "github.com/messagebird/go-rest-api/v8"
)

// ListOptions can be used to set pagination options in List().
type ListOptions struct {
	Limit, Offset int
}

// StartRequest contains the request data for the Start endpoint.
type StartRequest struct {
	ChannelID string          `json:"channelId"`
	Content   *MessageContent `json:"content"`
	To        string          `json:"to"`
	Type      MessageType     `json:"type"`
}

// UpdateRequest contains the request data for the Update endpoint.
type UpdateRequest struct {
	Status ConversationStatus `json:"status"`
}

// DefaultListOptions provide a reasonable default for paginated endpoints.
var DefaultListOptions = &ListOptions{10, 0}

// List gets a collection of Conversations. Pagination can be set in options.
func List(c *messagebird.Client, options *ListOptions) (*ConversationList, error) {
	query := paginationQuery(options)

	convList := &ConversationList{}
	if err := request(c, convList, http.MethodGet, fmt.Sprintf("%s?%s", path, query), nil); err != nil {
		return nil, err
	}

	return convList, nil
}

// Read fetches a single Conversation based on its ID.
func Read(c *messagebird.Client, id string) (*Conversation, error) {
	conv := &Conversation{}
	if err := request(c, conv, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return conv, nil
}

// Start creates a conversation by sending an initial message. If an active
// conversation exists for the recipient, it is resumed.
func Start(c *messagebird.Client, req *StartRequest) (*Conversation, error) {
	conv := &Conversation{}
	if err := request(c, conv, http.MethodPost, path+"/start", req); err != nil {
		return nil, err
	}

	return conv, nil
}

// Update changes the conversation's status, so this can be used to (un)archive
// conversations.
func Update(c *messagebird.Client, id string, req *UpdateRequest) (*Conversation, error) {
	conv := &Conversation{}
	if err := request(c, conv, http.MethodPatch, path+"/"+id, req); err != nil {
		return nil, err
	}

	return conv, nil
}
