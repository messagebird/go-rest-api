package verify

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/internal/mbtest"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func assertVerifyObject(t *testing.T, v *Verify) {
	if v == nil {
		t.Error("Got nil but expected valid object")
	}

	if v.ID != "15498233759288aaf929661v21936686" {
		t.Errorf("Unexpected Verify ID: %s", v.ID)
	}

	if v.HRef != "https://rest.messagebird.com/verify/15498233759288aaf929661v21936686" {
		t.Errorf("Unexpected HRef: %s", v.HRef)
	}

	if v.Recipient != 31612345678 {
		t.Errorf("Unexpected Recipient: %d", v.Recipient)
	}

	if v.Reference != "MyReference" {
		t.Errorf("Unexpected Reference: %s", v.Reference)
	}

	if len(v.Messages) != 1 {
		t.Errorf("Got %d messages, expected 1", len(v.Messages))
	}

	if v.Messages["href"] != "https://rest.messagebird.com/messages/c2bbd563759288aaf962910b56023756" {
		t.Errorf("Unexpected HRef value in messages: %s", v.Messages["href"])
	}

	if v.Status != "sent" {
		t.Errorf("Unexpected status: %s", v.Status)
	}

	if v.CreatedDatetime == nil || v.CreatedDatetime.Format(time.RFC3339) != "2017-05-26T20:06:07Z" {
		t.Errorf("Unexpected Verify created datetime: %s", v.CreatedDatetime.Format(time.RFC3339))
	}

	if v.ValidUntilDatetime == nil || v.ValidUntilDatetime.Format(time.RFC3339) != "2017-05-26T20:06:37Z" {
		t.Errorf("Unexpected Verify valid until datetime: %s", v.ValidUntilDatetime.Format(time.RFC3339))
	}
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "verifyObject.json", http.StatusOK)
	client := mbtest.Client(t)

	v, err := Create(client, "31612345678", nil)
	if err != nil {
		t.Fatalf("Didn't expect an error while requesting a verification: %s", err)
	}

	assertVerifyObject(t, v)
}

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	if err := Delete(client, "15498233759288aaf929661v21936686"); err != nil {
		t.Fatalf("unexpected error deleting Verify: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/verify/15498233759288aaf929661v21936686")
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "verifyObject.json", http.StatusOK)
	client := mbtest.Client(t)

	v, err := Read(client, "15498233759288aaf929661v21936686")
	if err != nil {
		t.Fatalf("unexpected error reading Verify: %s", err)
	}

	if v.ID != "15498233759288aaf929661v21936686" {
		t.Fatalf("got %s, expected 15498233759288aaf929661v21936686", v.ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/verify/15498233759288aaf929661v21936686")
}

func TestVerifyToken(t *testing.T) {
	mbtest.WillReturnTestdata(t, "verifyTokenObject.json", http.StatusOK)
	client := mbtest.Client(t)

	v, err := VerifyToken(client, "does not", "matter")
	if err != nil {
		t.Fatalf("Didn't expect an error while verifying token: %s", err)
	}

	assertVerifyTokenObject(t, v)
}

func assertVerifyTokenObject(t *testing.T, v *Verify) {
	if v == nil {
		t.Error("Got nil but expected valid object")
	}

	if v.ID != "a3f2edb23592d68163f9694v13904556" {
		t.Errorf("Unexpected Verify ID: %s", v.ID)
	}

	if v.HRef != "https://rest.messagebird.com/verify/a3f2edb23592d68163f9694v13904556" {
		t.Errorf("Unexpected HRef: %s", v.HRef)
	}

	if v.Recipient != 31612345678 {
		t.Errorf("Unexpected Recipient: %d", v.Recipient)
	}

	if v.Reference != "MyReference" {
		t.Errorf("Unexpected Reference: %s", v.Reference)
	}

	if len(v.Messages) != 1 {
		t.Errorf("Got %d messages, expected 1", len(v.Messages))
	}

	if v.Messages["href"] != "https://rest.messagebird.com/messages/63b168423592d681641eb07b76226648" {
		t.Errorf("Unexpected HRef value in messages: %s", v.Messages["href"])
	}

	if v.Status != "verified" {
		t.Errorf("Unexpected status: %s", v.Status)
	}

	if v.CreatedDatetime == nil || v.CreatedDatetime.Format(time.RFC3339) != "2017-05-30T12:39:50Z" {
		t.Errorf("Unexpected Verify created datetime: %s", v.CreatedDatetime.Format(time.RFC3339))
	}

	if v.ValidUntilDatetime == nil || v.ValidUntilDatetime.Format(time.RFC3339) != "2017-05-30T12:40:20Z" {
		t.Errorf("Unexpected Verify valid until datetime: %s", v.ValidUntilDatetime.Format(time.RFC3339))
	}
}

func TestRequestDataForVerify(t *testing.T) {
	verifyParams := &Params{
		Originator:  "MSGBIRD",
		Reference:   "MyReference",
		Type:        "sms",
		Template:    "Your code is: %token",
		DataCoding:  "plain",
		ReportURL:   "http://example.com/report",
		Voice:       "male",
		Language:    "en-gb",
		Timeout:     20,
		TokenLength: 8,
	}

	requestData, err := requestDataForVerify("31612345678", verifyParams)
	if err != nil {
		t.Fatalf("Didn't expect error while getting request data for message: %s", err)
	}

	if requestData.Recipient != "31612345678" {
		t.Errorf("Unexpected recipient: %s, expected: 31612345678", requestData.Recipient)
	}
	if requestData.Originator != "MSGBIRD" {
		t.Errorf("Unexpected originator: %s, expected: MSGBIRD", requestData.Originator)
	}
	if requestData.Reference != "MyReference" {
		t.Errorf("Unexpected reference: %s, expected: MyReference", requestData.Reference)
	}
	if requestData.Type != "sms" {
		t.Errorf("Unexpected type: %s, expected: sms", requestData.Type)
	}
	if requestData.DataCoding != "plain" {
		t.Errorf("Unexpected data coding: %s, expected: plain", requestData.DataCoding)
	}
	if requestData.ReportURL != "http://example.com/report" {
		t.Errorf("Unexpected report URL: %s, expected: http://example.com/repot", requestData.ReportURL)
	}
	if requestData.Voice != "male" {
		t.Errorf("Unexpected voice: %s, expected: male", requestData.Voice)
	}
	if requestData.Language != "en-gb" {
		t.Errorf("Unexpected language: %s, expected en-gb", requestData.Language)
	}
	if requestData.Timeout != 20 {
		t.Errorf("Unexepcted timeout: %d, expected 20", requestData.Timeout)
	}
	if requestData.TokenLength != 8 {
		t.Errorf("Unexpected token length: %d, expected 8", requestData.TokenLength)
	}
}
