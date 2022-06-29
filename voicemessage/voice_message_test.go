package voicemessage

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/v8"
	"github.com/messagebird/go-rest-api/v8/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func assertVoiceMessageObject(t *testing.T, message *VoiceMessage) {
	assert.Equal(t, "430c44a0354aab7ac9553f7a49907463", message.ID)
	assert.Equal(t, "https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463", message.HRef)
	assert.Equal(t, "MessageBird", message.Originator)

	assert.Equal(t, "Hello World", message.Body)
	assert.Equal(t, "", message.Reference)
	assert.Equal(t, "en-gb", message.Language)
	assert.Equal(t, "female", message.Voice)
	assert.Equal(t, 1, message.Repeat)
	assert.Equal(t, "continue", message.IfMachine)
	assert.Nil(t, message.ScheduledDatetime)

	assert.Equal(t, "2015-01-05T16:11:24Z", message.CreatedDatetime.Format(time.RFC3339))
	assert.Equal(t, 1, message.Recipients.TotalCount)
	assert.Equal(t, 1, message.Recipients.TotalSentCount)
	assert.Equal(t, int64(31612345678), message.Recipients.Items[0].Recipient)
	assert.Equal(t, "calling", message.Recipients.Items[0].Status)

	assert.Equal(t, "2015-01-05T16:11:24Z", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))

}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "voiceMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	message, err := Create(client, []string{"31612345678"}, "Hello World", nil)

	_, ok := err.(messagebird.ErrorResponse)
	assert.False(t, ok)

	assertVoiceMessageObject(t, message)
}

func TestCreateWithParams(t *testing.T) {
	mbtest.WillReturnTestdata(t, "voiceMessageObjectWithParams.json", http.StatusOK)
	client := mbtest.Client(t)

	params := &Params{
		Reference: "MyReference",
		Voice:     "male",
		Repeat:    5,
		IfMachine: "hangup",
	}

	message, err := Create(client, []string{"31612345678"}, "Hello World", params)
	assert.NoError(t, err)
	assert.Equal(t, "MyReference", message.Reference)
	assert.Equal(t, "male", message.Voice)
	assert.Equal(t, 5, message.Repeat)
	assert.Equal(t, "hangup", message.IfMachine)
}

func TestCreateWithScheduledDatetime(t *testing.T) {
	mbtest.WillReturnTestdata(t, "voiceMessageObjectWithCreatedDatetime.json", http.StatusOK)
	client := mbtest.Client(t)

	scheduledDatetime, _ := time.Parse(time.RFC3339, "2015-01-05T16:12:24+00:00")

	params := &Params{ScheduledDatetime: scheduledDatetime}

	message, err := Create(client, []string{"31612345678"}, "Hello World", params)
	assert.NoError(t, err)
	assert.Equal(t, scheduledDatetime.Format(time.RFC3339), message.ScheduledDatetime.Format(time.RFC3339))
	assert.Equal(t, 1, message.Recipients.TotalCount)
	assert.Equal(t, 0, message.Recipients.TotalSentCount)
	assert.Equal(t, int64(31612345678), message.Recipients.Items[0].Recipient)
	assert.Equal(t, "scheduled", message.Recipients.Items[0].Status)
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "voiceMessageListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	messageList, err := List(client)
	assert.NoError(t, err)
	assert.Equal(t, 0, messageList.Offset)
	assert.Equal(t, 20, messageList.Limit)
	assert.Equal(t, 2, messageList.Count)
	assert.Equal(t, 2, messageList.TotalCount)

	for _, message := range messageList.Items {
		assertVoiceMessageObject(t, &message)
	}
}

func TestRequestDataForVoiceMessage(t *testing.T) {
	currentTime := time.Now()
	voiceParams := &Params{
		Originator:        "MSGBIRD",
		Reference:         "MyReference",
		Language:          "en-gb",
		Voice:             "male",
		Repeat:            2,
		IfMachine:         "continue",
		ScheduledDatetime: currentTime,
	}

	request, err := paramsToRequest([]string{"31612345678"}, "MyBody", voiceParams)
	assert.NoError(t, err)
	assert.Equal(t, "31612345678", request.Recipients[0])
	assert.Equal(t, "MyBody", request.Body)
	assert.Equal(t, "MyReference", request.Reference)
	assert.Equal(t, "en-gb", request.Language)
	assert.Equal(t, "male", request.Voice)
	assert.Equal(t, 2, request.Repeat)
	assert.Equal(t, "continue", request.IfMachine)
	assert.Equal(t, voiceParams.ScheduledDatetime.Format(time.RFC3339), request.ScheduledDatetime)
}
