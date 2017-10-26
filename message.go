package messagebird

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

// TypeDetails is a hash with extra information.
// Is only used when a binary or premium message is sent.
type TypeDetails map[string]interface{}

// Message struct represents a message at MessageBird.com
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
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        Recipients
	Errors            []Error
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

// MessageParams provide additional message send options and used in URL as params.
type MessageParams struct {
	Type              string
	Reference         string
	Validity          int
	Gateway           int
	TypeDetails       TypeDetails
	DataCoding        string
	ScheduledDatetime time.Time
}

// MessageListParams provides additional message list options.
type MessageListParams struct {
	Originator string
	Direction  string
	Type       string
	Limit      int
	Offset     int
}

// paramsForMessage converts the specified MessageParams struct to a
// url.Values pointer and returns it.
func paramsForMessage(params *MessageParams) (*url.Values, error) {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams, nil
	}

	if params.Type != "" {
		urlParams.Set("type", params.Type)
		if params.Type == "flash" {
			urlParams.Set("mclass", "0")
		}
	}
	if params.Reference != "" {
		urlParams.Set("reference", params.Reference)
	}
	if params.Validity != 0 {
		urlParams.Set("validity", strconv.Itoa(params.Validity))
	}
	if params.Gateway != 0 {
		urlParams.Set("gateway", strconv.Itoa(params.Gateway))
	}

	for k, v := range params.TypeDetails {
		if vs, ok := v.(string); ok {
			urlParams.Set("typeDetails["+k+"]", vs)
		} else if vi, ok := v.(int); ok {
			urlParams.Set("typeDetails["+k+"]", strconv.Itoa(vi))
		} else {
			return nil, errors.New("Unknown type for typeDetails value")
		}
	}

	if params.DataCoding != "" {
		urlParams.Set("datacoding", params.DataCoding)
	}
	if params.ScheduledDatetime.Unix() > 0 {
		urlParams.Set("scheduledDatetime", params.ScheduledDatetime.Format(time.RFC3339))
	}

	return urlParams, nil
}

// paramsForMessageList converts the specified MessageListParams struct to a
// url.Values pointer and returns it.
func paramsForMessageList(params *MessageListParams) (*url.Values, error) {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams, nil
	}

	if params.Direction != "" {
		urlParams.Set("direction", params.Direction)
	}
	if params.Originator != "" {
		urlParams.Set("originator", params.Originator)
	}
	if params.Limit != 0 {
		urlParams.Set("limit", strconv.Itoa(params.Limit))
	}
	urlParams.Set("offset", strconv.Itoa(params.Offset))

	return urlParams, nil
}
