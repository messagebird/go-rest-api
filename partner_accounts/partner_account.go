package partner_accounts

import (
	"encoding/json"
	"fmt"
	messagebird "github.com/messagebird/go-rest-api/v8"
	"net/http"
)

const (
	// apiRoot is the absolute URL of the Converstations API. All paths are
	// relative to apiRoot (e.g.
	// https://conversations.messagebird.com/v1/webhooks).
	apiRoot = "https://partner-accounts.messagebird.com/v1"

	childAccountsPath = "child-accounts"
)

func init() {
	// The Partner Accounts API returns errors in a format that slightly differs from other APIs (as Voice API).
	// Here we instruct package messagebird to use our custom
	// voice.errorReader func, which has access to voice.ErrorResponse, to
	// unmarshal those. Package messagebird must not import the voice package to
	// safeguard against import cycles, so it can not use voice.ErrorResponse
	// directly.
	messagebird.SetErrorReader(errorReader)
}

type ErrorResponse struct {
	Type, Title, Detail string
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}

// errorReader takes a []byte representation of a Voice API JSON error and
// parses it to a voice.ErrorResponse.
func errorReader(b []byte) error {
	var er ErrorResponse
	if err := json.Unmarshal(b, &er); err != nil {
		return fmt.Errorf("encoding/json: Unmarshal: %v", err)
	}

	return er
}

type Account struct {
	ID                 int
	Name               string
	AccessKeys         []*AccessKey
	SigningKey         string
	InvoiceAggregation bool
}

type AccessKey struct {
	ID   string
	Key  string
	Mode string
}

type Accounts []Account

type createChildAccountRequest struct {
	name string
}

func CreateChildAccount(c messagebird.MessageBirdClient, name string) (*Account, error) {
	a := &Account{}

	req := &createChildAccountRequest{name}

	if err := c.Request(a, http.MethodPost, apiRoot+"/"+childAccountsPath, req); err != nil {
		return nil, err
	}

	return a, nil
}

func UpdateChildAccount(c messagebird.MessageBirdClient, id, name string) (*Account, error) {
	a := &Account{}

	req := &createChildAccountRequest{name}

	if err := c.Request(a, http.MethodPatch, apiRoot+"/"+childAccountsPath+"/"+id, req); err != nil {
		return nil, err
	}

	return a, nil
}

func ReadChildAccount(c messagebird.MessageBirdClient, id string) (*Account, error) {
	a := &Account{}

	if err := c.Request(a, http.MethodGet, apiRoot+"/"+childAccountsPath+"/"+id, nil); err != nil {
		return nil, err
	}

	return a, nil
}

// ListChildAccount fetch all the Child Accounts
func ListChildAccount(c messagebird.MessageBirdClient) (*Accounts, error) {
	a := &Accounts{}

	if err := c.Request(a, http.MethodGet, apiRoot+"/"+childAccountsPath, nil); err != nil {
		return nil, err
	}

	return a, nil
}

func DeleteChildAccount(c messagebird.MessageBirdClient, id string) error {
	return c.Request(nil, http.MethodDelete, apiRoot+"/"+childAccountsPath+"/"+id, nil)
}
