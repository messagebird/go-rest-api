package conversation

import (
	"fmt"
	messagebird "github.com/messagebird/go-rest-api/v7"
	"net/url"
	"strconv"
)

const (
	// apiRoot is the absolute URL of the Converstations API. All paths are
	// relative to apiRoot (e.g.
	// https://conversations.messagebird.com/v1/webhooks).
	apiRoot = "https://conversations.messagebird.com/v1"

	whatsappSandboxAPIRoot = "https://whatsapp-sandbox.messagebird.com/v1"

	// path is the path for the Conversation resource, relative to apiRoot.
	path = "conversations"

	// messagesPath is the path for the Message resource, relative to apiRoot
	// and path.
	messagesPath = "messages"

	// messagesPath is the path for the Message resource, relative to apiRoot
	// and path.
	sendMessage = "send"

	// webhooksPath is the path for the Webhook resource, relative to apiRoot.
	webhooksPath = "webhooks"
)

// ListRequestOptions can be used to set pagination options in List().
type ListRequestOptions struct {
	Limit, Offset int
}

// request does the exact same thing as Client.Request. It does, however,
// prefix the path with the Conversation API's root. This ensures the client
// doesn't "handle" this for us: by default, it uses the REST API.
func request(c *messagebird.Client, v interface{}, method, path string, data interface{}) error {
	var root string
	if c.IsFeatureEnabled(messagebird.FeatureConversationsAPIWhatsAppSandbox) {
		root = whatsappSandboxAPIRoot
	} else {
		root = apiRoot
	}
	return c.Request(v, method, fmt.Sprintf("%s/%s", root, path), data)
}

// paginationQuery builds the query string for paginated endpoints.
func paginationQuery(options *ListRequestOptions) string {
	if options == nil {
		return ""
	}

	query := url.Values{}
	query.Set("limit", strconv.Itoa(options.Limit))
	query.Set("offset", strconv.Itoa(options.Offset))

	return query.Encode()
}
