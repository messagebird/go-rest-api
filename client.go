//
// Copyright (c) 2014 MessageBird B.V.
// All rights reserved.
//
// Author: Maurice Nonnekes <maurice@messagebird.com>

// Package messagebird is an official library for interacting with MessageBird.com API.
// The MessageBird API connects your website or application to operators around the world. With our API you can integrate SMS, Chat & Voice.
// More documentation you can find on the MessageBird developers portal: https://developers.messagebird.com/
package messagebird

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

const (
	// ClientVersion is used in User-Agent request header to provide server with API level.
	ClientVersion = "4.2.0"

	// Endpoint points you to MessageBird REST API.
	Endpoint = "https://rest.messagebird.com"
)

const (
	// HLRPath represents the path to the HLR resource.
	HLRPath = "hlr"
	// MessagePath represents the path to the Message resource.
	MessagePath = "messages"
	// MMSPath represents the path to the MMS resource.
	MMSPath = "mms"
	// VoiceMessagePath represents the path to the VoiceMessage resource.
	VoiceMessagePath = "voicemessages"
	// VerifyPath represents the path to the Verify resource.
	VerifyPath = "verify"
	// LookupPath represents the path to the Lookup resource.
	LookupPath = "lookup"
)

var (
	// ErrResponse is returned when we were able to contact API but request was not successful and contains error details.
	ErrResponse = errors.New("The MessageBird API returned an error")

	// ErrUnexpectedResponse is used when there was an internal server error and nothing can be done at this point.
	ErrUnexpectedResponse = errors.New("The MessageBird API is currently unavailable")
)

// Client is used to access API with a given key.
// Uses standard lib HTTP client internally, so should be reused instead of created as needed and it is safe for concurrent use.
type Client struct {
	AccessKey  string       // The API access key
	HTTPClient *http.Client // The HTTP client to send requests on
	DebugLog   *log.Logger  // Optional logger for debugging purposes
}

// New creates a new MessageBird client object.
func New(AccessKey string) *Client {
	return &Client{AccessKey: AccessKey, HTTPClient: &http.Client{}}
}

func (c *Client) request(v interface{}, path string, params *url.Values) error {
	uri, err := url.Parse(Endpoint + "/" + path)
	if err != nil {
		return err
	}

	var request *http.Request
	if params != nil {
		body := params.Encode()
		if request, err = http.NewRequest("POST", uri.String(), strings.NewReader(body)); err != nil {
			return err
		}

		if c.DebugLog != nil {
			if unescapedBody, queryError := url.QueryUnescape(body); queryError == nil {
				log.Printf("HTTP REQUEST: POST %s %s", uri.String(), unescapedBody)
			} else {
				log.Printf("HTTP REQUEST: POST %s %s", uri.String(), body)
			}
		}

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		if request, err = http.NewRequest("GET", uri.String(), nil); err != nil {
			return err
		}

		if c.DebugLog != nil {
			log.Printf("HTTP REQUEST: GET %s", uri.String())
		}
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "AccessKey "+c.AccessKey)
	request.Header.Add("User-Agent", "MessageBird/ApiClient/"+ClientVersion+" Go/"+runtime.Version())

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if c.DebugLog != nil {
		log.Printf("HTTP RESPONSE: %s", string(responseBody))
	}

	// Status code 500 is a server error and means nothing can be done at this
	// point.
	if response.StatusCode == 500 {
		return ErrUnexpectedResponse
	}

	if err = json.Unmarshal(responseBody, &v); err != nil {
		return err
	}

	// Status codes 200 and 201 are indicative of being able to convert the
	// response body to the struct that was specified.
	if response.StatusCode == 200 || response.StatusCode == 201 {
		return nil
	}

	// Anything else than a 200/201/500 should be a JSON error.
	return ErrResponse
}

// Balance returns the balance information for the account that is associated
// with the access key.
func (c *Client) Balance() (*Balance, error) {
	balance := &Balance{}
	if err := c.request(balance, "balance", nil); err != nil {
		if err == ErrResponse {
			return balance, err
		}

		return nil, err
	}

	return balance, nil
}

// HLR looks up an existing HLR object for the specified id that was previously
// created by the NewHLR function.
func (c *Client) HLR(id string) (*HLR, error) {
	hlr := &HLR{}
	if err := c.request(hlr, HLRPath+"/"+id, nil); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// HLRs lists all HLR objects that were previously created by the NewHLR
// function.
func (c *Client) HLRs() (*HLRList, error) {
	hlrList := &HLRList{}
	if err := c.request(hlrList, HLRPath, nil); err != nil {
		if err == ErrResponse {
			return hlrList, err
		}

		return nil, err
	}

	return hlrList, nil
}

// NewHLR retrieves the information of an existing HLR.
func (c *Client) NewHLR(msisdn, reference string) (*HLR, error) {
	params := &url.Values{
		"msisdn":    {msisdn},
		"reference": {reference},
	}

	hlr := &HLR{}
	if err := c.request(hlr, HLRPath, params); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// Message retrieves the information of an existing Message.
func (c *Client) Message(id string) (*Message, error) {
	message := &Message{}
	if err := c.request(message, MessagePath+"/"+id, nil); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// Messages retrieves all messages of the user represented as a MessageList object.
func (c *Client) Messages(msgListParams *MessageListParams) (*MessageList, error) {
	messageList := &MessageList{}
	params, err := paramsForMessageList(msgListParams)
	if err != nil {
		return messageList, err
	}

	if err := c.request(messageList, MessagePath+"?"+params.Encode(), nil); err != nil {
		if err == ErrResponse {
			return messageList, err
		}

		return nil, err
	}

	return messageList, nil
}

// NewMessage creates a new message for one or more recipients.
func (c *Client) NewMessage(originator string, recipients []string, body string, msgParams *MessageParams) (*Message, error) {
	params, err := paramsForMessage(msgParams)
	if err != nil {
		return nil, err
	}

	params.Set("originator", originator)
	params.Set("body", body)
	params.Set("recipients", strings.Join(recipients, ","))

	message := &Message{}
	if err := c.request(message, MessagePath, params); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// MMSMessage retrieves the information of an existing MmsMessage.
func (c *Client) MMSMessage(id string) (*MMSMessage, error) {
	mmsMessage := &MMSMessage{}
	if err := c.request(mmsMessage, MMSPath+"/"+id, nil); err != nil {
		if err == ErrResponse {
			return mmsMessage, err
		}

		return nil, err
	}

	return mmsMessage, nil
}

// NewMMSMessage creates a new MMS message for one or more recipients.
func (c *Client) NewMMSMessage(originator string, recipients []string, msgParams *MMSMessageParams) (*MMSMessage, error) {
	params, err := paramsForMMSMessage(msgParams)
	if err != nil {
		return nil, err
	}

	params.Set("originator", originator)
	params.Set("recipients", strings.Join(recipients, ","))

	mmsMessage := &MMSMessage{}
	if err := c.request(mmsMessage, MMSPath, params); err != nil {
		if err == ErrResponse {
			return mmsMessage, err
		}

		return nil, err
	}

	return mmsMessage, nil
}

// VoiceMessage retrieves the information of an existing VoiceMessage.
func (c *Client) VoiceMessage(id string) (*VoiceMessage, error) {
	message := &VoiceMessage{}
	if err := c.request(message, VoiceMessagePath+"/"+id, nil); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// VoiceMessages retrieves all VoiceMessages of the user.
func (c *Client) VoiceMessages() (*VoiceMessageList, error) {
	messageList := &VoiceMessageList{}
	if err := c.request(messageList, VoiceMessagePath, nil); err != nil {
		if err == ErrResponse {
			return messageList, err
		}

		return nil, err
	}

	return messageList, nil
}

// NewVoiceMessage creates a new voice message for one or more recipients.
func (c *Client) NewVoiceMessage(recipients []string, body string, params *VoiceMessageParams) (*VoiceMessage, error) {
	urlParams := paramsForVoiceMessage(params)
	urlParams.Set("body", body)
	urlParams.Set("recipients", strings.Join(recipients, ","))

	message := &VoiceMessage{}
	if err := c.request(message, VoiceMessagePath, urlParams); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// NewVerify generates a new One-Time-Password for one recipient.
func (c *Client) NewVerify(recipient string, params *VerifyParams) (*Verify, error) {
	urlParams := paramsForVerify(params)
	urlParams.Set("recipient", recipient)

	verify := &Verify{}
	if err := c.request(verify, VerifyPath, urlParams); err != nil {
		if err == ErrResponse {
			return verify, err
		}

		return nil, err
	}

	return verify, nil
}

// VerifyToken performs token value check against MessageBird API.
func (c *Client) VerifyToken(id, token string) (*Verify, error) {
	params := &url.Values{}
	params.Set("token", token)

	path := VerifyPath + "/" + id + "?" + params.Encode()

	verify := &Verify{}
	if err := c.request(verify, path, nil); err != nil {
		if err == ErrResponse {
			return verify, err
		}

		return nil, err
	}

	return verify, nil
}

// Lookup performs a new lookup for the specified number.
func (c *Client) Lookup(phoneNumber string, params *LookupParams) (*Lookup, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "?" + urlParams.Encode()

	lookup := &Lookup{}
	if err := c.request(lookup, path, nil); err != nil {
		if err == ErrResponse {
			return lookup, err
		}

		return nil, err
	}

	return lookup, nil
}

// NewLookupHLR creates a new HLR lookup for the specified number.
func (c *Client) NewLookupHLR(phoneNumber string, params *LookupParams) (*HLR, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "/" + HLRPath

	hlr := &HLR{}
	if err := c.request(hlr, path, urlParams); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// LookupHLR performs a HLR lookup for the specified number.
func (c *Client) LookupHLR(phoneNumber string, params *LookupParams) (*HLR, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "/" + HLRPath + "?" + urlParams.Encode()

	hlr := &HLR{}
	if err := c.request(hlr, path, nil); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}
