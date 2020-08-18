package mms

import (
	"net/http"
	"testing"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v6"
	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "mmsMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	params := &Params{
		Body:              "Hello World",
		MediaUrls:         []string{"http://w3.org/1.gif", "http://w3.org/2.gif"},
		Subject:           "TestSubject",
		Reference:         "TestReference",
		ScheduledDatetime: time.Now(),
	}

	message, err := Create(client, "TestName", []string{"31612345678"}, params)

	if err != nil {
		t.Fatalf("Didn't expect error while creating a new MMS message: %s", err)
	}
	if message.ID != "6d9e7100b1f9406c81a3c303c30ccf05" {
		t.Errorf("Unexpected message id: %s", message.ID)
	}
	if message.HRef != "https://rest.messagebird.com/mms/6d9e7100b1f9406c81a3c303c30ccf05" {
		t.Errorf("Unexpected message href: %s", message.HRef)
	}
	if message.Direction != "mt" {
		t.Errorf("Unexpected message direction: %s", message.Direction)
	}
	if message.Originator != "TestName" {
		t.Errorf("Unexpected message originator: %s", message.Originator)
	}
	if message.Body != "Hello World" {
		t.Errorf("Unexpected message body: %s", message.Body)
	}
	if message.MediaUrls[0] != "http://w3.org/1.gif" {
		t.Errorf("Unexpected message mediaUrl: %s", message.MediaUrls[0])
	}
	if message.MediaUrls[1] != "http://w3.org/2.gif" {
		t.Errorf("Unexpected message mediaUrl: %s", message.MediaUrls[1])
	}
	if message.Reference != "TestReference" {
		t.Errorf("Unexpected message reference: %s", message.Reference)
	}
	if message.Subject != "TestSubject" {
		t.Errorf("Unexpected message reference: %s", message.Subject)
	}
	if message.ScheduledDatetime != nil {
		t.Errorf("Unexpected message scheduled datetime: %s", message.ScheduledDatetime)
	}
	if message.CreatedDatetime == nil || message.CreatedDatetime.Format(time.RFC3339) != "2017-10-20T12:50:28Z" {
		t.Errorf("Unexpected message created datetime: %s", message.CreatedDatetime)
	}
	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d", message.Recipients.TotalCount)
	}
	if message.Recipients.TotalSentCount != 1 {
		t.Errorf("Unexpected number of total sent count: %d", message.Recipients.TotalSentCount)
	}
	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected message recipient: %d", message.Recipients.Items[0].Recipient)
	}
	if message.Recipients.Items[0].Status != "sent" {
		t.Errorf("Unexpected message recipient status: %s", message.Recipients.Items[0].Status)
	}
	if message.Recipients.Items[0].StatusDatetime == nil || message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339) != "2017-10-20T12:50:28Z" {
		t.Errorf("Unexpected datetime status for message recipient: %s", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}

	errorResponse, ok := err.(messagebird.ErrorResponse)
	if ok {
		t.Errorf("Unexpected error returned with message %#v", errorResponse)
	}
}

func TestCreateError(t *testing.T) {
	mbtest.WillReturnAccessKeyError()
	client := mbtest.Client(t)

	params := &Params{
		Body:              "Hello World",
		MediaUrls:         nil,
		Subject:           "",
		Reference:         "",
		ScheduledDatetime: time.Now(),
	}

	_, err := Create(client, "TestName", []string{"31612345678"}, params)

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

func TestCreateWithEmptyParams(t *testing.T) {
	client := mbtest.Client(t)

	params := &Params{
		Body:              "",
		MediaUrls:         nil,
		Subject:           "",
		Reference:         "",
		ScheduledDatetime: time.Now(),
	}

	_, err := Create(client, "TestName", []string{"31612345678"}, params)

	if err == nil {
		t.Fatalf("Expected error to be returned, instead I got nil")
	}
	if err.Error() != "Body or MediaUrls is required" {
		t.Errorf("Unexpected error message, I got %s", err)
	}
}
