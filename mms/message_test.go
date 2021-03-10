package mms

import (
	"net/http"
	"testing"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v6"
	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, "6d9e7100b1f9406c81a3c303c30ccf05", message.ID)
	assert.Equal(t, "https://rest.messagebird.com/mms/6d9e7100b1f9406c81a3c303c30ccf05", message.HRef)
	assert.Equal(t, "mt", message.Direction)
	assert.Equal(t, "TestName", message.Originator)
	assert.Equal(t, "Hello World", message.Body)
	assert.Equal(t, "http://w3.org/1.gif", message.MediaUrls[0])
	assert.Equal(t, "http://w3.org/2.gif", message.MediaUrls[1])
	assert.Equal(t, "TestReference", message.Reference)
	assert.Equal(t, "TestSubject", message.Subject)
	assert.Nil(t, message.ScheduledDatetime)
	assert.Equal(t, "2017-10-20T12:50:28Z", message.CreatedDatetime.Format(time.RFC3339))
	assert.Equal(t, 1, message.Recipients.TotalCount)
	assert.Equal(t, 1, message.Recipients.TotalSentCount)
	assert.Equal(t, int64(31612345678), message.Recipients.Items[0].Recipient)
	assert.Equal(t, "sent", message.Recipients.Items[0].Status)
	assert.Equal(t, "2017-10-20T12:50:28Z", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))

	_, ok := err.(messagebird.ErrorResponse)
	assert.False(t, ok)
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
	assert.True(t, ok)
	assert.Len(t, errorResponse.Errors, 1)
	assert.Equal(t, 2, errorResponse.Errors[0].Code)
	assert.Equal(t, "access_key", errorResponse.Errors[0].Parameter)
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
	assert.EqualError(t, err, "Body or MediaUrls is required")
}
