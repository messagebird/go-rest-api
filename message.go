package messagebird

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

type TypeDetails map[string]interface{}

type Message struct {
	Id                string
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

type MessageParams struct {
	Type              string
	Reference         string
	Validity          int
	Gateway           int
	TypeDetails       TypeDetails
	DataCoding        string
	ScheduledDatetime time.Time
}

type MessageQueryParams struct {
	Originator string
	Direction  string
	Limit      int
	Offset     int
}

type Messages struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []Message
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

func paramsForMessageQuery(params *MessageQueryParams) (*url.Values, error) {
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
