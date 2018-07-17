package lookup

import (
	"net/http"
	"net/url"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/hlr"
)

// Formats represents phone number in multiple formats.
type Formats struct {
	E164          string
	International string
	National      string
	Rfc3966       string
}

// Lookup is used to validate and look up a mobile number.
type Lookup struct {
	Href          string
	CountryCode   string
	CountryPrefix int
	PhoneNumber   int64
	Type          string
	Formats       Formats
	HLR           *hlr.HLR
}

// LookupParams provide additional lookup information.
type LookupParams struct {
	CountryCode string
	Reference   string
}

type lookupRequest struct {
	CountryCode string `json:"countryCode,omitempty"`
	Reference   string `json:"reference,omitempty"`
}

// LookupPath represents the path to the Lookup resource.
const LookupPath = "lookup"

// Create performs a new lookup for the specified number.
func Create(c *messagebird.Client, phoneNumber string, params *LookupParams) (*Lookup, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "?" + urlParams.Encode()

	lookup := &Lookup{}
	if err := c.Request(lookup, http.MethodPost, path, nil); err != nil {
		return nil, err
	}

	return lookup, nil
}

// CreateHLR creates a new HLR lookup for the specified number.
func CreateHLR(c *messagebird.Client, phoneNumber string, params *LookupParams) (*hlr.HLR, error) {
	requestData := requestDataForLookup(params)
	path := LookupPath + "/" + phoneNumber + "/" + hlr.HLRPath

	hlr := &hlr.HLR{}
	if err := c.Request(hlr, http.MethodPost, path, requestData); err != nil {
		return nil, err
	}

	return hlr, nil
}

// ReadHLR performs a HLR lookup for the specified number.
func ReadHLR(c *messagebird.Client, phoneNumber string, params *LookupParams) (*hlr.HLR, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "/" + hlr.HLRPath + "?" + urlParams.Encode()

	hlr := &hlr.HLR{}
	if err := c.Request(hlr, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return hlr, nil
}

func requestDataForLookup(params *LookupParams) *lookupRequest {
	request := &lookupRequest{}

	if params == nil {
		return request
	}

	request.CountryCode = params.CountryCode
	request.Reference = params.Reference

	return request
}

func paramsForLookup(params *LookupParams) *url.Values {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams
	}

	if params.CountryCode != "" {
		urlParams.Set("countryCode", params.CountryCode)
	}
	if params.Reference != "" {
		urlParams.Set("reference", params.Reference)
	}

	return urlParams
}
