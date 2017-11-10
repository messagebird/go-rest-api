package messagebird

import (
	"strconv"
	"testing"
)

var lookupObject = []byte(`{
    "href":"https://rest.messagebird.com/lookup/31624971134",
    "countryCode":"NL",
    "countryPrefix":31,
    "phoneNumber":31624971134,
    "type":"mobile",
    "formats":{
        "e164":"+31624971134",
        "international":"+31 6 24971134",
        "national":"06 24971134",
        "rfc3966":"tel:+31-6-24971134"
    },
    "hlr":{
        "id":"6118d3f06566fcd0cdc8962h65065907",
        "network":20416,
        "reference":"referece2000",
        "status":"active",
        "createdDatetime":"2015-12-15T08:19:24+00:00",
        "statusDatetime":"2015-12-15T08:19:25+00:00"
    }
}`)

var lookupHLRObject = []byte(`{
    "id":"6118d3f06566fcd0cdc8962h65065907",
    "network":20416,
    "reference":"referece2000",
    "status":"active",
    "createdDatetime":"2015-12-15T08:19:24+00:00",
    "statusDatetime":"2015-12-15T08:19:25+00:00"
}`)

func TestLookup(t *testing.T) {
	SetServerResponse(200, lookupObject)

	phoneNumber := "31624971134"
	lookup, err := mbClient.Lookup(phoneNumber, &LookupParams{CountryCode: "NL"})
	if err != nil {
		t.Fatalf("Didn't expect error while doing the lookup: %s", err)
	}

	if lookup.Href != "https://rest.messagebird.com/lookup/31624971134" {
		t.Errorf("Unexpected lookup href: %s", lookup.Href)
	}
	if strconv.FormatInt(lookup.PhoneNumber, 10) != phoneNumber {
		t.Errorf("Unexpected lookup phoneNumber: %d", lookup.PhoneNumber)
	}

	if lookup.Formats.International != "+31 6 24971134" {
		t.Errorf("Unexpected International format: %s", lookup.HLR.Reference)
	}

	if lookup.HLR != nil {
		if lookup.HLR.Reference != "referece2000" {
			t.Errorf("Unexpected hlr reference: %s", lookup.HLR.Reference)
		}
	} else {
		t.Errorf("Unexpected empty hlr")
	}
}

func checkHLR(t *testing.T, hlr *HLR) {
	if hlr.ID != "6118d3f06566fcd0cdc8962h65065907" {
		t.Errorf("Unexpected hlr id: %s", hlr.ID)
	}
	if hlr.Network != 20416 {
		t.Errorf("Unexpected hlr network: %d", hlr.Network)
	}
	if hlr.Reference != "referece2000" {
		t.Errorf("Unexpected hlr reference: %s", hlr.Reference)
	}
	if hlr.Status != "active" {
		t.Errorf("Unexpected hlr status: %s", hlr.Status)
	}
}

func TestLookupHLR(t *testing.T) {
	SetServerResponse(200, lookupHLRObject)

	hlr, err := mbClient.LookupHLR("31624971134", &LookupParams{CountryCode: "NL"})
	if err != nil {
		t.Fatalf("Didn't expect error while doing the lookup: %s", err)
	}
	checkHLR(t, hlr)
}

func TestRequestDataForLookupHLR(t *testing.T) {
	lookupParams := &LookupParams{
		CountryCode: "NL",
		Reference:   "MyReference",
	}
	request := requestDataForLookup(lookupParams)

	if request.CountryCode != "NL" {
		t.Errorf("Unexpected country code: %s, expected: NL", request.CountryCode)
	}
	if request.Reference != "MyReference" {
		t.Errorf("Unexpected reference: %s, expected: MyReference", request.Reference)
	}
}

func TestNewLookupHLR(t *testing.T) {
	SetServerResponse(201, lookupHLRObject)

	hlr, err := mbClient.NewLookupHLR("31624971134", &LookupParams{CountryCode: "NL", Reference: "reference2000"})
	if err != nil {
		t.Fatalf("Didn't expect error while doing the lookup: %s", err)
	}

	checkHLR(t, hlr)
}
