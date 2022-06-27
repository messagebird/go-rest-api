package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v7"
)

const (
	// path represents the path to the Verify resource.
	path = "verify"

	emailMessagesPath = path + "/messages/email"
)

type VerifyMessage struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// Params handles optional verification parameters.
type Params struct {
	Originator  string
	Reference   string
	Type        string
	Template    string
	DataCoding  string
	ReportURL   string
	Voice       string
	Language    string
	Timeout     int
	TokenLength int
	Subject     string
}

type verifyRequest struct {
	Recipient   string `json:"recipient"`
	Originator  string `json:"originator,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Type        string `json:"type,omitempty"`
	Template    string `json:"template,omitempty"`
	DataCoding  string `json:"dataCoding,omitempty"`
	ReportURL   string `json:"reportUrl,omitempty"`
	Voice       string `json:"voice,omitempty"`
	Language    string `json:"language,omitempty"`
	Timeout     int    `json:"timeout,omitempty"`
	TokenLength int    `json:"tokenLength,omitempty"`
	Subject     string `json:"subject,omitempty"`
}

// Create generates a new One-Time-Password for one recipient.
func Create(c *messagebird.Client, recipient string, params *Params) (*Verify, error) {
	requestData, err := paramsToVerifyRequest(recipient, params)
	if err != nil {
		return nil, err
	}

	verify := &Verify{}
	if err := c.Request(verify, http.MethodPost, path, requestData); err != nil {
		return nil, err
	}

	return verify, nil
}

// Delete deletes an existing Verify object by its ID.
func Delete(c *messagebird.Client, id string) error {
	return c.Request(nil, http.MethodDelete, path+"/"+id, nil)
}

// Read retrieves an existing Verify object by its ID.
func Read(c *messagebird.Client, id string) (*Verify, error) {
	verify := &Verify{}

	if err := c.Request(verify, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return verify, nil
}

// VerifyToken performs token value check against MessageBird API.
func VerifyToken(c *messagebird.Client, id, token string) (*Verify, error) {
	pathWithParams := path + "/" + id + "?token=" + token

	verify := &Verify{}
	if err := c.Request(verify, http.MethodGet, pathWithParams, nil); err != nil {
		return nil, err
	}

	return verify, nil
}

func ReadVerifyEmailMessage(c *messagebird.Client, id string) (*VerifyMessage, error) {
	verifyMessage := &VerifyMessage{}
	if err := c.Request(verifyMessage, http.MethodGet, emailMessagesPath+"/"+id, nil); err != nil {
		return nil, err
	}

	return verifyMessage, nil
}

func paramsToVerifyRequest(recipient string, params *Params) (*verifyRequest, error) {
	if recipient == "" {
		return nil, errors.New("recipient is required")
	}

	request := &verifyRequest{
		Recipient: recipient,
	}

	if params == nil {
		return request, nil
	}

	request.Originator = params.Originator
	request.Reference = params.Reference
	request.Type = params.Type
	request.Template = params.Template
	request.DataCoding = params.DataCoding
	request.ReportURL = params.ReportURL
	request.Voice = params.Voice
	request.Language = params.Language
	request.Timeout = params.Timeout
	request.TokenLength = params.TokenLength
	request.Subject = params.Subject

	return request, nil
}

// Verify object represents MessageBird server response.
type Verify struct {
	ID                 string
	HRef               string
	Reference          string
	Status             string
	Messages           map[string]string
	CreatedDatetime    *time.Time
	ValidUntilDatetime *time.Time
	Recipient          string
}

// UnmarshalJSON
// The type of the Verify.Recipient object changed from int to string
// but the api still returns a recipent numeric value whne sms type is used.
// This was the best way to ensure backward compatibility with the previous versions.
func (v *Verify) UnmarshalJSON(b []byte) error {
	if v == nil {
		return errors.New("cannot unmarshal to nil pointer")
	}

	// Need a type alias so we get a type the same memory layout, but without Verify's method set.
	// Otherwise encoding/json will recursively invoke this UnmarshalJSON() implementation.
	type Alias Verify
	var wrapper struct {
		Alias
		Recipient interface{}
	}
	if err := json.Unmarshal(b, &wrapper); err != nil {
		return err
	}

	switch wrapper.Recipient.(type) {
	case float64:
		const noExponent = 'f'
		const precision = -1
		const bitSize = 64
		asFloat := wrapper.Recipient.(float64)

		wrapper.Alias.Recipient = strconv.FormatFloat(asFloat, noExponent, precision, bitSize)
	case string:
		wrapper.Alias.Recipient = wrapper.Recipient.(string)
	default:
		return fmt.Errorf("recipient is unknown type %T", wrapper.Recipient)
	}

	*v = Verify(wrapper.Alias)

	return nil
}
