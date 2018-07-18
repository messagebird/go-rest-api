package voicemessage

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api"

	"github.com/messagebird/go-rest-api/internal/messagebirdtest"
)

var voiceMessageObject = []byte(`{
	"body": "Hello World",
	"createdDatetime": "2015-01-05T16:11:24+00:00",
	"href": "https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
	"id": "430c44a0354aab7ac9553f7a49907463",
	"ifMachine": "continue",
	"language": "en-gb",
	"originator": "MessageBird",
	"recipients": {
			"items": [
					{
							"recipient": 31612345678,
							"status": "calling",
							"statusDatetime": "2015-01-05T16:11:24+00:00"
					}
			],
			"totalCount": 1,
			"totalDeliveredCount": 0,
			"totalDeliveryFailedCount": 0,
			"totalSentCount": 1
	},
	"reference": null,
	"repeat": 1,
	"scheduledDatetime": null,
	"voice": "female"
}`)

var voiceMessageObjectWithParams = []byte(`{
    "body": "Hello World",
    "createdDatetime": "2015-01-05T16:11:24+00:00",
    "href": "https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
    "id": "430c44a0354aab7ac9553f7a49907463",
    "ifMachine": "hangup",
    "language": "en-gb",
    "recipients": {
        "items": [
            {
                "recipient": 31612345678,
                "status": "calling",
                "statusDatetime": "2015-01-05T16:11:24+00:00"
            }
        ],
        "totalCount": 1,
        "totalDeliveredCount": 0,
        "totalDeliveryFailedCount": 0,
        "totalSentCount": 1
    },
    "reference": "MyReference",
    "repeat": 5,
    "scheduledDatetime": null,
    "voice": "male"
}`)

var voiceMessageObjectWithCreatedDatetime = []byte(`{
    "body": "Hello World",
    "createdDatetime": "2015-01-05T16:11:24+00:00",
    "href": "https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
    "id": "430c44a0354aab7ac9553f7a49907463",
    "ifMachine": "continue",
    "language": "en-gb",
    "recipients": {
        "items": [
            {
                "recipient": 31612345678,
                "status": "scheduled",
                "statusDatetime": null
            }
        ],
        "totalCount": 1,
        "totalDeliveredCount": 0,
        "totalDeliveryFailedCount": 0,
        "totalSentCount": 0
    },
    "reference": null,
    "repeat": 1,
    "scheduledDatetime": "2015-01-05T16:12:24+00:00",
    "voice": "female"
}`)

var voiceMessageListObject = []byte(`{
    "count": 2,
    "items": [
        {
            "body": "Hello World",
            "createdDatetime": "2015-01-05T16:11:24+00:00",
            "href": "https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
            "id": "430c44a0354aab7ac9553f7a49907463",
            "ifMachine": "continue",
            "language": "en-gb",
            "originator": "MessageBird",
            "recipients": {
                "items": [
                    {
                        "recipient": 31612345678,
                        "status": "calling",
                        "statusDatetime": "2015-01-05T16:11:24+00:00"
                    }
                ],
                "totalCount": 1,
                "totalDeliveredCount": 0,
                "totalDeliveryFailedCount": 0,
                "totalSentCount": 1
            },
            "reference": null,
            "repeat": 1,
            "scheduledDatetime": null,
            "voice": "female"
        },
        {
            "body": "Hello World",
            "createdDatetime": "2015-01-05T16:11:24+00:00",
            "href": "https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463",
            "id": "430c44a0354aab7ac9553f7a49907463",
            "ifMachine": "continue",
            "language": "en-gb",
            "originator": "MessageBird",
            "recipients": {
                "items": [
                    {
                        "recipient": 31612345678,
                        "status": "calling",
                        "statusDatetime": "2015-01-05T16:11:24+00:00"
                    }
                ],
                "totalCount": 1,
                "totalDeliveredCount": 0,
                "totalDeliveryFailedCount": 0,
                "totalSentCount": 1
            },
            "reference": null,
            "repeat": 1,
            "scheduledDatetime": null,
            "voice": "female"
        }
    ],
    "limit": 20,
    "links": {
        "first": "https://rest.messagebird.com/voicemessages/?offset=0",
        "last": "https://rest.messagebird.com/voicemessages/?offset=0",
        "next": null,
        "previous": null
    },
    "offset": 0,
    "totalCount": 2
}`)

func TestMain(m *testing.M) {
	messagebirdtest.EnableServer(m)
}

func assertVoiceMessageObject(t *testing.T, message *VoiceMessage) {
	if message.ID != "430c44a0354aab7ac9553f7a49907463" {
		t.Errorf("Unexpected voice message id: %s, expected: 430c44a0354aab7ac9553f7a49907463", message.ID)
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

}

func TestCreate(t *testing.T) {
	messagebirdtest.WillReturn(voiceMessageObject, http.StatusOK)
	client := messagebirdtest.Client(t)

	message, err := Create(client, []string{"31612345678"}, "Hello World", nil)

	errorResponse, ok := err.(messagebird.ErrorResponse)
	if ok {
		t.Errorf("Unexpected error returned with voiceMessage %#v", errorResponse)
	}

	assertVoiceMessageObject(t, message)
}

func TestCreateWithParams(t *testing.T) {
	messagebirdtest.WillReturn(voiceMessageObjectWithParams, http.StatusOK)
	client := messagebirdtest.Client(t)

	params := &Params{
		Reference: "MyReference",
		Voice:     "male",
		Repeat:    5,
		IfMachine: "hangup",
	}

	message, err := Create(client, []string{"31612345678"}, "Hello World", params)
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

func TestCreateWithScheduledDatetime(t *testing.T) {
	messagebirdtest.WillReturn(voiceMessageObjectWithCreatedDatetime, http.StatusOK)
	client := messagebirdtest.Client(t)

	scheduledDatetime, _ := time.Parse(time.RFC3339, "2015-01-05T16:12:24+00:00")

	params := &Params{ScheduledDatetime: scheduledDatetime}

	message, err := Create(client, []string{"31612345678"}, "Hello World", params)
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

func TestList(t *testing.T) {
	messagebirdtest.WillReturn(voiceMessageListObject, http.StatusOK)
	client := messagebirdtest.Client(t)

	messageList, err := List(client)
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

	request, err := requestDataForVoiceMessage([]string{"31612345678"}, "MyBody", voiceParams)
	if err != nil {
		t.Fatalf("Didn't expect error while getting request data for voice message: %s", err)
	}

	if request.Recipients[0] != "31612345678" {
		t.Errorf("Unexpected recipient: %s, expected: 31612345678", request.Recipients[0])
	}
	if request.Body != "MyBody" {
		t.Errorf("Unexpected body: %s, expected: MyBody", request.Body)
	}
	if request.Reference != "MyReference" {
		t.Errorf("Unexpected reference: %s, expected: MyReference", request.Reference)
	}
	if request.Language != "en-gb" {
		t.Errorf("Unexpected language: %s, expected: en-gb", request.Language)
	}
	if request.Voice != "male" {
		t.Errorf("Unexpected voice: %s, expected: male", request.Voice)
	}
	if request.Repeat != 2 {
		t.Errorf("Unexpected repeat: %d, expected: 2", request.Repeat)
	}
	if request.IfMachine != "continue" {
		t.Errorf("Unexpected if machine: %s, expected: continue", request.IfMachine)
	}
	if request.ScheduledDatetime != voiceParams.ScheduledDatetime.Format(time.RFC3339) {
		t.Errorf("Unexpected scheduled date time: %s, expected: %s", request.ScheduledDatetime, voiceParams.ScheduledDatetime.Format(time.RFC3339))
	}
}
