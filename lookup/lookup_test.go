package lookup

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/messagebird/go-rest-api/v6/hlr"
	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "lookupObject.json", http.StatusOK)
	client := mbtest.Client(t)

	phoneNumber := "31624971134"
	lookup, err := Read(client, phoneNumber, &Params{CountryCode: "NL"})
	assert.NoError(t, err)
	assert.Equal(t, "https://rest.messagebird.com/lookup/31624971134", lookup.Href)

	assert.Equal(t, phoneNumber, strconv.FormatInt(lookup.PhoneNumber, 10))
	assert.Equal(t, "+31 6 24971134", lookup.Formats.International)

	assert.Equal(t, "referece2000", lookup.HLR.Reference)
}

func checkHLR(t *testing.T, hlr *hlr.HLR) {
	assert.Equal(t, "6118d3f06566fcd0cdc8962h65065907", hlr.ID)
	assert.Equal(t, 20416, hlr.Network)
	assert.Equal(t, "referece2000", hlr.Reference)
	assert.Equal(t, "active", hlr.Status)
}

func TestReadHLR(t *testing.T) {
	mbtest.WillReturnTestdata(t, "lookupHLRObject.json", http.StatusOK)
	client := mbtest.Client(t)

	hlr, err := ReadHLR(client, "31624971134", &Params{CountryCode: "NL"})
	assert.NoError(t, err)
	checkHLR(t, hlr)
}

func TestRequestDataForLookupHLR(t *testing.T) {
	lookupParams := &Params{
		CountryCode: "NL",
		Reference:   "MyReference",
	}
	request := requestDataForLookup(lookupParams)
	assert.Equal(t, "NL", request.CountryCode)
	assert.Equal(t, "MyReference", request.Reference)
}

func TestCreateHLR(t *testing.T) {
	mbtest.WillReturnTestdata(t, "lookupHLRObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	hlr, err := CreateHLR(client, "31624971134", &Params{CountryCode: "NL", Reference: "reference2000"})
	assert.NoError(t, err)

	checkHLR(t, hlr)
}
