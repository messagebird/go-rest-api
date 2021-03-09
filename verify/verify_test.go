package verify

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func assertVerifyObject(t *testing.T, v *Verify) {
	assert.NotNil(t, v)
	assert.Equal(t, "15498233759288aaf929661v21936686", v.ID)
	assert.Equal(t, "https://rest.messagebird.com/verify/15498233759288aaf929661v21936686", v.HRef)
	assert.Equal(t, 31612345678, v.Recipient)
	assert.Equal(t, "MyReference", v.Reference)
	assert.Len(t, v.Messages, 1)
	assert.Equal(t, "https://rest.messagebird.com/messages/c2bbd563759288aaf962910b56023756", v.Messages["href"])
	assert.Equal(t, "sent", v.Status)

	assert.Equal(t, "2017-05-26T20:06:07Z", v.CreatedDatetime.Format(time.RFC3339))

	assert.Equal(t, "2017-05-26T20:06:37Z", v.ValidUntilDatetime.Format(time.RFC3339))
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "verifyObject.json", http.StatusOK)
	client := mbtest.Client(t)

	v, err := Create(client, "31612345678", nil)
	assert.NoError(t, err)

	assertVerifyObject(t, v)
}

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	err := Delete(client, "15498233759288aaf929661v21936686")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/verify/15498233759288aaf929661v21936686")
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "verifyObject.json", http.StatusOK)
	client := mbtest.Client(t)

	v, err := Read(client, "15498233759288aaf929661v21936686")
	assert.NoError(t, err)
	assert.Equal(t, "15498233759288aaf929661v21936686", v.ID)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/verify/15498233759288aaf929661v21936686")
}

func TestVerifyToken(t *testing.T) {
	mbtest.WillReturnTestdata(t, "verifyTokenObject.json", http.StatusOK)
	client := mbtest.Client(t)

	v, err := VerifyToken(client, "does not", "matter")
	assert.NoError(t, err)

	assertVerifyTokenObject(t, v)
}

func assertVerifyTokenObject(t *testing.T, v *Verify) {
	assert.NotNil(t, v)
	assert.Equal(t, "a3f2edb23592d68163f9694v13904556", v.ID)
	assert.Equal(t, "https://rest.messagebird.com/verify/a3f2edb23592d68163f9694v13904556", v.HRef)
	assert.Equal(t, 31612345678, v.Recipient)
	assert.Equal(t, "MyReference", v.Reference)
	assert.Len(t, v.Messages, 1)
	assert.Equal(t, "https://rest.messagebird.com/messages/63b168423592d681641eb07b76226648", v.Messages["href"])
	assert.Equal(t, "verified", v.Status)

	assert.Equal(t, "2017-05-30T12:39:50Z", v.CreatedDatetime.Format(time.RFC3339))
	assert.Equal(t, "2017-05-30T12:40:20Z", v.ValidUntilDatetime.Format(time.RFC3339))
}

func TestRequestDataForVerify(t *testing.T) {
	verifyParams := &Params{
		Originator:  "MSGBIRD",
		Reference:   "MyReference",
		Type:        "sms",
		Template:    "Your code is: %token",
		DataCoding:  "plain",
		ReportURL:   "http://example.com/report",
		Voice:       "male",
		Language:    "en-gb",
		Timeout:     20,
		TokenLength: 8,
	}

	requestData, err := requestDataForVerify("31612345678", verifyParams)
	assert.NoError(t, err)
	assert.Equal(t, "31612345678", requestData.Recipient)
	assert.Equal(t, "MSGBIRD", requestData.Originator)
	assert.Equal(t, "MyReference", requestData.Reference)
	assert.Equal(t, "sms", requestData.Type)
	assert.Equal(t, "plain", requestData.DataCoding)
	assert.Equal(t, "http://example.com/report", requestData.ReportURL)
	assert.Equal(t, "male", requestData.Voice)
	assert.Equal(t, "en-gb", requestData.Language)
	assert.Equal(t, 20, requestData.Timeout)
	assert.Equal(t, 8, requestData.TokenLength)
}
