package hlr

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/v6"
	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func assertHLRObject(t *testing.T, hlr *HLR) {
	assert.Equal(t, "27978c50354a93ca0ca8de6h54340177", hlr.ID)
	assert.Equal(t, "https://rest.messagebird.com/hlr/27978c50354a93ca0ca8de6h54340177", hlr.HRef)
	assert.Equal(t, 31612345678, hlr.MSISDN)
	assert.Equal(t, 20406, hlr.Network)
	assert.Equal(t, "MyReference", hlr.Reference)
	assert.Equal(t, "sent", hlr.Status)

	assert.Equal(t, "2015-01-04T13:14:08Z", hlr.CreatedDatetime.Format(time.RFC3339))
	assert.Equal(t, "2015-01-04T13:14:09Z", hlr.StatusDatetime.Format(time.RFC3339))

}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "hlrObject.json", http.StatusOK)
	client := mbtest.Client(t)

	hlr, err := Read(client, "27978c50354a93ca0ca8de6h54340177")
	assert.NoError(t, err)

	assertHLRObject(t, hlr)
}

func TestRequestDataForHLR(t *testing.T) {
	requestData, err := requestDataForHLR("31612345678", "MyReference")
	assert.NoError(t, err)
	assert.Equal(t, "31612345678", requestData.MSISDN)
	assert.Equal(t, "MyReference", requestData.Reference)
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "hlrObject.json", http.StatusOK)
	client := mbtest.Client(t)

	hlr, err := Create(client, "31612345678", "MyReference")
	assert.NoError(t, err)

	assertHLRObject(t, hlr)
}

func TestHLRError(t *testing.T) {
	mbtest.WillReturnAccessKeyError()
	client := mbtest.Client(t)

	_, err := Read(client, "dummy_hlr_id")

	errorResponse, ok := err.(messagebird.ErrorResponse)
	assert.True(t, ok)
	assert.Len(t, errorResponse.Errors, 1)
	assert.Equal(t, 2, errorResponse.Errors[0].Code)
	assert.Equal(t, "access_key", errorResponse.Errors[0].Parameter)
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "hlrListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	hlrList, err := List(client)
	assert.NoError(t, err)
	assert.Equal(t, 0, hlrList.Offset)
	assert.Equal(t, 20, hlrList.Limit)
	assert.Equal(t, 2, hlrList.Count)
	assert.Equal(t, 2, hlrList.TotalCount)

	for _, hlr := range hlrList.Items {
		assertHLRObject(t, &hlr)
	}
}
