package messagebird

import "net/url"

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
	HLR           *HLR
}

// LookupParams provide additional lookup information.
type LookupParams struct {
	CountryCode string
	Reference   string
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
