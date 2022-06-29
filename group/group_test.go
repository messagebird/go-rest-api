package group

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/v8/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "groupObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	group, err := Create(client, &Request{"Friends"})

	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/groups")
	mbtest.AssertTestData(t, "groupRequestCreateObject.json", mbtest.Request.Body)
	assert.Equal(t, "Friends", group.Name)
}

func TestCreateWithEmptyName(t *testing.T) {
	client := mbtest.Client(t)

	_, err := Create(client, &Request{""})
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	err := Delete(client, "group-id")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/groups/group-id")
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "groupListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := List(client, DefaultListOptions)
	assert.NoError(t, err)
	assert.Equal(t, 0, list.Offset)
	assert.Equal(t, 10, list.Limit)
	assert.Equal(t, 2, list.Count)
	assert.Equal(t, 2, list.TotalCount)

	assert.Len(t, list.Items, 2)
	assert.Equal(t, "first-id", list.Items[0].ID)
	assert.Equal(t, "second-id", list.Items[1].ID)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/groups")
}

func TestListPagination(t *testing.T) {
	client := mbtest.Client(t)

	tt := []struct {
		expected string
		options  *ListOptions
	}{
		{"limit=10&offset=0", DefaultListOptions},
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
	mbtest.WillReturnTestdata(t, "groupObject.json", http.StatusOK)
	client := mbtest.Client(t)

	group, err := Read(client, "group-id")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/groups/group-id")
	assert.Equal(t, "group-id", group.ID)
	assert.Equal(t, "https://rest.messagebird.com/groups/group-id", group.HRef)
	assert.Equal(t, "Friends", group.Name)
	assert.Equal(t, 3, group.Contacts.TotalCount)
	assert.Equal(t, "https://rest.messagebird.com/groups/group-id", group.Contacts.HRef)

	created, _ := time.Parse(time.RFC3339, "2018-07-25T12:16:10+00:00")
	assert.True(t, created.Equal(*group.CreatedDatetime))

	updated, _ := time.Parse(time.RFC3339, "2018-07-25T12:16:23+00:00")
	assert.True(t, updated.Equal(*group.UpdatedDatetime))
}

func TestUpdate(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	err := Update(client, "group-id", &Request{"Family"})
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodPatch, "/groups/group-id")
	mbtest.AssertTestData(t, "groupRequestUpdateObject.json", mbtest.Request.Body)
	assert.Equal(t, "application/json", mbtest.Request.ContentType)
}

func TestAddContacts(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	err := AddContacts(client, "group-id", []string{"first-contact-id", "second-contact-id"})
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodPut, "/groups/group-id/contacts")
	mbtest.AssertTestData(t, "groupRequestAddContactsObject.txt", mbtest.Request.Body)
	assert.Equal(t, "application/x-www-form-urlencoded", mbtest.Request.ContentType)
}

func TestAddContactsWithEmptyContacts(t *testing.T) {
	client := mbtest.Client(t)

	tt := []struct {
		contactIDs []string
	}{
		{[]string{}},
		{nil},
	}

	for _, tc := range tt {
		err := AddContacts(client, "group-id", tc.contactIDs)
		assert.Error(t, err)
	}
}

func TestAddContactsWithTooManyContacts(t *testing.T) {
	client := mbtest.Client(t)

	// Only 50 contacts are allowed at a time.
	contactIDs := make([]string, 51)

	err := AddContacts(client, "group-id", contactIDs)
	assert.Error(t, err)
}

func TestListContacts(t *testing.T) {
	mbtest.WillReturnTestdata(t, "groupContactListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := ListContacts(client, "group-id", DefaultListOptions)
	assert.NoError(t, err)
	assert.Equal(t, 0, list.Offset)
	assert.Equal(t, 20, list.Limit)
	assert.Equal(t, 3, list.Count)
	assert.Equal(t, 3, list.TotalCount)
	assert.Equal(t, "first-contact-id", list.Items[0].ID)
	assert.Equal(t, "second-contact-id", list.Items[1].ID)
	assert.Equal(t, "third-contact-id", list.Items[2].ID)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/groups/group-id/contacts")
}

func TestRemoveContact(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	err := RemoveContact(client, "group-id", "contact-id")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/groups/group-id/contacts/contact-id")
}
