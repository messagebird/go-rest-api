package verify

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

// Verify object represents MessageBird server response.
type Verify struct {
	ID                 string
	HRef               string
	Reference          string
	Status             string
	Messages           map[string]string
	CreatedDatetime    *time.Time
	ValidUntilDatetime *time.Time
	Recipient          int
}

// VerifyParams handles optional verification parameters.
type VerifyParams struct {
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
}

// path represents the path to the Verify resource.
const path = "verify"

// Create generates a new One-Time-Password for one recipient.
func Create(c *messagebird.Client, recipient string, params *VerifyParams) (*Verify, error) {
	requestData, err := requestDataForVerify(recipient, params)
	if err != nil {
		return nil, err
	}

	verify := &Verify{}
	if err := c.Request(verify, http.MethodPost, path, requestData); err != nil {
		return nil, err
	}

	return verify, nil
}

// VerifyToken performs token value check against MessageBird API.
func VerifyToken(c *messagebird.Client, id, token string) (*Verify, error) {
	params := &url.Values{}
	params.Set("token", token)

	pathWithParams := path + "/" + id + "?" + params.Encode()

	verify := &Verify{}
	if err := c.Request(verify, http.MethodGet, pathWithParams, nil); err != nil {
		return nil, err
	}

	return verify, nil
}

func requestDataForVerify(recipient string, params *VerifyParams) (*verifyRequest, error) {
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

	return request, nil
}
