package messagebird

import (
	"net/http"
	"testing"
	"time"
)

var messageObject = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"sms",
  "originator":"TestName",
  "body":"Hello World",
  "reference":null,
  "validity":null,
  "gateway":239,
  "typeDetails":{

  },
  "datacoding":"plain",
  "mclass":1,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func assertMessageObject(t *testing.T, message *Message) {
	if message.ID != "6fe65f90454aa61536e6a88b88972670" {
		t.Errorf("Unexpected message id: %s, expected: 6fe65f90454aa61536e6a88b88972670", message.ID)
	}

	if message.HRef != "https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670" {
		t.Errorf("Unexpected message href: %s, expected: https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670", message.HRef)
	}

	if message.Direction != "mt" {
		t.Errorf("Unexpected message direction: %s, expected: mt", message.Direction)
	}

	if message.Type != "sms" {
		t.Errorf("Unexpected message type: %s, expected: sms", message.Type)
	}

	if message.Originator != "TestName" {
		t.Errorf("Unexpected message originator: %s, expected: TestName", message.Originator)
	}

	if message.Body != "Hello World" {
		t.Errorf("Unexpected message body: %s, expected: Hello World", message.Body)
	}

	if message.Reference != "" {
		t.Errorf("Unexpected message reference: %s, expected: \"\"", message.Reference)
	}

	if message.Validity != nil {
		t.Errorf("Unexpected message validity: %d, expected: nil", *message.Validity)
	}

	if message.Gateway != 239 {
		t.Errorf("Unexpected message gateway: %d, expected: 239", message.Gateway)
	}

	if len(message.TypeDetails) != 0 {
		t.Errorf("Unexpected number of items in message typedetails: %d, expected: 0", len(message.TypeDetails))
	}

	if message.DataCoding != "plain" {
		t.Errorf("Unexpected message datacoding: %s, expected: plain", message.DataCoding)
	}

	if message.MClass != 1 {
		t.Errorf("Unexpected message mclass: %d, expected: 1", message.MClass)
	}

	if message.ScheduledDatetime != nil {
		t.Errorf("Unexpected message scheduled datetime: %s, expected: nil", message.ScheduledDatetime)
	}

	if message.CreatedDatetime == nil || message.CreatedDatetime.Format(time.RFC3339) != "2015-01-05T10:02:59Z" {
		t.Errorf("Unexpected message created datetime: %s, expected: 2015-01-05T10:02:59Z", message.CreatedDatetime)
	}

	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d, expected: 1", message.Recipients.TotalCount)
	}

	if message.Recipients.TotalSentCount != 1 {
		t.Errorf("Unexpected number of total sent count: %d, expected: 1", message.Recipients.TotalSentCount)
	}

	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected message recipient: %d, expected: 31612345678", message.Recipients.Items[0].Recipient)
	}

	if message.Recipients.Items[0].Status != "sent" {
		t.Errorf("Unexpected message recipient status: %s, expected: sent", message.Recipients.Items[0].Status)
	}

	if message.Recipients.Items[0].StatusDatetime == nil || message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339) != "2015-01-05T10:02:59Z" {
		t.Errorf("Unexpected datetime status for message recipient: %s, expected: 2015-01-05T10:02:59Z", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}
}

func TestNewMessage(t *testing.T) {
	SetServerResponse(http.StatusOK, messageObject)

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", nil)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	assertMessageObject(t, message)
}

func TestNewMessageError(t *testing.T) {
	SetServerResponse(http.StatusMethodNotAllowed, accessKeyErrorObject)
	_, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", nil)

	errorResponse, ok := err.(ErrorResponse)
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

var messageWithParamsObject = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"sms",
  "originator":"TestName",
  "body":"Hello World",
  "reference":"TestReference",
  "validity":13,
  "gateway":10,
  "typeDetails":{

  },
  "datacoding":"unicode",
  "mclass":1,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessageWithParams(t *testing.T) {
	SetServerResponse(http.StatusOK, messageWithParamsObject)

	params := &MessageParams{
		Type:       "sms",
		Reference:  "TestReference",
		Validity:   13,
		Gateway:    10,
		DataCoding: "unicode",
	}

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Type != "sms" {
		t.Errorf("Unexpected message type: %s, expected: sms", message.Type)
	}

	if message.Reference != "TestReference" {
		t.Errorf("Unexpected message reference: %s, expected: TestReference", message.Reference)
	}

	if *message.Validity != 13 {
		t.Errorf("Unexpected message validity: %d, expected: 13", *message.Validity)
	}

	if message.Gateway != 10 {
		t.Errorf("Unexpected message gateway: %d, expected: 10", message.Gateway)
	}

	if message.DataCoding != "unicode" {
		t.Errorf("Unexpected message datacoding: %s, expected: unicode", message.DataCoding)
	}
}

var binaryMessageObject = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"binary",
  "originator":"TestName",
  "body":"Hello World",
  "reference":"TestReference",
  "validity":13,
  "gateway":10,
  "typeDetails":{
    "udh":"050003340201"
  },
  "datacoding":"unicode",
  "mclass":1,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessageWithBinaryType(t *testing.T) {
	SetServerResponse(http.StatusOK, binaryMessageObject)

	params := &MessageParams{
		Type:        "binary",
		TypeDetails: TypeDetails{"udh": "050003340201"},
	}

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Type != "binary" {
		t.Errorf("Unexpected message type: %s, expected: binary", message.Type)
	}

	if len(message.TypeDetails) != 1 {
		t.Fatalf("Unexpected number of message typedetails: %d, expected: 1", len(message.TypeDetails))
	}

	if message.TypeDetails["udh"] != "050003340201" {
		t.Errorf("Unexpected 'udh' value in message typedetails: %s, expected: 050003340201", message.TypeDetails["udh"])
	}
}

var premiumMessageObject = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"premium",
  "originator":"TestName",
  "body":"Hello World",
  "reference":"TestReference",
  "validity":13,
  "gateway":10,
  "typeDetails":{
    "tariff":150,
    "shortcode":1008,
    "keyword":"RESTAPI"
  },
  "datacoding":"unicode",
  "mclass":1,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessageWithPremiumType(t *testing.T) {
	SetServerResponse(http.StatusOK, premiumMessageObject)

	params := &MessageParams{
		Type:        "premium",
		TypeDetails: TypeDetails{"keyword": "RESTAPI", "shortcode": 1008, "tariff": 150},
	}

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Type != "premium" {
		t.Errorf("Unexpected message type: %s, expected: premium", message.Type)
	}

	if len(message.TypeDetails) != 3 {
		t.Fatalf("Unexpected number of message typedetails: %d, expected: 3", len(message.TypeDetails))
	}

	if message.TypeDetails["tariff"] != 150.0 {
		t.Errorf("Unexpected 'tariff' value in message typedetails: %d, expected: 150.0", message.TypeDetails["tariff"])
	}

	if message.TypeDetails["shortcode"] != 1008.0 {
		t.Errorf("Unexpected 'shortcode' value in message typedetails: %d, expected: 1008.0", message.TypeDetails["shortcode"])
	}

	if message.TypeDetails["keyword"] != "RESTAPI" {
		t.Errorf("Unexpected 'keyword' value in message typedetails: %s, expected: RESTAPI", message.TypeDetails["keyword"])
	}
}

var flashMessageObject = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"flash",
  "originator":"TestName",
  "body":"Hello World",
  "reference":"TestReference",
  "validity":13,
  "gateway":10,
  "typeDetails":{

  },
  "datacoding":"unicode",
  "mclass":0,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessageWithFlashType(t *testing.T) {
	SetServerResponse(http.StatusOK, flashMessageObject)

	params := &MessageParams{Type: "flash"}

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Type != "flash" {
		t.Errorf("Unexpected message type: %s, expected: flash", message.Type)
	}
}

var messageObjectWithCreatedDatetime = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"sms",
  "originator":"TestName",
  "body":"Hello World",
  "reference":null,
  "validity":null,
  "gateway":239,
  "typeDetails":{

  },
  "datacoding":"plain",
  "mclass":1,
  "scheduledDatetime":"2015-01-05T10:03:59+00:00",
  "createdDatetime":"2015-01-05T10:02:59+00:00",
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

func TestNewMessageWithScheduledDatetime(t *testing.T) {
	SetServerResponse(http.StatusOK, messageObjectWithCreatedDatetime)

	scheduledDatetime, _ := time.Parse(time.RFC3339, "2015-01-05T10:03:59+00:00")

	params := &MessageParams{ScheduledDatetime: scheduledDatetime}
	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.ScheduledDatetime.Format(time.RFC3339) != scheduledDatetime.Format(time.RFC3339) {
		t.Errorf("Unexpected message scheduled datetime: %s, expected: %s", message.ScheduledDatetime.Format(time.RFC3339), scheduledDatetime.Format(time.RFC3339))
	}

	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d, expected: 1", message.Recipients.TotalCount)
	}

	if message.Recipients.TotalSentCount != 0 {
		t.Errorf("Unexpected number of total sent count: %d, expected: 0", message.Recipients.TotalSentCount)
	}

	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected message recipient: %d, expected: 31612345678", message.Recipients.Items[0].Recipient)
	}

	if message.Recipients.Items[0].Status != "scheduled" {
		t.Errorf("Unexpected message recipient status: %s, expected: scheduled", message.Recipients.Items[0].Status)
	}

	if message.Recipients.Items[0].StatusDatetime != nil {
		t.Errorf("Unexpected datetime status for message recipient: %s, expected: nil", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}
}

var messageListObject = []byte(`{
  "offset":0,
  "limit":20,
  "count":2,
  "totalCount":2,
  "links":{
    "first":"https://rest.messagebird.com/messages/?offset=0",
    "previous":null,
    "next":null,
    "last":"https://rest.messagebird.com/messages/?offset=0"
  },
  "items":[
    {
      "id":"6fe65f90454aa61536e6a88b88972670",
      "href":"https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670",
      "direction":"mt",
      "type":"sms",
      "originator":"TestName",
      "body":"Hello World",
      "reference":null,
      "validity":null,
      "gateway":239,
      "typeDetails":{

      },
      "datacoding":"plain",
      "mclass":1,
      "scheduledDatetime":null,
      "createdDatetime":"2015-01-05T10:02:59+00:00",
      "recipients":{
        "totalCount":1,
        "totalSentCount":1,
        "totalDeliveredCount":0,
        "totalDeliveryFailedCount":0,
        "items":[
          {
            "recipient":31612345678,
            "status":"sent",
            "statusDatetime":"2015-01-05T10:02:59+00:00"
          }
        ]
      }
    },
    {
      "id":"6fe65f90454aa61536e6a88b88972670",
      "href":"https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670",
      "direction":"mt",
      "type":"sms",
      "originator":"TestName",
      "body":"Hello World",
      "reference":null,
      "validity":null,
      "gateway":239,
      "typeDetails":{

      },
      "datacoding":"plain",
      "mclass":1,
      "scheduledDatetime":null,
      "createdDatetime":"2015-01-05T10:02:59+00:00",
      "recipients":{
        "totalCount":1,
        "totalSentCount":1,
        "totalDeliveredCount":0,
        "totalDeliveryFailedCount":0,
        "items":[
          {
            "recipient":31612345678,
            "status":"sent",
            "statusDatetime":"2015-01-05T10:02:59+00:00"
          }
        ]
      }
    }
  ]
}`)

func TestMessageList(t *testing.T) {
	SetServerResponse(http.StatusOK, messageListObject)

	messageList, err := mbClient.Messages(nil)
	if err != nil {
		t.Fatalf("Didn't expect an error while requesting Messages: %s", err)
	}

	if messageList.Offset != 0 {
		t.Errorf("Unexpected result for the MessageList offset: %d, expected: 0", messageList.Offset)
	}
	if messageList.Limit != 20 {
		t.Errorf("Unexpected result for the MessageList limit: %d, expected: 20", messageList.Limit)
	}
	if messageList.Count != 2 {
		t.Errorf("Unexpected result for the MessageList count: %d, expected: 2", messageList.Count)
	}
	if messageList.TotalCount != 2 {
		t.Errorf("Unexpected result for the MessageList total count: %d, expected: 2", messageList.TotalCount)
	}

	for _, message := range messageList.Items {
		assertMessageObject(t, &message)
	}
}

func TestRequestDataForMessage(t *testing.T) {
	currentTime := time.Now()
	messageParams := &MessageParams{
		Type:              "binary",
		Reference:         "MyReference",
		Validity:          1,
		Gateway:           0,
		TypeDetails:       nil,
		DataCoding:        "unicode",
		ScheduledDatetime: currentTime,
	}

	request, err := requestDataForMessage("MSGBIRD", []string{"31612345678"}, "MyBody", messageParams)
	if err != nil {
		t.Fatalf("Didn't expect error while getting request data for message: %s", err)
	}

	if request.Originator != "MSGBIRD" {
		t.Errorf("Unexpected originator: %s, expected: MSGBIRD", request.Originator)
	}
	if request.Recipients[0] != "31612345678" {
		t.Errorf("Unexpected recipient: %s, expected: 31612345678", request.Recipients[0])
	}
	if request.Body != "MyBody" {
		t.Errorf("Unexpected body: %s, expected: MyBody", request.Body)
	}
	if request.Type != messageParams.Type {
		t.Errorf("Unexpected type: %s, expected: %s", request.Type, messageParams.Type)
	}
	if request.Reference != messageParams.Reference {
		t.Errorf("Unexpected reference: %s, expected: %s", request.Reference, messageParams.Reference)
	}
	if request.Validity != messageParams.Validity {
		t.Errorf("Unexpected validity: %d, expected: %d", request.Validity, messageParams.Validity)
	}
	if request.Gateway != messageParams.Gateway {
		t.Errorf("Unexpected gateway: %d, expected: %d", request.Gateway, messageParams.Gateway)
	}
	if request.TypeDetails != nil {
		t.Errorf("Unexpected type details: %s, expected: nil", request.TypeDetails)
	}
	if request.DataCoding != messageParams.DataCoding {
		t.Errorf("Unexpected data coding: %s, expected: %s", request.DataCoding, messageParams.DataCoding)
	}
	if request.ScheduledDatetime != messageParams.ScheduledDatetime.Format(time.RFC3339) {
		t.Errorf("Unexepcted scheduled date time: %s, expected: %s", request.ScheduledDatetime, messageParams.ScheduledDatetime.Format(time.RFC3339))
	}

}
