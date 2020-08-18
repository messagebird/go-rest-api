package hlr

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/v6"
	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func assertHLRObject(t *testing.T, hlr *HLR) {
	if hlr.ID != "27978c50354a93ca0ca8de6h54340177" {
		t.Errorf("Unexpected result for HLR Id: %s, expected: 27978c50354a93ca0ca8de6h54340177", hlr.ID)
	}

	if hlr.HRef != "https://rest.messagebird.com/hlr/27978c50354a93ca0ca8de6h54340177" {
		t.Errorf("Unexpected HLR href: %s, expected: https://rest.messagebird.com/hlr/27978c50354a93ca0ca8de6h54340177", hlr.HRef)
	}

	if hlr.MSISDN != 31612345678 {
		t.Errorf("Unexpected HLR msisdn: %d, expected: 31612345678", hlr.MSISDN)
	}

	if hlr.Network != 20406 {
		t.Errorf("Unexpected HLR network: %d, expected: 20406", hlr.Network)
	}

	if hlr.Reference != "MyReference" {
		t.Errorf("Unexpected HLR reference: %s, expected: MyReference", hlr.Reference)
	}

	if hlr.Status != "sent" {
		t.Errorf("Unexpected HLR status: %s, expected: sent", hlr.Status)
	}

	if hlr.CreatedDatetime == nil || hlr.CreatedDatetime.Format(time.RFC3339) != "2015-01-04T13:14:08Z" {
		t.Errorf("Unexpected HLR created datetime: %s, expected: 2015-01-04T13:14:08Z", hlr.CreatedDatetime.Format(time.RFC3339))
	}

	if hlr.StatusDatetime == nil || hlr.StatusDatetime.Format(time.RFC3339) != "2015-01-04T13:14:09Z" {
		t.Errorf("Unexpected HLR status datetime: %s, expected: 2015-01-04T13:14:09Z", hlr.StatusDatetime.Format(time.RFC3339))
	}
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "hlrObject.json", http.StatusOK)
	client := mbtest.Client(t)

	hlr, err := Read(client, "27978c50354a93ca0ca8de6h54340177")
	if err != nil {
		t.Fatalf("Didn't expect an error while requesting a HLR: %s", err)
	}

	assertHLRObject(t, hlr)
}

func TestRequestDataForHLR(t *testing.T) {
	requestData, err := requestDataForHLR("31612345678", "MyReference")
	if err != nil {
		t.Fatalf("Didn't expect an error while getting the request data for a HLR: %s", err)
	}

	if requestData.MSISDN != "31612345678" {
		t.Errorf("Unexpected msisdn: %s, expected: 31612345678", requestData.MSISDN)
	}

	if requestData.Reference != "MyReference" {
		t.Errorf("Unexpected reference: %s, expected: MyReference", requestData.Reference)
	}
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "hlrObject.json", http.StatusOK)
	client := mbtest.Client(t)

	hlr, err := Create(client, "31612345678", "MyReference")
	if err != nil {
		t.Fatalf("Didn't expect an error while creating a new HLR: %s", err)
	}

	assertHLRObject(t, hlr)
}

func TestHLRError(t *testing.T) {
	mbtest.WillReturnAccessKeyError()
	client := mbtest.Client(t)

	_, err := Read(client, "dummy_hlr_id")

	errorResponse, ok := err.(messagebird.ErrorResponse)
	if !ok {
		t.Fatalf("Expected ErrorResponse to be returned, instead I got %s", err)
	}

	if len(errorResponse.Errors) != 1 {
		t.Fatalf("Unexpected number of errors: %d, expected: 1", len(errorResponse.Errors))
	}

	if errorResponse.Errors[0].Code != 2 {
		t.Errorf("Unexpected error code: %d, expected: 2", errorResponse.Errors[0].Code)
	}

	if errorResponse.Errors[0].Parameter != "access_key" {
		t.Errorf("Unexpected error parameter: %s, expected: access_key", errorResponse.Errors[0].Parameter)
	}
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "hlrListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	hlrList, err := List(client)
	if err != nil {
		t.Fatalf("Didn't expect an error while requesting HLRs: %s", err)
	}

	if hlrList.Offset != 0 {
		t.Errorf("Unexpected result for the HLRList offset: %d, expected: 0", hlrList.Offset)
	}
	if hlrList.Limit != 20 {
		t.Errorf("Unexpected result for the HLRList limit: %d, expected: 20", hlrList.Limit)
	}
	if hlrList.Count != 2 {
		t.Errorf("Unexpected result for the HLRList count: %d, expected: 2", hlrList.Count)
	}
	if hlrList.TotalCount != 2 {
		t.Errorf("Unexpected result for the HLRList total count: %d, expected: 2", hlrList.TotalCount)
	}

	for _, hlr := range hlrList.Items {
		assertHLRObject(t, &hlr)
	}
}
