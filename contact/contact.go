package contact

import (
	"net/http"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v8"
)

// path represents the path to the Contacts resource.
const path = "contacts"

// Contact gets returned by the API.
type Contact struct {
	ID            string
	HRef          string
	MSISDN        int64
	FirstName     string
	LastName      string
	CustomDetails struct {
		Custom1 string
		Custom2 string
		Custom3 string
		Custom4 string
	}
	Groups struct {
		TotalCount int
		HRef       string
	}
	Messages struct {
		TotalCount int
		HRef       string
	}
	CreatedDatetime *time.Time
	UpdatedDatetime *time.Time
}

type Contacts struct {
	Limit, Offset     int
	Count, TotalCount int
	Items             []Contact
}

// Request represents a contact for write operations, e.g. for creating a new
// contact or updating an existing one.
type Request struct {
	MSISDN    string `json:"msisdn,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Custom1   string `json:"custom1,omitempty"`
	Custom2   string `json:"custom2,omitempty"`
	Custom3   string `json:"custom3,omitempty"`
	Custom4   string `json:"custom4,omitempty"`
}

func Create(c *messagebird.Client, contactRequest *Request) (*Contact, error) {
	contact := &Contact{}
	if err := c.Request(contact, http.MethodPost, path, contactRequest); err != nil {
		return nil, err
	}

	return contact, nil
}

// Delete attempts deleting the contact with the provided ID. If nil is returned,
// the resource was deleted successfully.
func Delete(c *messagebird.Client, id string) error {
	return c.Request(nil, http.MethodDelete, path+"/"+id, nil)
}

// List retrieves a paginated list of contacts, based on the options provided.
// It's worth noting DefaultListOptions.
func List(c *messagebird.Client, options *messagebird.CommonPaginationRequest) (*Contacts, error) {
	contactList := &Contacts{}
	if err := c.Request(contactList, http.MethodGet, path+"?"+options.QueryParams(), nil); err != nil {
		return nil, err
	}

	return contactList, nil
}

// Read retrieves the information of an existing contact.
func Read(c *messagebird.Client, id string) (*Contact, error) {
	contact := &Contact{}
	if err := c.Request(contact, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return contact, nil
}

// Update updates the record referenced by id with any values set in contactRequest.
// Do not set any values that should not be updated.
func Update(c *messagebird.Client, id string, contactRequest *Request) (*Contact, error) {
	contact := &Contact{}
	if err := c.Request(contact, http.MethodPatch, path+"/"+id, contactRequest); err != nil {
		return nil, err
	}

	return contact, nil
}
