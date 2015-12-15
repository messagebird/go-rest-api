package messagebird

import "testing"

var lookupObject []byte = []byte(`{
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
        "reference":"yoloswag2000",
        "status":"active",
        "createdDatetime":"2015-12-15T08:19:24+00:00",
        "statusDatetime":"2015-12-15T08:19:25+00:00"
    }
}`)

var lookupHLRObject []byte = []byte(`{
    "id":"6118d3f06566fcd0cdc8962h65065907",
    "network":20416,
    "reference":"yoloswag2000",
    "status":"active",
    "createdDatetime":"2015-12-15T08:19:24+00:00",
    "statusDatetime":"2015-12-15T08:19:25+00:00"
}`)

func TestLookup(t *testing.T) {
	SetServerResponse(200, lookupObject)

	phoneNumber := 31624971134
	lookup, err := mbClient.Lookup(phoneNumber, "NL")
	if err != nil {
		t.Fatalf("Didn't expect error while doing the lookup: %s", err)
	}

	if lookup.Href != "https://rest.messagebird.com/lookup/31624971134" {
		t.Errorf("Unexpected lookup href: %s", lookup.Href)
	}
	if lookup.PhoneNumber != phoneNumber {
		t.Errorf("Unexpected lookup phoneNumber: %s", lookup.PhoneNumber)
	}
	if lookup.Formats != nil {
		if lookup.Formats.International != "+31 6 24971134" {
			t.Errorf("Unexpected International format: %s", lookup.HLR.Reference)
		}
	} else {
		t.Errorf("Unexpected empty Formats object")
	}
	if lookup.HLR != nil {
		if lookup.HLR.Reference != "yoloswag2000" {
			t.Errorf("Unexpected hlr reference: %s", lookup.HLR.Reference)
		}
	} else {
		t.Errorf("Unexpected empty hlr")
	}
}

func checkHLR(t *testing.T, hlr *HLR) {
	if hlr.Id != "6118d3f06566fcd0cdc8962h65065907" {
		t.Errorf("Unexpected hlr id: %s", hlr.Id)
	}
	if hlr.Network != 20416 {
		t.Errorf("Unexpected hlr network: %d", hlr.Network)
	}
	if hlr.Reference != "yoloswag2000" {
		t.Errorf("Unexpected hlr reference: %s", hlr.Reference)
	}
	if hlr.Status != "active" {
		t.Errorf("Unexpected hlr status: %s", hlr.Status)
	}
}

func TestHLRLookup(t *testing.T) {
	SetServerResponse(200, lookupHLRObject)

	hlr, err := mbClient.HLRLookup(31624971134, "NL")
	if err != nil {
		t.Fatalf("Didn't expect error while doing the lookup: %s", err)
	}
	checkHLR(t, hlr)
}

func TestNewHLRLookup(t *testing.T) {
	SetServerResponse(201, lookupHLRObject)

	hlr, err := mbClient.NewHLRLookup(31624971134, "NL", "yoloswag2000")
	if err != nil {
		t.Fatalf("Didn't expect error while doing the lookup: %s", err)
	}

	checkHLR(t, hlr)
}
