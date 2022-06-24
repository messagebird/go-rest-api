package sms

import (
	"net/http"
	"strings"
	"testing"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v7"
	"github.com/messagebird/go-rest-api/v7/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func assertMessageObject(t *testing.T, message *Message, expectedStatus string) {
	assert.Equal(t, "6fe65f90454aa61536e6a88b88972670", message.ID)
	assert.Equal(t, "https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670", message.HRef)
	assert.Equal(t, "mt", message.Direction)
	assert.Equal(t, "sms", message.Type)
	assert.Equal(t, "TestName", message.Originator)
	assert.Equal(t, "Hello World", message.Body)
	assert.Equal(t, "", message.Reference)
	assert.Nil(t, message.Validity)
	assert.Equal(t, 239, message.Gateway)
	assert.Len(t, message.TypeDetails, 0)
	assert.Equal(t, "plain", message.DataCoding)
	assert.Equal(t, 1, message.MClass)

	assert.Equal(t, expectedStatus, message.Recipients.Items[0].Status)
	if expectedStatus == "scheduled" {
		assert.NotNil(t, message.ScheduledDatetime)
		assert.Equal(t, 0, message.Recipients.TotalSentCount)
		assert.Nil(t, message.Recipients.Items[0].StatusDatetime)
	} else {
		assert.Nil(t, message.ScheduledDatetime)
		assert.Equal(t, 1, message.Recipients.TotalSentCount)
		assert.Equal(t, "2022-01-05T10:02:59Z", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}

	if message.CreatedDatetime == nil || message.CreatedDatetime.Format(time.RFC3339) != "2022-01-05T10:02:59Z" {
		t.Errorf("Unexpected message created datetime: %s, expected: 2022-01-05T10:02:59Z", message.CreatedDatetime)
	}
	assert.Equal(t, 1, message.Recipients.TotalCount)
	assert.Equal(t, int64(31612345678), message.Recipients.Items[0].Recipient)
	assert.Equal(t, 1, message.Recipients.Items[0].MessagePartCount)
}

func assertExtendedMessageObject(t *testing.T, message *Message) {
	assertMessageObject(t, message, "sent")

	assert.Equal(t, "Ukraine", *message.Recipients.Items[0].RecipientCountry)
	assert.Equal(t, 380, *message.Recipients.Items[0].RecipientCountryPrefix)
	assert.Equal(t, "life:)", *message.Recipients.Items[0].RecipientOperator)
	assert.Equal(t, 22, *message.Recipients.Items[0].MessageLength)
	assert.Equal(t, "successfully delivered", *message.Recipients.Items[0].StatusReason)
	assert.Equal(t, "25506", *message.Recipients.Items[0].Mccmnc)
	assert.Equal(t, "255", *message.Recipients.Items[0].Mcc)
	assert.Equal(t, "06", *message.Recipients.Items[0].Mnc)

	assert.Nil(t, message.Recipients.Items[0].StatusErrorCode)
	assert.NotNil(t, message.Recipients.Items[0].Price)
	assert.Equal(t, 22, message.Recipients.Items[0].Price.Amount)
	assert.Equal(t, "UAH", message.Recipients.Items[0].Price.Currency)
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	message, err := Create(client, "TestName", []string{"31612345678"}, "Hello World", nil)
	assert.NoError(t, err)

	assertMessageObject(t, message, "sent")
}

func TestCreateError(t *testing.T) {
	mbtest.WillReturnAccessKeyError()
	client := mbtest.Client(t)

	_, err := Create(client, "TestName", []string{"31612345678"}, "Hello World", nil)

	errorResponse, ok := err.(messagebird.ErrorResponse)
	assert.True(t, ok)
	assert.Equal(t, 1, len(errorResponse.Errors))
	assert.Equal(t, 2, errorResponse.Errors[0].Code)
	assert.Equal(t, "access_key", errorResponse.Errors[0].Parameter)
}

func TestCreateWithParams(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageWithParamsObject.json", http.StatusOK)
	client := mbtest.Client(t)

	params := &Params{
		Type:       "sms",
		Reference:  "TestReference",
		Validity:   13,
		Gateway:    10,
		DataCoding: "unicode",
	}

	message, err := Create(client, "TestName", []string{"31612345678"}, "Hello World", params)
	assert.NoError(t, err)
	assert.Equal(t, "sms", message.Type)
	assert.Equal(t, "TestReference", message.Reference)
	assert.Equal(t, 13, *message.Validity)
	assert.Equal(t, 10, message.Gateway)
	assert.Equal(t, "unicode", message.DataCoding)
}

func TestCreateWithBinaryType(t *testing.T) {
	mbtest.WillReturnTestdata(t, "binaryMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	params := &Params{
		Type:        "binary",
		TypeDetails: TypeDetails{"udh": "050003340201"},
	}

	message, err := Create(client, "TestName", []string{"31612345678"}, "Hello World", params)
	assert.NoError(t, err)
	assert.Equal(t, "binary", message.Type)
	assert.Len(t, message.TypeDetails, 1)
	assert.Equal(t, "050003340201", message.TypeDetails["udh"])
}

func TestCreateWithPremiumType(t *testing.T) {
	mbtest.WillReturnTestdata(t, "premiumMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	params := &Params{
		Type:        "premium",
		TypeDetails: TypeDetails{"keyword": "RESTAPI", "shortcode": 1008, "tariff": 150},
	}

	message, err := Create(client, "TestName", []string{"31612345678"}, "Hello World", params)
	assert.NoError(t, err)
	assert.Equal(t, "premium", message.Type)
	assert.Equal(t, 3, len(message.TypeDetails))
	assert.Equal(t, 150.0, message.TypeDetails["tariff"])
	assert.Equal(t, 1008.0, message.TypeDetails["shortcode"])
	assert.Equal(t, "RESTAPI", message.TypeDetails["keyword"])
}

func TestCreateWithFlashType(t *testing.T) {
	mbtest.WillReturnTestdata(t, "flashMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	params := &Params{Type: "flash"}

	message, err := Create(client, "TestName", []string{"31612345678"}, "Hello World", params)
	assert.NoError(t, err)
	assert.Equal(t, "flash", message.Type)
}

func TestCreateWithScheduledDatetime(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageObjectWithScheduledDatetime.json", http.StatusOK)
	client := mbtest.Client(t)

	scheduledDatetime, _ := time.Parse(time.RFC3339, "2022-01-05T10:03:59+00:00")

	params := &Params{ScheduledDatetime: scheduledDatetime}

	message, err := Create(client, "TestName", []string{"31612345678"}, "Hello World", params)
	assert.NoError(t, err)
	assert.Equal(t, scheduledDatetime.Format(time.RFC3339), message.ScheduledDatetime.Format(time.RFC3339))
	assert.Equal(t, 1, message.Recipients.TotalCount)
	assert.Equal(t, 0, message.Recipients.TotalSentCount)
	assert.Equal(t, int64(31612345678), message.Recipients.Items[0].Recipient)
	assert.Equal(t, "scheduled", message.Recipients.Items[0].Status)
	assert.Nil(t, message.Recipients.Items[0].StatusDatetime)
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "readMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	message, err := Read(client, "6fe65f90454aa61536e6a88b88972670")
	assert.NoError(t, err)

	assertExtendedMessageObject(t, message)
}

func TestReadScheduled(t *testing.T) {
	mbtest.WillReturnTestdata(t, "readScheduledMessageObject.json", http.StatusOK)
	client := mbtest.Client(t)

	message, err := Read(client, "6fe65f90454aa61536e6a88b88972670")
	assert.NoError(t, err)

	assertMessageObject(t, message, "scheduled")
}

func TestReadNotFound(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageNotFound.json", http.StatusNotFound)
	client := mbtest.Client(t)

	message, err := Read(client, "6fe65f90454aa61536e6a88b88972670")
	assert.Nil(t, message)
	assert.Errorf(t, err, "API errors: message not found")
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "messageListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	messageList, err := List(client, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, messageList.Offset)
	assert.Equal(t, 20, messageList.Limit)
	assert.Equal(t, 2, messageList.Count)
	assert.Equal(t, 2, messageList.TotalCount)
	assert.Equal(t, len(messageList.Items), messageList.Count)

	for _, message := range messageList.Items {
		assertMessageObject(t, &message, "sent")
	}
}

func TestListScheduled(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedStatusFilter := "status=scheduled"
		if !strings.Contains(r.URL.String(), expectedStatusFilter) {
			t.Errorf("API call should contain filter by status (%v), but is is not %v", expectedStatusFilter, r.URL.String())
		}
		_, err := w.Write(mbtest.Testdata(t, "messageListScheduledObject.json"))
		assert.NoError(t, err)
	})
	transport, teardown := mbtest.HTTPTestTransport(h)
	defer teardown()

	client := mbtest.Client(t)
	client.HTTPClient.Transport = transport

	messageList, err := List(client, &ListParams{Status: "scheduled"})
	assert.NoError(t, err)
	assert.Equal(t, 2, messageList.Count)
	assert.Equal(t, 2, messageList.TotalCount)
	assert.Equal(t, len(messageList.Items), messageList.Count)
}

func TestRequestDataForMessage(t *testing.T) {
	currentTime := time.Now()
	messageParams := &Params{
		Type:              "binary",
		Reference:         "MyReference",
		Validity:          1,
		Gateway:           0,
		TypeDetails:       nil,
		DataCoding:        "unicode",
		ScheduledDatetime: currentTime,
		ShortenURLs:       true,
	}

	request, err := requestDataForMessage("MSGBIRD", []string{"31612345678"}, "MyBody", messageParams)
	assert.NoError(t, err)
	assert.Equal(t, "MSGBIRD", request.Originator)
	assert.Equal(t, "31612345678", request.Recipients[0])
	assert.Equal(t, "MyBody", request.Body)
	assert.Equal(t, messageParams.Type, request.Type)
	assert.Equal(t, messageParams.Reference, request.Reference)
	assert.Equal(t, messageParams.GroupIds, request.GroupIds)
	assert.Equal(t, messageParams.Validity, request.Validity)
	assert.Equal(t, messageParams.Gateway, request.Gateway)
	assert.Nil(t, request.TypeDetails)
	assert.Equal(t, messageParams.DataCoding, request.DataCoding)
	assert.Equal(t, messageParams.ScheduledDatetime.Format(time.RFC3339), request.ScheduledDatetime)
	assert.True(t, request.ShortenURLs)
}

func TestDelete(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNoContent)
	client := mbtest.Client(t)

	err := Delete(client, "6fe65f90454aa61536e6a88b88972670")
	assert.NoError(t, err)
}

func TestDeleteNotFound(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNotFound)
	client := mbtest.Client(t)

	err := Delete(client, "6fe65f90454aa61536e6a88b88972670")
	assert.Errorf(t, err, "API errors: message not found")
}
