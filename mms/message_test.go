package mms

import (
	"net/http"
	"testing"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v9"
	"github.com/messagebird/go-rest-api/v9/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "mmsMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	scheduledDateTime := time.Now()

	req := &CreateRequest{
		Body:              "Hello World",
		MediaUrls:         []string{"https://media.giphy.com/media/Vuw9m5wXviFIQ/giphy.gif", "https://media.giphy.com/media/pxy9QQUMF0glq/giphy.gif"},
		Subject:           "TestSubject",
		Reference:         "TestReference",
		ScheduledDatetime: &scheduledDateTime,
	}

	message, err := Create(client, req)

	assert.NoError(t, err)
	assert.Equal(t, "6d9e7100b1f9406c81a3c303c30ccf05", message.ID)
	assert.Equal(t, "https://rest.messagebird.com/mms/6d9e7100b1f9406c81a3c303c30ccf05", message.HRef)
	assert.Equal(t, "mt", message.Direction)
	assert.Equal(t, "TestName", message.Originator)
	assert.Equal(t, "Hello World", message.Body)
	assert.Equal(t, "https://media.giphy.com/media/Vuw9m5wXviFIQ/giphy.gif", message.MediaUrls[0])
	assert.Equal(t, "https://media.giphy.com/media/pxy9QQUMF0glq/giphy.gif", message.MediaUrls[1])
	assert.Equal(t, "TestReference", message.Reference)
	assert.Equal(t, "TestSubject", message.Subject)
	assert.Nil(t, message.ScheduledDatetime)
	assert.Equal(t, "2022-05-20T12:50:28Z", message.CreatedDatetime.Format(time.RFC3339))
	assert.Equal(t, 1, message.Recipients.TotalCount)
	assert.Equal(t, 1, message.Recipients.TotalSentCount)
	assert.Equal(t, int64(31612345678), message.Recipients.Items[0].Recipient)
	assert.Equal(t, "sent", message.Recipients.Items[0].Status)
	assert.Equal(t, "2022-05-20T12:50:28Z", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))

	_, ok := err.(messagebird.ErrorResponse)
	assert.False(t, ok)
}

func TestCreateNilRequest(t *testing.T) {
	mbtest.WillReturnTestdata(t, "mmsMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	message, err := Create(client, nil)

	assert.Nil(t, message)
	assert.Error(t, err)
	assert.EqualError(t, err, "create request should not be nil")
}

func TestCreateWithoutBodyAndMediaUrls(t *testing.T) {
	mbtest.WillReturnTestdata(t, "mmsMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	scheduledDateTime := time.Now()

	req := &CreateRequest{
		Subject:           "TestSubject",
		Reference:         "TestReference",
		ScheduledDatetime: &scheduledDateTime,
	}

	message, err := Create(client, req)

	assert.Nil(t, message)
	assert.Error(t, err)
	assert.EqualError(t, err, "body or mediaUrls is required")
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "mmsMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	message, err := Read(client, "6d9e7100b1f9406c81a3c303c30ccf05")

	assert.NoError(t, err)
	assert.Equal(t, "6d9e7100b1f9406c81a3c303c30ccf05", message.ID)
	assert.Equal(t, "https://rest.messagebird.com/mms/6d9e7100b1f9406c81a3c303c30ccf05", message.HRef)
	assert.Equal(t, "mt", message.Direction)
	assert.Equal(t, "TestName", message.Originator)
	assert.Equal(t, "Hello World", message.Body)
	assert.Equal(t, "https://media.giphy.com/media/Vuw9m5wXviFIQ/giphy.gif", message.MediaUrls[0])
	assert.Equal(t, "https://media.giphy.com/media/pxy9QQUMF0glq/giphy.gif", message.MediaUrls[1])
	assert.Equal(t, "TestReference", message.Reference)
	assert.Equal(t, "TestSubject", message.Subject)
	assert.Nil(t, message.ScheduledDatetime)
	assert.Equal(t, "2022-05-20T12:50:28Z", message.CreatedDatetime.Format(time.RFC3339))
	assert.Equal(t, 1, message.Recipients.TotalCount)
	assert.Equal(t, 1, message.Recipients.TotalSentCount)
	assert.Equal(t, int64(31612345678), message.Recipients.Items[0].Recipient)
	assert.Equal(t, "sent", message.Recipients.Items[0].Status)
	assert.Equal(t, "2022-05-20T12:50:28Z", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))

	_, ok := err.(messagebird.ErrorResponse)
	assert.False(t, ok)
}
