package contact

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/internal/messagebirdtest"
)

func TestMain(m *testing.M) {
	messagebirdtest.EnableServer(m)
}

func TestCreateWithEmptyMSISDN(t *testing.T) {
	client := messagebirdtest.Client(t)

	if _, err := Create(client, &Request{}); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCreate(t *testing.T) {
	messagebirdtest.WillReturnTestdata(t, "contactObject.json", http.StatusCreated)
	client := messagebirdtest.Client(t)

	contact, err := Create(client, &Request{
		MSISDN:    "31612345678",
		FirstName: "Foo",
		LastName:  "Bar",
		Custom1:   "First",
		Custom2:   "Second",
	})
	if err != nil {
		t.Fatalf("unexpected error creating Contact: %s", err)
	}

	messagebirdtest.AssertEndpointCalled(t, http.MethodPost, "/contacts")
	messagebirdtest.AssertTestdata(t, "contactRequestObjectCreate.json", messagebirdtest.Request.Body)

	if contact.MSISDN != 31612345678 {
		t.Fatalf("expected 31612345678, got %d", contact.MSISDN)
	}

	if contact.FirstName != "Foo" {
		t.Fatalf("expected Foo, got %s", contact.FirstName)
	}

	if contact.LastName != "Bar" {
		t.Fatalf("expected Bar, got %s", contact.LastName)
	}

	if contact.CustomDetails.Custom1 != "First" {
		t.Fatalf("expected First, got %s", contact.CustomDetails.Custom1)
	}

	if contact.CustomDetails.Custom2 != "Second" {
		t.Fatalf("expected Second, got %s", contact.CustomDetails.Custom2)
	}

	if contact.CustomDetails.Custom3 != "Third" {
		t.Fatalf("expected Third, got %s", contact.CustomDetails.Custom3)
	}

	if contact.CustomDetails.Custom4 != "Fourth" {
		t.Fatalf("expected Fourth, got %s", contact.CustomDetails.Custom4)
	}
}

func TestDelete(t *testing.T) {
	messagebirdtest.WillReturn([]byte(""), http.StatusNoContent)
	client := messagebirdtest.Client(t)

	if err := Delete(client, "contact-id"); err != nil {
		t.Fatalf("unexpected error deleting Contact: %s", err)
	}

	messagebirdtest.AssertEndpointCalled(t, http.MethodDelete, "/contacts/contact-id")
}

func TestDeleteWithEmptyID(t *testing.T) {
	client := messagebirdtest.Client(t)

	if err := Delete(client, ""); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestList(t *testing.T) {
	messagebirdtest.WillReturnTestdata(t, "contactListObject.json", http.StatusOK)
	client := messagebirdtest.Client(t)

	list, err := List(client, DefaultListOptions)
	if err != nil {
		t.Fatalf("unexpected error retrieving Contact list: %s", err)
	}

	if list.Offset != 0 {
		t.Fatalf("expected 0, got %d", list.Offset)
	}

	if list.Limit != 20 {
		t.Fatalf("expected 0, got %d", list.Limit)
	}

	if list.Count != 2 {
		t.Fatalf("expected 2, got %d", list.Count)
	}

	if list.TotalCount != 2 {
		t.Fatalf("expected 2, got %d", list.TotalCount)
	}

	if actualCount := len(list.Items); actualCount != 2 {
		t.Fatalf("expected 2, got %d", actualCount)
	}

	if list.Items[0].ID != "first-id" {
		t.Fatalf("expected first-id, got %s", list.Items[0].ID)
	}

	if list.Items[1].ID != "second-id" {
		t.Fatalf("expected second-id, got %s", list.Items[1].ID)
	}

	messagebirdtest.AssertEndpointCalled(t, http.MethodGet, "/contacts")
}

func TestListPagination(t *testing.T) {
	client := messagebirdtest.Client(t)

	tt := []struct {
		expected string
		options  *ListOptions
	}{
		{"limit=20&offset=0", DefaultListOptions},
		{"limit=10&offset=25", &ListOptions{10, 25}},
		{"limit=50&offset=10", &ListOptions{50, 10}},
	}

	for _, tc := range tt {
		List(client, tc.options)

		if query := messagebirdtest.Request.URL.RawQuery; query != tc.expected {
			t.Fatalf("expected %s, got %s", tc.expected, query)
		}
	}
}

func TestRead(t *testing.T) {
	messagebirdtest.WillReturnTestdata(t, "contactObject.json", http.StatusOK)
	client := messagebirdtest.Client(t)

	contact, err := Read(client, "contact-id")
	if err != nil {
		t.Fatalf("unexpected error reading Contact: %s", err)
	}

	messagebirdtest.AssertEndpointCalled(t, http.MethodGet, "/contacts/contact-id")

	if contact.ID != "contact-id" {
		t.Fatalf("expected contact-id, got %s", contact.ID)
	}

	if contact.HRef != "https://rest.messagebird.com/contacts/contact-id" {
		t.Fatalf("expected https://rest.messagebird.com/contacts/contact-id, got %s", contact.HRef)
	}

	if contact.MSISDN != 31612345678 {
		t.Fatalf("expected 31612345678, got %d", contact.MSISDN)
	}

	if contact.FirstName != "Foo" {
		t.Fatalf("expected Foo, got %s", contact.FirstName)
	}

	if contact.LastName != "Bar" {
		t.Fatalf("expected Bar, got %s", contact.LastName)
	}

	if contact.Groups.TotalCount != 3 {
		t.Fatalf("expected 3, got %d", contact.Groups.TotalCount)
	}

	if contact.Groups.HRef != "https://rest.messagebird.com/contacts/contact-id/groups" {
		t.Fatalf("expected https://rest.messagebird.com/contacts/contact-id/groups, got %s", contact.Groups.HRef)
	}

	if contact.Messages.TotalCount != 5 {
		t.Fatalf("expected 5, got %d", contact.Messages.TotalCount)
	}

	if contact.Messages.HRef != "https://rest.messagebird.com/contacts/contact-id/messages" {
		t.Fatalf("expected https://rest.messagebird.com/contacts/contact-id/messages, got %s", contact.Messages.HRef)
	}

	expectedCreatedDatetime, _ := time.Parse(time.RFC3339, "2018-07-13T10:34:08+00:00")
	if !contact.CreatedDatetime.Equal(expectedCreatedDatetime) {
		t.Fatalf("expected %s, got %s", expectedCreatedDatetime, contact.CreatedDatetime)
	}

	expectedUpdatedDatetime, _ := time.Parse(time.RFC3339, "2018-07-13T10:44:08+00:00")
	if !contact.UpdatedDatetime.Equal(expectedUpdatedDatetime) {
		t.Fatalf("expected %s, got %s", expectedUpdatedDatetime, contact.UpdatedDatetime)
	}
}

func TestReadWithCustomDetails(t *testing.T) {
	messagebirdtest.WillReturnTestdata(t, "contactObjectWithCustomDetails.json", http.StatusOK)
	client := messagebirdtest.Client(t)

	contact, err := Read(client, "contact-id")
	if err != nil {
		t.Fatalf("unexpected error reading Contact with custom details: %s", err)
	}

	messagebirdtest.AssertEndpointCalled(t, http.MethodGet, "/contacts/contact-id")

	if contact.CustomDetails.Custom1 != "First" {
		t.Fatalf("expected First, got %s", contact.CustomDetails.Custom1)
	}

	if contact.CustomDetails.Custom2 != "Second" {
		t.Fatalf("expected Second, got %s", contact.CustomDetails.Custom2)
	}

	if contact.CustomDetails.Custom3 != "Third" {
		t.Fatalf("expected Third, got %s", contact.CustomDetails.Custom3)
	}

	if contact.CustomDetails.Custom4 != "Fourth" {
		t.Fatalf("expected Fourth, got %s", contact.CustomDetails.Custom4)
	}
}

func TestUpdate(t *testing.T) {
	client := messagebirdtest.Client(t)

	tt := []struct {
		expectedTestdata string
		contactRequest   *Request
	}{
		{"contactRequestObjectUpdateCustom.json", &Request{Custom1: "Foo", Custom4: "Bar"}},
		{"contactRequestObjectUpdateMSISDN.json", &Request{MSISDN: "31687654321"}},
		{"contactRequestObjectUpdateName.json", &Request{FirstName: "Message", LastName: "Bird"}},
	}

	for _, tc := range tt {
		if _, err := Update(client, "contact-id", tc.contactRequest); err != nil {
			t.Fatalf("unexpected error updating Contact: %s\n", err)
		}

		messagebirdtest.AssertEndpointCalled(t, http.MethodPatch, "/contacts/contact-id")
		messagebirdtest.AssertTestdata(t, tc.expectedTestdata, messagebirdtest.Request.Body)
	}
}
