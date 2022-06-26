package conversation

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v7"
)

const (
	// MessageDirectionReceived indicates an inbound message received from the customer.
	MessageDirectionReceived MessageDirection = "received"

	// MessageDirectionSent indicates an outbound message sent from the API.
	MessageDirectionSent MessageDirection = "sent"
)

const (
	MessageStatusAccepted        MessageStatus = "accepted"
	MessageStatusPending         MessageStatus = "pending"
	MessageStatusSent            MessageStatus = "sent"
	MessageStatusRejected        MessageStatus = "rejected"
	MessageStatusFailed          MessageStatus = "failed"
	MessageStatusRead            MessageStatus = "read"
	MessageStatusReceived        MessageStatus = "received"
	MessageStatusDeleted         MessageStatus = "deleted"
	MessageStatusUnknown         MessageStatus = "unknown"
	MessageStatusTransmitted     MessageStatus = "transmitted"
	MessageStatusDeliveryFailed  MessageStatus = "delivery_failed"
	MessageStatusBuffered        MessageStatus = "buffered"
	MessageStatusExpired         MessageStatus = "expired"
	MessageStatusClicked         MessageStatus = "clicked"
	MessageStatusOpened          MessageStatus = "opened"
	MessageStatusBounce          MessageStatus = "bounce"
	MessageStatusSpamComplaint   MessageStatus = "spam_complaint"
	MessageStatusOutOfBounded    MessageStatus = "out_of_bounded"
	MessageStatusDelayed         MessageStatus = "delayed"
	MessageStatusListUnsubscribe MessageStatus = "list_unsubscribe"
	MessageStatusDispatched      MessageStatus = "dispatched"
)

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeFile     MessageType = "file"
	MessageTypeLocation MessageType = "location"
	MessageTypeEvent    MessageType = "event"
	MessageTypeRich     MessageType = "rich"
	MessageTypeMenu     MessageType = "menu"
	MessageTypeButtons  MessageType = "buttons"
	MessageTypeLink     MessageType = "link"

	MessageTypeHSM             MessageType = "hsm"
	MessageTypeWhatsAppSticker MessageType = "whatsappSticker"
	MessageTypeInteractive     MessageType = "interactive"
	MessageTypeWhatsappOrder   MessageType = "whatsappOrder"
	MessageTypeWhatsappText    MessageType = "whatsappText"

	MessageTypeExternalAttachment MessageType = "externalAttachment"
	MessageTypeEmail              MessageType = "email"
)

// MessageType indicates what kind of content a Message has, e.g. audio or text.
type MessageType string

// MessageStatus is a field set by the API. It indicates what the state of the
// message is, e.g. whether it has been successfully delivered or read.
type MessageStatus string
type MessageRecipient string
type MessageTag string
type MessageDirection string

type MessageCreateRequest struct {
	ChannelID string          `json:"channelid"`
	Content   *MessageContent `json:"content"`
	Type      MessageType     `json:"type"`
}

type Message struct {
	ID              string
	ConversationID  string
	ChannelID       string
	Platform        string
	To              MessageRecipient
	From            string
	Direction       MessageDirection
	Status          MessageStatus
	Type            MessageType
	Content         *MessageContent
	CreatedDatetime *time.Time
	UpdatedDatetime *time.Time
	Source          map[string]interface{}
	Tag             MessageTag
	Fallback        *Fallback
	TTL             string
}

// MessageContent holds a message's actual content. Only one field can be set
// per request.
type MessageContent struct {
	Audio    *Audio    `json:"audio,omitempty"`
	File     *File     `json:"file,omitempty"`
	Image    *Image    `json:"image,omitempty"`
	Location *Location `json:"location,omitempty"`
	Video    *Video    `json:"video,omitempty"`
	Text     string    `json:"text,omitempty"`

	// HSM is a highly structured message for WhatsApp. Its definition lives in
	// hsm.go.
	HSM *HSM `json:"hsm,omitempty"`

	Interactive     *WhatsAppInteractive `json:"interactive,omitempty"`
	WhatsAppSticker *WhatsAppSticker     `json:"whatsappSticker,omitempty"`
	WhatsAppOrder   *WhatsAppOrder       `json:"whatsappOrder,omitempty"`
	WhatsAppText    *WhatsAppText        `json:"whatsappText,omitempty"`

	FacebookQuickReply      *FacebookMessage `json:"facebookQuickReply,omitempty"`
	FacebookMediaTemplate   *FacebookMessage `json:"facebookMediaTemplate,omitempty"`
	FacebookGenericTemplate *FacebookMessage `json:"facebookGenericTemplate,omitempty"`

	Email               *Email   `json:"email,omitempty"`
	ExternalAttachments []*Media `json:"externalAttachments,omitempty"`
	DisableUrlPreview   bool     `json:"disableUrlPreview,omitempty"`
}

type Audio Media
type File Media
type Image Media
type Video Media

type Media struct {
	URL     string `json:"url"`
	Caption string `json:"caption,omitempty"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Fallback struct {
	From  string `json:"from"`
	After string `json:"after"`
}

type MessageList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []*Message
}

// SendMessageRequest contains the request data for the Reply endpoint.
type SendMessageRequest struct {
	To        string                 `json:"to"`
	From      string                 `json:"from"`
	Type      MessageType            `json:"type"`
	Content   *MessageContent        `json:"content"`
	ReportUrl string                 `json:"reportUrl,omitempty"`
	Fallback  *Fallback              `json:"fallback,omitempty"`
	Source    map[string]interface{} `json:"source,omitempty"`
	Tag       MessageTag             `json:"tag,omitempty"`
	TrackId   string                 `json:"trackId,omitempty"`
	TTL       string                 `json:"ttl,omitempty"`
}

type ListConversationMessagesRequest struct {
	PaginationRequest
	ExcludePlatforms string
}

func (lr *ListConversationMessagesRequest) GetParams() string {
	if lr == nil {
		return ""
	}

	query := url.Values{}

	query.Set("limit", strconv.Itoa(lr.Limit))
	query.Set("offset", strconv.Itoa(lr.Offset))
	query.Set("excludePlatforms", lr.ExcludePlatforms)

	return query.Encode()
}

type ListMessagesRequest struct {
	Ids  string
	From *time.Time
}

func (lr *ListMessagesRequest) GetParams() string {
	if lr == nil {
		return ""
	}

	query := url.Values{}

	query.Set("ids", lr.Ids)
	query.Set("ids", lr.From.Format(time.RFC3339))

	return query.Encode()
}

// SendMessage send a message to a specific recipient in a specific platform.
// If an active conversation already exists for the recipient, the conversation will be resumed.
// In case there's no active conversation a new one is created.
func SendMessage(c *messagebird.Client, options *SendMessageRequest) (*Message, error) {
	message := &Message{}
	if err := request(c, message, http.MethodPost, sendMessagePath, options); err != nil {
		return nil, err
	}

	return message, nil
}

// ListConversationMessages gets a collection of messages from a conversation.
// Pagination can be set in the options.
func ListConversationMessages(c *messagebird.Client, conversationID string, options *ListConversationMessagesRequest) (*MessageList, error) {
	uri := fmt.Sprintf("%s/%s/%s?%s", path, conversationID, messagesPath, options.GetParams())

	messageList := &MessageList{}
	if err := request(c, messageList, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return messageList, nil
}

// ListMessages gets a collection of messages from a conversation.
// Pagination can be set in the options.
func ListMessages(c *messagebird.Client, options *ListMessagesRequest) (*MessageList, error) {
	uri := fmt.Sprintf("%s/%s?%s", path, messagesPath, options.GetParams())

	messageList := &MessageList{}
	if err := request(c, messageList, http.MethodGet, uri, nil); err != nil {
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
