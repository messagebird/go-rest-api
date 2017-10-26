package messagebird

import (
	"testing"
	"time"
)

var mmsMessageObject = []byte(`{
    "id": "6d9e7100b1f9406c81a3c303c30ccf05",
    "href": "https://rest.messagebird.com/mms/6d9e7100b1f9406c81a3c303c30ccf05",
    "direction": "mt",
    "originator": "TestName",
    "subject": "TestSubject",
    "body": "Hello World",
    "mediaUrls": ["http://w3.org/1.gif", "http://w3.org/2.gif"],
    "reference": "TestReference",
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

func TestNewMMSMessage(t *testing.T) {
	SetServerResponse(200, mmsMessageObject)

	params := &MMSMessageParams{
		Body:              "Hello World",
		MediaUrls:         []string{"http://w3.org/1.gif", "http://w3.org/2.gif"},
		Subject:           "TestSubject",
		Reference:         "TestReference",
		ScheduledDatetime: time.Now(),
	}
	mmsMessage, err := mbClient.NewMMSMessage("TestName", []string{"31612345678"}, params)

	if err != nil {
		t.Fatalf("Didn't expect error while creating a new MMS message: %s", err)
	}
	if mmsMessage.ID != "6d9e7100b1f9406c81a3c303c30ccf05" {
		t.Errorf("Unexpected mmsMessage id: %s", mmsMessage.ID)
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
	if mmsMessage.MediaUrls[0] != "http://w3.org/1.gif" {
		t.Errorf("Unexpected mmsMessage mediaUrl: %s", mmsMessage.MediaUrls[0])
	}
	if mmsMessage.MediaUrls[1] != "http://w3.org/2.gif" {
		t.Errorf("Unexpected mmsMessage mediaUrl: %s", mmsMessage.MediaUrls[1])
	}
	if mmsMessage.Reference != "TestReference" {
		t.Errorf("Unexpected mmsMessage reference: %s", mmsMessage.Reference)
	}
	if mmsMessage.Subject != "TestSubject" {
		t.Errorf("Unexpected mmsMessage reference: %s", mmsMessage.Subject)
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

func TestNewMMSMessageError(t *testing.T) {
	SetServerResponse(405, accessKeyErrorObject)

	params := &MMSMessageParams{
		Body:              "Hello World",
		MediaUrls:         nil,
		Subject:           "",
		Reference:         "",
		ScheduledDatetime: time.Now(),
	}
	mmsMessage, err := mbClient.NewMMSMessage("TestName", []string{"31612345678"}, params)

	if err != ErrResponse {
		t.Fatalf("Expected ErrResponse to be returned, instead I got %s", err)
	}
	if len(mmsMessage.Errors) != 1 {
		t.Fatalf("Unexpected number of errors: %d", len(mmsMessage.Errors))
	}
	if mmsMessage.Errors[0].Code != 2 {
		t.Errorf("Unexpected error code: %d", mmsMessage.Errors[0].Code)
	}
	if mmsMessage.Errors[0].Parameter != "access_key" {
		t.Errorf("Unexpected error parameter %s", mmsMessage.Errors[0].Parameter)
	}
}

func TestNewMMSMessageWithEmptyParams(t *testing.T) {
	params := &MMSMessageParams{
		Body:              "",
		MediaUrls:         nil,
		Subject:           "",
		Reference:         "",
		ScheduledDatetime: time.Now(),
	}
	_, err := mbClient.NewMMSMessage("TestName", []string{"31612345678"}, params)

	if err == nil {
		t.Fatalf("Expected error to be returned, instead I got nil")
	}
	if err.Error() != "Body or MediaUrls is required" {
		t.Errorf("Unexpected error message, I got %s", err)
	}
}
