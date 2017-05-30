package messagebird

import (
	"testing"
	"time"
)

var verifyObject = []byte(`{
  "id": "15498233759288aaf929661v21936686",
  "href": "https://rest.messagebird.com/verify/15498233759288aaf929661v21936686",
  "recipient": 31612345678,
  "reference": "MyReference",
  "messages": {
    "href": "https://rest.messagebird.com/messages/c2bbd563759288aaf962910b56023756"
  },
  "status": "sent",
  "createdDatetime": "2017-05-26T20:06:07+00:00",
  "validUntilDatetime": "2017-05-26T20:06:37+00:00"
}`)

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
		t.Errorf("Unexpected Recipient: %s", v.Recipient)
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

func TestVerify(t *testing.T) {
	SetServerResponse(200, verifyObject)

	v, err := mbClient.NewVerify("31612345678", nil)
	if err != nil {
		t.Fatalf("Didn't expect an error while requesting a verification: %s", err)
	}

	assertVerifyObject(t, v)
}

var verifyTokenObject = []byte(`{
  "id": "a3f2edb23592d68163f9694v13904556",
  "href": "https://rest.messagebird.com/verify/a3f2edb23592d68163f9694v13904556",
  "recipient": 31612345678,
  "reference": "MyReference",
  "messages": {
    "href": "https://rest.messagebird.com/messages/63b168423592d681641eb07b76226648"
  },
  "status": "verified",
  "createdDatetime": "2017-05-30T12:39:50+00:00",
  "validUntilDatetime": "2017-05-30T12:40:20+00:00"
}`)

func TestVerifyToken(t *testing.T) {
	SetServerResponse(200, verifyTokenObject)

	v, err := mbClient.VerifyToken("does not", "matter")
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
		t.Errorf("Unexpected Recipient: %s", v.Recipient)
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
