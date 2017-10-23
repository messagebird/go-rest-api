package messagebird

import (
	"net/http"
	"testing"
	"time"
)

var voiceMessageObject []byte = []byte(`{
  "id":"430c44a0354aab7ac9553f7a49907463",
  "href":"https:\/\/rest.messagebird.com\/voicemessages\/430c44a0354aab7ac9553f7a49907463",
  "originator":"MessageBird",
  "body":"Hello World",
  "reference":null,
  "language":"en-gb",
  "voice":"female",
  "repeat":1,
  "ifMachine":"continue",
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T16:11:24+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"calling",
        "statusDatetime":"2015-01-05T16:11:24+00:00"
      }
    ]
  }
}`)

func assertVoiceMessageObject(t *testing.T, message *VoiceMessage) {
	if message.Id != "430c44a0354aab7ac9553f7a49907463" {
		t.Errorf("Unexpected voice message id: %s, expected: 430c44a0354aab7ac9553f7a49907463", message.Id)
	}

	if message.HRef != "https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463" {
		t.Errorf("Unexpected voice message href: %s, expected: https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463", message.HRef)
	}

	if message.Originator != "MessageBird" {
		t.Errorf("Unexpected voice message originator: %s, expected: MessageBird", message.Originator)
	}

	if message.Body != "Hello World" {
		t.Errorf("Unexpected voice message body: %s, expected: Hello World", message.Body)
	}

	if message.Reference != "" {
		t.Errorf("Unexpected voice message reference: %s, expected: \"\"", message.Reference)
	}

	if message.Language != "en-gb" {
		t.Errorf("Unexpected voice message language: %s, expected: en-gb", message.Language)
	}

	if message.Voice != "female" {
		t.Errorf("Unexpected voice message voice: %s, expected: female", message.Voice)
	}

	if message.Repeat != 1 {
		t.Errorf("Unexpected voice message repeat: %d, expected: 1", message.Repeat)
	}

	if message.IfMachine != "continue" {
		t.Errorf("Unexpected voice message ifmachine: %s, expected: continue", message.IfMachine)
	}

	if message.ScheduledDatetime != nil {
		t.Errorf("Unexpected voice message scheduled datetime: %s, expected: nil", message.ScheduledDatetime)
	}

	if message.CreatedDatetime == nil || message.CreatedDatetime.Format(time.RFC3339) != "2015-01-05T16:11:24Z" {
		t.Errorf("Unexpected voice message created datetime: %s, expected: 2015-01-05T16:11:24Z", message.CreatedDatetime.Format(time.RFC3339))
	}

	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d, expected: 1", message.Recipients.TotalCount)
	}

	if message.Recipients.TotalSentCount != 1 {
		t.Errorf("Unexpected number of total sent count: %d, expected: 1", message.Recipients.TotalSentCount)
	}

	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected voice message recipient: %d, expected: 31612345678", message.Recipients.Items[0].Recipient)
	}

	if message.Recipients.Items[0].Status != "calling" {
		t.Errorf("Unexpected voice message recipient status: %s, expected: calling", message.Recipients.Items[0].Status)
	}

	if message.Recipients.Items[0].StatusDatetime == nil || message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339) != "2015-01-05T16:11:24Z" {
		t.Errorf("Unexpected datetime status for voice message recipient: %s, expected: 2015-01-05T16:11:24Z", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}

	if len(message.Errors) != 0 {
		t.Errorf("Unexpected number of errors in voice message: %d, expected: 0", len(message.Errors))
	}
}

func TestNewVoiceMessage(t *testing.T) {
	SetServerResponse(http.StatusOK, voiceMessageObject)

	message, err := mbClient.NewVoiceMessage([]string{"31612345678"}, "Hello World", nil)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new voice message: %s", err)
	}

	assertVoiceMessageObject(t, message)
}

var voiceMessageObjectWithParams []byte = []byte(`{
  "id":"430c44a0354aab7ac9553f7a49907463",
  "href":"https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
  "body":"Hello World",
  "reference":"MyReference",
  "language":"en-gb",
  "voice":"male",
  "repeat":5,
  "ifMachine":"hangup",
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T16:11:24+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"calling",
        "statusDatetime":"2015-01-05T16:11:24+00:00"
      }
    ]
  }
}`)

func TestNewVoiceMessageWithParams(t *testing.T) {
	SetServerResponse(http.StatusOK, voiceMessageObjectWithParams)

	params := &VoiceMessageParams{
		Reference: "MyReference",
		Voice:     "male",
		Repeat:    5,
		IfMachine: "hangup",
	}

	message, err := mbClient.NewVoiceMessage([]string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new voice message: %s", err)
	}

	if message.Reference != "MyReference" {
		t.Errorf("Unexpected voice message reference: %s, expected: MyReference", message.Reference)
	}

	if message.Voice != "male" {
		t.Errorf("Unexpected voice message voice: %s, expected: male", message.Voice)
	}

	if message.Repeat != 5 {
		t.Errorf("Unexpected voice message repeat: %d, expected: 5", message.Repeat)
	}

	if message.IfMachine != "hangup" {
		t.Errorf("Unexpected voice message ifmachine: %s, expected: hangup", message.IfMachine)
	}
}

var voiceMessageObjectWithCreatedDatetime []byte = []byte(`{
  "id":"430c44a0354aab7ac9553f7a49907463",
  "href":"https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
  "body":"Hello World",
  "reference":null,
  "language":"en-gb",
  "voice":"female",
  "repeat":1,
  "ifMachine":"continue",
  "scheduledDatetime":"2015-01-05T16:12:24+00:00",
  "createdDatetime":"2015-01-05T16:11:24+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":0,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"scheduled",
        "statusDatetime":null
      }
    ]
  }
}`)

func TestNewVoiceMessageWithScheduledDatetime(t *testing.T) {
	SetServerResponse(http.StatusOK, voiceMessageObjectWithCreatedDatetime)

	scheduledDatetime, _ := time.Parse(time.RFC3339, "2015-01-05T16:12:24+00:00")

	params := &VoiceMessageParams{ScheduledDatetime: scheduledDatetime}
	message, err := mbClient.NewVoiceMessage([]string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new voice message: %s", err)
	}

	if message.ScheduledDatetime.Format(time.RFC3339) != scheduledDatetime.Format(time.RFC3339) {
		t.Errorf("Unexpected scheduled datetime: %s, expected: %s", message.ScheduledDatetime.Format(time.RFC3339), scheduledDatetime.Format(time.RFC3339))
	}

	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d, expected: 1", message.Recipients.TotalCount)
	}

	if message.Recipients.TotalSentCount != 0 {
		t.Errorf("Unexpected number of total sent count: %d, expected: 0", message.Recipients.TotalSentCount)
	}

	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected voice message recipient: %d, expected: 31612345678", message.Recipients.Items[0].Recipient)
	}

	if message.Recipients.Items[0].Status != "scheduled" {
		t.Errorf("Unexpected voice message recipient status: %s, expected: scheduled", message.Recipients.Items[0].Status)
	}
}

var voiceMessageListObject = []byte(`{
  "offset":0,
  "limit":20,
  "count":2,
  "totalCount":2,
  "links":{
    "first":"https://rest.messagebird.com/voicemessages/?offset=0",
    "previous":null,
    "next":null,
    "last":"https://rest.messagebird.com/voicemessages/?offset=0"
  },
  "items":[
    {
      "id":"430c44a0354aab7ac9553f7a49907463",
      "href":"https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
      "originator":"MessageBird",
      "body":"Hello World",
      "reference":null,
      "language":"en-gb",
      "voice":"female",
      "repeat":1,
      "ifMachine":"continue",
      "scheduledDatetime":null,
      "createdDatetime":"2015-01-05T16:11:24+00:00",
      "recipients":{
        "totalCount":1,
        "totalSentCount":1,
        "totalDeliveredCount":0,
        "totalDeliveryFailedCount":0,
        "items":[
          {
            "recipient":31612345678,
            "status":"calling",
            "statusDatetime":"2015-01-05T16:11:24+00:00"
          }
        ]
      }
    },
    {
      "id":"430c44a0354aab7ac9553f7a49907463",
      "href":"https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
      "originator":"MessageBird",
      "body":"Hello World",
      "reference":null,
      "language":"en-gb",
      "voice":"female",
      "repeat":1,
      "ifMachine":"continue",
      "scheduledDatetime":null,
      "createdDatetime":"2015-01-05T16:11:24+00:00",
      "recipients":{
        "totalCount":1,
        "totalSentCount":1,
        "totalDeliveredCount":0,
        "totalDeliveryFailedCount":0,
        "items":[
          {
            "recipient":31612345678,
            "status":"calling",
            "statusDatetime":"2015-01-05T16:11:24+00:00"
          }
        ]
      }
    }
  ]
}`)

func TestVoiceMessageList(t *testing.T) {
	SetServerResponse(http.StatusOK, voiceMessageListObject)

	messageList, err := mbClient.VoiceMessages()
	if err != nil {
		t.Fatalf("Didn't expect an error while requesting VoiceMessages: %s", err)
	}

	if messageList.Offset != 0 {
		t.Errorf("Unexpected result for the VoiceMessages offset: %d, expected: 0", messageList.Offset)
	}
	if messageList.Limit != 20 {
		t.Errorf("Unexpected result for the VoiceMessages limit: %d, expected: 20", messageList.Limit)
	}
	if messageList.Count != 2 {
		t.Errorf("Unexpected result for the VoiceMessages count: %d, expected: 2", messageList.Count)
	}
	if messageList.TotalCount != 2 {
		t.Errorf("Unexpected result for the VoiceMessages total count: %d, expected: 2", messageList.TotalCount)
	}

	for _, message := range messageList.Items {
		assertVoiceMessageObject(t, &message)
	}
}
