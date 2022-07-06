package lookup

import (
	messagebird "github.com/messagebird/go-rest-api/v9"
	"github.com/messagebird/go-rest-api/v9/hlr"
	"net/http"
	"net/url"
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

// Params provide additional lookup information.
type Params struct {
	CountryCode string
	Reference   string
}

func (p *Params) QueryParams() string {
	if p == nil {
		return ""
	}

	query := url.Values{}

	if p.CountryCode != "" {
		query.Set("countryCode", p.CountryCode)
	}
	if p.Reference != "" {
		query.Set("reference", p.Reference)
	}

	return query.Encode()
}

type lookupRequest struct {
	CountryCode string `json:"countryCode,omitempty"`
	Reference   string `json:"reference,omitempty"`
}

// hlrPath represents the path to the HLR resource within the lookup resource.
const hlrPath = "hlr"

// lookupPath represents the path to the Lookup resource.
const lookupPath = "lookup"

// Read performs a new lookup for the specified number.
func Read(c messagebird.Client, phoneNumber string, params *Params) (*Lookup, error) {
	path := lookupPath + "/" + phoneNumber + "?" + params.QueryParams()

	lookup := &Lookup{}
	if err := c.Request(lookup, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return lookup, nil
}

// CreateHLR creates a new HLR lookup for the specified number.
func CreateHLR(c messagebird.Client, phoneNumber string, params *Params) (*hlr.HLR, error) {
	requestData := requestDataForLookup(params)
	path := lookupPath + "/" + phoneNumber + "/" + hlrPath

	val := &hlr.HLR{}
	if err := c.Request(val, http.MethodPost, path, requestData); err != nil {
		return nil, err
	}

	return val, nil
}

// ReadHLR performs a HLR lookup for the specified number.
func ReadHLR(c messagebird.Client, phoneNumber string, params *Params) (*hlr.HLR, error) {
	path := lookupPath + "/" + phoneNumber + "/" + hlrPath + "?" + params.QueryParams()

	val := &hlr.HLR{}
	if err := c.Request(val, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return val, nil
}

func requestDataForLookup(params *Params) *lookupRequest {
	request := &lookupRequest{}

	if params == nil {
		return request
	}

	request.CountryCode = params.CountryCode
	request.Reference = params.Reference

	return request
}
