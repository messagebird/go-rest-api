package hlr

import (
	"errors"
	"net/http"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v7"
)

// HLR stands for Home Location Register. Contains information about the
// subscribers identity, telephone number, the associated services and general
// information about the location of the subscriber.
type HLR struct {
	ID              string
	HRef            string
	MSISDN          int
	Network         int
	Reference       string
	Status          string
	Details         map[string]interface{}
	CreatedDatetime *time.Time
	StatusDatetime  *time.Time
}

// HLRList represents a list of HLR requests.
type HLRList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      map[string]*string
	Items      []HLR
}

type hlrRequest struct {
	MSISDN    string `json:"msisdn"`
	Reference string `json:"reference"`
}

// path represents the path to the HLR resource.
const path = "hlr"

// Read looks up an existing HLR object for the specified id that was previously
// created by the NewHLR function.
func Read(c *messagebird.Client, id string) (*HLR, error) {
	hlr := &HLR{}
	if err := c.Request(hlr, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return hlr, nil
}

// List all HLR objects that were previously created by the Create function.
func List(c *messagebird.Client) (*HLRList, error) {
	hlrList := &HLRList{}
	if err := c.Request(hlrList, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return hlrList, nil
}

// Create creates a new HLR object.
func Create(c *messagebird.Client, msisdn string, reference string) (*HLR, error) {
	requestData, err := requestDataForHLR(msisdn, reference)
	if err != nil {
		return nil, err
	}

	hlr := &HLR{}

	if err := c.Request(hlr, http.MethodPost, path, requestData); err != nil {
		return nil, err
	}

	return hlr, nil
}

func requestDataForHLR(msisdn string, reference string) (*hlrRequest, error) {
	if msisdn == "" {
		return nil, errors.New("msisdn is required")
	}
	if reference == "" {
		return nil, errors.New("reference is required")
	}

	request := &hlrRequest{
		MSISDN:    msisdn,
		Reference: reference,
	}

	return request, nil
}
