package lookup

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/messagebird/go-rest-api/hlr"
	"github.com/messagebird/go-rest-api/internal/messagebirdtest"
)

func TestMain(m *testing.M) {
	messagebirdtest.EnableServer(m)
}

func TestCreate(t *testing.T) {
	messagebirdtest.WillReturnTestdata(t, "lookupObject.json", http.StatusOK)
	client := messagebirdtest.Client(t)

	phoneNumber := "31624971134"
	lookup, err := Create(client, phoneNumber, &Params{CountryCode: "NL"})
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

func checkHLR(t *testing.T, hlr *hlr.HLR) {
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

func TestReadHLR(t *testing.T) {
	messagebirdtest.WillReturnTestdata(t, "lookupHLRObject.json", http.StatusOK)
	client := messagebirdtest.Client(t)

	hlr, err := ReadHLR(client, "31624971134", &Params{CountryCode: "NL"})
	if err != nil {
		t.Fatalf("Didn't expect error while doing the lookup: %s", err)
	}
	checkHLR(t, hlr)
}

func TestRequestDataForLookupHLR(t *testing.T) {
	lookupParams := &Params{
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

func TestCreateHLR(t *testing.T) {
	messagebirdtest.WillReturnTestdata(t, "lookupHLRObject.json", http.StatusCreated)
	client := messagebirdtest.Client(t)

	hlr, err := CreateHLR(client, "31624971134", &Params{CountryCode: "NL", Reference: "reference2000"})
	if err != nil {
		t.Fatalf("Didn't expect error while doing the lookup: %s", err)
	}

	checkHLR(t, hlr)
}
