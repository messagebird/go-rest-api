package messagebird

import (
	"errors"
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

type messageRequest struct {
	Originator        string      `json:"originator"`
	Body              string      `json:"body"`
	Recipients        []string    `json:"recipients"`
	Type              string      `json:"type,omitempty"`
	Reference         string      `json:"reference,omitempty"`
	Validity          int         `json:"validity,omitempty"`
	Gateway           int         `json:"gateway,omitempty"`
	TypeDetails       TypeDetails `json:"typeDetails,omitempty"`
	DataCoding        string      `json:"dataCoding,omitempty"`
	MClass            int         `json:"mclass,omitempty"`
	ScheduledDatetime string      `json:"scheduledDatetime,omitempty"`
}

func requestDataForMessage(originator string, recipients []string, body string, params *MessageParams) (*messageRequest, error) {
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
	request.Reference = params.Reference
	request.Validity = params.Validity
	request.Gateway = params.Gateway
	request.TypeDetails = params.TypeDetails
	request.DataCoding = params.DataCoding
	request.ScheduledDatetime = params.ScheduledDatetime.Format(time.RFC3339)

	return request, nil
}
