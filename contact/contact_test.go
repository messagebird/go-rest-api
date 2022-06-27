package contact

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/v7/internal/mbtest"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestCreateWithEmptyMSISDN(t *testing.T) {
	client := mbtest.Client(t)

	_, err := Create(client, &Request{})
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "contactObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	contact, err := Create(client, &Request{
		MSISDN:    "31612345678",
		FirstName: "Foo",
		LastName:  "Bar",
		Custom1:   "First",
		Custom2:   "Second",
	})
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/contacts")
	mbtest.AssertTestData(t, "contactRequestObjectCreate.json", mbtest.Request.Body)

	assert.Equal(t, int64(31612345678), contact.MSISDN)

	assert.Equal(t, "Foo", contact.FirstName)
	assert.Equal(t, "Bar", contact.LastName)
	assert.Equal(t, "First", contact.CustomDetails.Custom1)
	assert.Equal(t, "Second", contact.CustomDetails.Custom2)
	assert.Equal(t, "Third", contact.CustomDetails.Custom3)
	assert.Equal(t, "Fourth", contact.CustomDetails.Custom4)
}

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	err := Delete(client, "contact-id")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/contacts/contact-id")
}

func TestDeleteWithEmptyID(t *testing.T) {
	client := mbtest.Client(t)

	err := Delete(client, "")
	assert.Error(t, err)
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "contactListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := List(client, DefaultListOptions)
	assert.NoError(t, err)

	assert.Equal(t, 0, list.Offset)
	assert.Equal(t, 20, list.Limit)
	assert.Equal(t, 2, list.Count)
	assert.Equal(t, 2, list.TotalCount)
	assert.Len(t, list.Items, 2)

	assert.Equal(t, "first-id", list.Items[0].ID)
	assert.Equal(t, "second-id", list.Items[1].ID)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/contacts")
}

func TestListPagination(t *testing.T) {
	client := mbtest.Client(t)

	tt := []struct {
		expected string
		options  *ListOptions
	}{
		{"limit=20&offset=0", DefaultListOptions},
		{"limit=10&offset=25", &ListOptions{10, 25}},
		{"limit=50&offset=10", &ListOptions{50, 10}},
	}

	for _, tc := range tt {
		_, err := List(client, tc.options)
		assert.NoError(t, err)

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, tc.expected, query)
	}
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "contactObject.json", http.StatusOK)
	client := mbtest.Client(t)

	contact, err := Read(client, "contact-id")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/contacts/contact-id")

	assert.Equal(t, "contact-id", contact.ID)
	assert.Equal(t, "https://rest.messagebird.com/contacts/contact-id", contact.HRef)
	assert.Equal(t, int64(31612345678), contact.MSISDN)
	assert.Equal(t, "Foo", contact.FirstName)
	assert.Equal(t, "Bar", contact.LastName)
	assert.Equal(t, 3, contact.Groups.TotalCount)
	assert.Equal(t, "https://rest.messagebird.com/contacts/contact-id/groups", contact.Groups.HRef)
	assert.Equal(t, 5, contact.Messages.TotalCount)
	assert.Equal(t, "https://rest.messagebird.com/contacts/contact-id/messages", contact.Messages.HRef)

	expectedCreatedDatetime, _ := time.Parse(time.RFC3339, "2018-07-13T10:34:08+00:00")
	assert.True(t, contact.CreatedDatetime.Equal(expectedCreatedDatetime))

	expectedUpdatedDatetime, _ := time.Parse(time.RFC3339, "2018-07-13T10:44:08+00:00")
	assert.True(t, contact.UpdatedDatetime.Equal(expectedUpdatedDatetime))
}

func TestReadWithCustomDetails(t *testing.T) {
	mbtest.WillReturnTestdata(t, "contactObjectWithCustomDetails.json", http.StatusOK)
	client := mbtest.Client(t)

	contact, err := Read(client, "contact-id")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/contacts/contact-id")

	assert.Equal(t, "First", contact.CustomDetails.Custom1)
	assert.Equal(t, "Second", contact.CustomDetails.Custom2)
	assert.Equal(t, "Third", contact.CustomDetails.Custom3)
	assert.Equal(t, "Fourth", contact.CustomDetails.Custom4)
}

func TestUpdate(t *testing.T) {
	client := mbtest.Client(t)

	tt := []struct {
		expectedTestdata string
		contactRequest   *Request
	}{
		{"contactRequestObjectUpdateCustom.json", &Request{Custom1: "Foo", Custom4: "Bar"}},
		{"contactRequestObjectUpdateMSISDN.json", &Request{MSISDN: "31687654321"}},
		{"contactRequestObjectUpdateName.json", &Request{FirstName: "Message", LastName: "Bird"}},
	}

	for _, tc := range tt {
		_, err := Update(client, "contact-id", tc.contactRequest)
		assert.NoError(t, err)

		mbtest.AssertEndpointCalled(t, http.MethodPatch, "/contacts/contact-id")
		mbtest.AssertTestData(t, tc.expectedTestdata, mbtest.Request.Body)
	}
}
