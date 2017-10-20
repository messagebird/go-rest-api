package messagebird

import (
	"testing"
	"time"
)

var mmsMessageObject []byte = []byte(`{
    "id": "6d9e7100b1f9406c81a3c303c30ccf05",
    "href": "https://rest.messagebird.com/mms/6d9e7100b1f9406c81a3c303c30ccf05",
    "direction": "mt",
    "originator": "TestName",
    "subject": "SBJCT",
    "body": "Hello World",
    "mediaUrls": [],
    "reference": null,
    "scheduledDatetime": null,
    "createdDatetime": "2017-10-20T12:50:28+00:00",
    "recipients": {
        "totalCount": 1,
        "totalSentCount": 1,
        "totalDeliveredCount": 0,
        "totalDeliveryFailedCount": 0,
        "items": [
            {
                "recipient": 31612345678,
                "status": "sent",
                "statusDatetime": "2017-10-20T12:50:28+00:00"
            }
        ]
    }
}`)

func TestNewMmsMessage(t *testing.T) {
	SetServerResponse(200, mmsMessageObject)

	mmsMessage, err := mbClient.NewMmsMessage("TestName", []string{"31612345678"}, "Hello World", nil, nil)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new MMS message: %s", err)
	}

	if mmsMessage.Id != "6d9e7100b1f9406c81a3c303c30ccf05" {
		t.Errorf("Unexpected mmsMessage id: %s", mmsMessage.Id)
	}

	if mmsMessage.HRef != "https://rest.messagebird.com/mms/6d9e7100b1f9406c81a3c303c30ccf05" {
		t.Errorf("Unexpected mmsMessage href: %s", mmsMessage.HRef)
	}

	if mmsMessage.Direction != "mt" {
		t.Errorf("Unexpected mmsMessage direction: %s", mmsMessage.Direction)
	}

	if mmsMessage.Originator != "TestName" {
		t.Errorf("Unexpected mmsMessage originator: %s", mmsMessage.Originator)
	}

	if mmsMessage.Body != "Hello World" {
		t.Errorf("Unexpected mmsMessage body: %s", mmsMessage.Body)
	}

	if mmsMessage.Reference != "" {
		t.Errorf("Unexpected mmsMessage reference: %s", mmsMessage.Reference)
	}

	if mmsMessage.ScheduledDatetime != nil {
		t.Errorf("Unexpected mmsMessage scheduled datetime: %s", mmsMessage.ScheduledDatetime)
	}

	if mmsMessage.CreatedDatetime == nil || mmsMessage.CreatedDatetime.Format(time.RFC3339) != "2017-10-20T12:50:28Z" {
		t.Errorf("Unexpected mmsMessage created datetime: %s", mmsMessage.CreatedDatetime)
	}

	if mmsMessage.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d", mmsMessage.Recipients.TotalCount)
	}

	if mmsMessage.Recipients.TotalSentCount != 1 {
		t.Errorf("Unexpected number of total sent count: %d", mmsMessage.Recipients.TotalSentCount)
	}

	if mmsMessage.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected mmsMessage recipient: %d", mmsMessage.Recipients.Items[0].Recipient)
	}

	if mmsMessage.Recipients.Items[0].Status != "sent" {
		t.Errorf("Unexpected mmsMessage recipient status: %s", mmsMessage.Recipients.Items[0].Status)
	}

	if mmsMessage.Recipients.Items[0].StatusDatetime == nil || mmsMessage.Recipients.Items[0].StatusDatetime.Format(time.RFC3339) != "2017-10-20T12:50:28Z" {
		t.Errorf("Unexpected datetime status for mmsMessage recipient: %s", mmsMessage.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}

	if len(mmsMessage.Errors) != 0 {
		t.Errorf("Unexpected number of errors in mmsMessage: %d", len(mmsMessage.Errors))
	}
}

func TestNewMmsMessageError(t *testing.T) {
	SetServerResponse(405, accessKeyErrorObject)

	message, err := mbClient.NewMmsMessage("TestName", []string{"31612345678"}, "Hello World", nil, nil)
	if err != ErrResponse {
		t.Fatalf("Expected ErrResponse to be returned, instead I got %s", err)
	}

	if len(message.Errors) != 1 {
		t.Fatalf("Unexpected number of errors: %d", len(message.Errors))
	}

	if message.Errors[0].Code != 2 {
		t.Errorf("Unexpected error code: %d", message.Errors[0].Code)
	}

	if message.Errors[0].Parameter != "access_key" {
		t.Errorf("Unexpected error parameter %s", message.Errors[0].Parameter)
	}
}

func TestNewMmsMessageWithParams(t *testing.T) {
	SetServerResponse(200, mmsMessageObjectWithParams)

	params := &MmsMessageParams{
		Subject:   "Test-Subject",
		Reference: "Test-Reference",
	}

	mmsMessage, err := mbClient.NewMmsMessage("TestName", []string{"31612345678"}, "", []string{"http://w3.org/1.gif", "http://w3.org/2.gif"}, params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new MMS message: %s", err)
	}

	if mmsMessage.Subject != "Test-Subject" {
		t.Errorf("Unexpected message subject: %s", mmsMessage.Subject)
	}
	if mmsMessage.Reference != "Test-Reference" {
		t.Errorf("Unexpected message reference: %s", mmsMessage.Reference)
	}
}

var mmsMessageObjectWithParams []byte = []byte(`{
    "id": "6d9e7100b1f9406c81a3c303c30ccf05",
    "href": "https://rest.messagebird.com/mms/6d9e7100b1f9406c81a3c303c30ccf05",
    "direction": "mt",
    "originator": "TestName",
    "subject": "Test-Subject",
    "body": "Hello World",
    "mediaUrls": [],
    "reference": "Test-Reference",
    "scheduledDatetime": null,
    "createdDatetime": "2017-10-20T12:50:28+00:00",
    "recipients": {
        "totalCount": 1,
        "totalSentCount": 1,
        "totalDeliveredCount": 0,
        "totalDeliveryFailedCount": 0,
        "items": [
            {
                "recipient": 31612345678,
                "status": "sent",
                "statusDatetime": "2017-10-20T12:50:28+00:00"
            }
        ]
    }
}`)
