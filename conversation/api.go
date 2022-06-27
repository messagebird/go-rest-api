package conversation

import (
	"fmt"
	messagebird "github.com/messagebird/go-rest-api/v8"
)

const (
	// apiRoot is the absolute URL of the Converstations API. All paths are
	// relative to apiRoot (e.g.
	// https://conversations.messagebird.com/v1/webhooks).
	apiRoot = "https://conversations.messagebird.com/v1"

	whatsappSandboxAPIRoot = "https://whatsapp-sandbox.messagebird.com/v1"

	// path is the path for the Conversation resource, relative to apiRoot.
	path = "conversations"

	// startConversationPath is the path for starting new conversation
	startConversationPath = "start"

	// contactPath is the path for fetching a collection of conversations by contact ID
	contactPath = "contact"

	// messagesPath is the path for the Message resource, relative to apiRoot
	// and path.
	messagesPath = "messages"

	// sendMessagePath is the path for creating the Message resource relative to apiRoot
	sendMessagePath = "send"

	// webhooksPath is the path for the Webhook resource, relative to apiRoot.
	webhooksPath = "webhooks"
)

// request does the exact same thing as Client.Request. It does, however,
// prefix the path with the Conversation API's root. This ensures the client
// doesn't "handle" this for us: by default, it uses the REST API.
func request(c messagebird.ClientInterface, v interface{}, method, path string, data interface{}) error {
	var root string
	if c.IsFeatureEnabled(messagebird.FeatureConversationsAPIWhatsAppSandbox) {
		root = whatsappSandboxAPIRoot
	} else {
		root = apiRoot
	}
	return c.Request(v, method, fmt.Sprintf("%s/%s", root, path), data)
}
