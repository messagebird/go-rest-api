//
// Copyright (c) 2014 MessageBird B.V.
// All rights reserved.
//
// Author: Maurice Nonnekes <maurice@messagebird.com>

package messagebird

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

const (
	ClientVersion = "1.0.1"
	Endpoint      = "https://rest.messagebird.com"
)

type Recipient struct {
	Recipient      int
	Status         string
	StatusDatetime *time.Time
}

type Recipients struct {
	TotalCount               int
	TotalSentCount           int
	TotalDeliveredCount      int
	TotalDeliveryFailedCount int
	Items                    []Recipient
}

type Error struct {
	Code        int
	Description string
	Parameter   string
}

type Balance struct {
	Payment string
	Type    string
	Amount  int
	Errors  []Error
}

type HLR struct {
	Id              string
	HRef            string
	MSISDN          int
	Reference       string
	Status          string
	CreatedDatetime *time.Time
	StatusDatetime  *time.Time
	Errors          []Error
}

type Message struct {
	Id                string
	HRef              string
	Direction         string
	Type              string
	Originator        string
	Body              string
	Reference         string
	Validity          string
	Gateway           int
	TypeDetails       map[string]interface{}
	DataCoding        string
	MClass            int
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        Recipients
	Errors            []Error
}

type VoiceMessage struct {
	Id                string
	HRef              string
	Body              string
	Reference         string
	Language          string
	Voice             string
	Repeat            int
	IfMachine         string
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        Recipients
	Errors            []Error
}

type Client struct {
	AccessKey string
}

// Create a new Client.
func New(AccessKey string) *Client {
	return &Client{AccessKey}
}

// This function performs a call to MessageBird's HTTP API and expects JSON in
// return. It then tries to unmarshal the JSON body of the response to the
// specified struct.
func (c Client) request(v interface{}, path string, params *url.Values) error {
	var request *http.Request

	// Construct the URI of the request.
	uri, err := url.Parse(Endpoint + "/" + path)
	if err != nil {
		return err
	}

	// Construct a new request.
	if params == nil {
		request, err = http.NewRequest("GET", uri.String(), nil)
	} else {
		request, err = http.NewRequest("POST", uri.String(), strings.NewReader(params.Encode()))
	}
	if err != nil {
		return err
	}

	// Add a basic set of headers to the request.
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "AccessKey "+c.AccessKey)
	request.Header.Add("User-Agent", "MessageBird/ApiClient/"+ClientVersion+" Go/"+runtime.Version())

	// Add the Content-Type header if this is a POST request.
	if params != nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// Create an http.Client, execute the HTTP request and wait for a response.
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// Be sure to close the filedescriptor when all is done.
	defer response.Body.Close()

	// Read out the entire body.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Convert the JSON body to the specified struct.
	if err = json.Unmarshal(body, v); err != nil {
		return err
	}

	return nil
}

// This function retrieves your balance.
func (c Client) Balance() (*Balance, error) {
	balance := &Balance{}

	if err := c.request(balance, "balance", nil); err != nil {
		return nil, err
	}

	return balance, nil
}

// This function retrieves the information of a specific HLR.
func (c Client) HLR(id string) (*HLR, error) {
	hlr := &HLR{}

	if err := c.request(hlr, "hlr/"+id, nil); err != nil {
		return nil, err
	}

	return hlr, nil
}

// This function creates a new HLR.
func (c Client) CreateHLR(msisdn string, reference string) (*HLR, error) {
	params := &url.Values{
		"msisdn":    {msisdn},
		"reference": {reference}}

	hlr := &HLR{}

	if err := c.request(hlr, "hlr", params); err != nil {
		return nil, err
	}

	return hlr, nil
}

// This function retrieves the information of a specific message.
func (c Client) Message(id string) (*Message, error) {
	msg := &Message{}

	if err := c.request(msg, "messages/"+id, nil); err != nil {
		return nil, err
	}

	return msg, nil
}

// This function creates a new message, which is sent to one or more recipients.
func (c Client) CreateMessage(originator string, recipients []string, body string, params *url.Values) (*Message, error) {
	recips := strings.Join(recipients, ",")

	if params == nil {
		params = &url.Values{
			"originator": {originator},
			"body":       {body},
			"recipients": {recips}}
	} else {
		params.Set("originator", originator)
		params.Set("body", body)
		params.Set("recipients", recips)
	}

	msg := &Message{}

	if err := c.request(msg, "messages", params); err != nil {
		return nil, err
	}

	return msg, nil
}

// This function retrieves the information of a specific voice message.
func (c Client) VoiceMessage(id string) (*VoiceMessage, error) {
	vmsg := &VoiceMessage{}

	if err := c.request(vmsg, "voicemessages/"+id, nil); err != nil {
		return nil, err
	}

	return vmsg, nil
}

// This function creates a new voice message
func (c Client) CreateVoiceMessage(recipients []string, body string, params *url.Values) (*VoiceMessage, error) {
	recips := strings.Join(recipients, ",")

	if params == nil {
		params = &url.Values{
			"body":       {body},
			"recipients": {recips}}
	} else {
		params.Set("body", body)
		params.Set("recipients", recips)
	}

	vmsg := &VoiceMessage{}

	if err := c.request(vmsg, "voicemessages", params); err != nil {
		return nil, err
	}

	return vmsg, nil
}
