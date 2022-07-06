package sms

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v9"
)

// TypeDetails is a hash with extra information.
// Is only used when a binary or premium message is sent.
type TypeDetails map[string]interface{}

// Message struct represents a message at messagebird.com.
type Message struct {
	ID                string
	HRef              string
	Direction         string
	Type              string
	Originator        string
	Body              string
	Reference         string
	Validity          *int
	Gateway           int
	TypeDetails       TypeDetails
	DataCoding        string
	MClass            int
	ReportURL         string
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        messagebird.Recipients
}

// MessageList represents a list of Messages.
type MessageList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      map[string]*string
	Items      []Message
}

// Params provide additional message send options and used in URL as params.
type Params struct {
	GroupIds          []string
	Type              string
	Reference         string
	Validity          int
	Gateway           int
	TypeDetails       TypeDetails
	DataCoding        string
	ReportURL         string
	ScheduledDatetime time.Time
	ShortenURLs       bool
}

// ListParams provides additional message list options.
type ListParams struct {
	Originator string
	Direction  string
	Type       string
	Status     string
	Limit      int
	Offset     int
}

func (lp *ListParams) QueryParams() string {
	if lp == nil {
		return ""
	}

	query := url.Values{}

	if len(lp.Originator) > 0 {
		query.Set("originator", lp.Originator)
	}

	if len(lp.Direction) > 0 {
		query.Set("direction", lp.Direction)
	}

	if len(lp.Type) > 0 {
		query.Set("type", lp.Type)
	}

	if len(lp.Status) > 0 {
		query.Set("status", lp.Status)
	}

	if lp.Limit > 0 {
		query.Set("limit", strconv.Itoa(lp.Limit))
	}

	if lp.Offset > 0 {
		query.Set("offset", strconv.Itoa(lp.Offset))
	}

	return query.Encode()
}

type messageRequest struct {
	Originator        string      `json:"originator"`
	Body              string      `json:"body"`
	Recipients        []string    `json:"recipients"`
	GroupIds          []string    `json:"groupIds"`
	Type              string      `json:"type,omitempty"`
	Reference         string      `json:"reference,omitempty"`
	Validity          int         `json:"validity,omitempty"`
	Gateway           int         `json:"gateway,omitempty"`
	TypeDetails       TypeDetails `json:"typeDetails,omitempty"`
	DataCoding        string      `json:"datacoding,omitempty"`
	MClass            int         `json:"mclass,omitempty"`
	ShortenURLs       bool        `json:"shortenUrls"`
	ReportURL         string      `json:"reportUrl,omitempty"`
	ScheduledDatetime string      `json:"scheduledDatetime,omitempty"`
}

// path represents the path to the Message resource.
const path = "messages"

// Read retrieves the information of an existing Message.
func Read(c messagebird.Client, id string) (*Message, error) {
	message := &Message{}
	if err := c.Request(message, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return message, nil
}

// Delete Cancel sending Scheduled Sms.
// Return true if have been successfully deleted.
func Delete(c messagebird.Client, id string) error {
	return c.Request(&Message{}, http.MethodDelete, path+"/"+id, nil)
}

// List retrieves all messages of the user represented as a MessageList object.
func List(c messagebird.Client, params *ListParams) (*MessageList, error) {
	messageList := &MessageList{}

	if err := c.Request(messageList, http.MethodGet, path+"?"+params.QueryParams(), nil); err != nil {
		return nil, err
	}

	return messageList, nil
}

// Create creates a new message for one or more recipients.
func Create(c messagebird.Client, originator string, recipients []string, body string, msgParams *Params) (*Message, error) {
	requestData, err := paramsToRequest(originator, recipients, body, msgParams)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	if err := c.Request(message, http.MethodPost, path, requestData); err != nil {
		return nil, err
	}

	return message, nil
}

func paramsToRequest(originator string, recipients []string, body string, params *Params) (*messageRequest, error) {
	if originator == "" {
		return nil, errors.New("originator is required")
	}
	if len(recipients) == 0 {
		return nil, errors.New("at least 1 recipient is required")
	}
	if body == "" {
		return nil, errors.New("body is required")
	}

	request := &messageRequest{
		Originator: originator,
		Recipients: recipients,
		Body:       body,
	}

	if params == nil {
		return request, nil
	}

	request.Type = params.Type
	if request.Type == "flash" {
		request.MClass = 0
	} else {
		request.MClass = 1
	}

	if !params.ScheduledDatetime.IsZero() {
		request.ScheduledDatetime = params.ScheduledDatetime.Format(time.RFC3339)
	}

	request.GroupIds = params.GroupIds
	request.Reference = params.Reference
	request.Validity = params.Validity
	request.Gateway = params.Gateway
	request.TypeDetails = params.TypeDetails
	request.DataCoding = params.DataCoding
	request.ReportURL = params.ReportURL
	request.ShortenURLs = params.ShortenURLs

	return request, nil
}
