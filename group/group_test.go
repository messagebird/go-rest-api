package group

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/internal/mbtest"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestCreate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "groupObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	group, err := Create(client, &Request{"Friends"})
	if err != nil {
		t.Fatalf("unexpected error creating Group: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/groups")
	mbtest.AssertTestdata(t, "groupRequestCreateObject.json", mbtest.Request.Body)

	if group.Name != "Friends" {
		t.Fatalf("expected Friends, got %s", group.Name)
	}
}

func TestCreateWithEmptyName(t *testing.T) {
	client := mbtest.Client(t)

	if _, err := Create(client, &Request{""}); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	if err := Delete(client, "group-id"); err != nil {
		t.Fatalf("unexpected error deleting Group: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/groups/group-id")
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "groupListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := List(client, DefaultListOptions)
	if err != nil {
		t.Fatalf("unexpected error retrieving Contact list: %s", err)
	}

	if list.Offset != 0 {
		t.Fatalf("expected 0, got %d", list.Offset)
	}

	if list.Limit != 10 {
		t.Fatalf("expected 10, got %d", list.Limit)
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
		List(client, tc.options)

		if query := mbtest.Request.URL.RawQuery; query != tc.expected {
			t.Fatalf("expected %s, got %s", tc.expected, query)
		}
	}
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "groupObject.json", http.StatusOK)
	client := mbtest.Client(t)

	group, err := Read(client, "group-id")
	if err != nil {
		t.Fatalf("unexpected error reading Group: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/groups/group-id")

	if group.ID != "group-id" {
		t.Fatalf("expected group-id, got %s", group.ID)
	}

	if group.HRef != "https://rest.messagebird.com/groups/group-id" {
		t.Fatalf("expected https://rest.messagebird.com/groups/group-id, got %s", group.HRef)
	}

	if group.Name != "Friends" {
		t.Fatalf("expected Friends, got %s", group.Name)
	}

	if group.Contacts.TotalCount != 3 {
		t.Fatalf("expected 3, got %d", group.Contacts.TotalCount)
	}

	if group.Contacts.HRef != "https://rest.messagebird.com/groups/group-id" {
		t.Fatalf("expected https://rest.messagebird.com/groups/group-id, got %s", group.Contacts.HRef)
	}

	if created, _ := time.Parse(time.RFC3339, "2018-07-25T12:16:10+00:00"); !created.Equal(group.CreatedDatetime) {
		t.Fatalf("expected 2018-07-25T12:16:10+00:00, got %s", group.CreatedDatetime)
	}

	if updated, _ := time.Parse(time.RFC3339, "2018-07-25T12:16:23+00:00"); !updated.Equal(group.UpdatedDatetime) {
		t.Fatalf("expected 2018-07-25T12:16:23+00:00, got %s", group.UpdatedDatetime)
	}
}

func TestUpdate(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	if err := Update(client, "group-id", &Request{"Family"}); err != nil {
		t.Fatalf("unexpected error updating Group: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPatch, "/groups/group-id")
	mbtest.AssertTestdata(t, "groupRequestUpdateObject.json", mbtest.Request.Body)
}

func TestAddContacts(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	if err := AddContacts(client, "group-id", []string{"first-contact-id", "second-contact-id"}); err != nil {
		t.Fatalf("unexpected error removing Contacts from Group: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/groups/group-id/contacts")

	if mbtest.Request.URL.RawQuery != "_method=PUT&ids[]=first-contact-id&ids[]=second-contact-id" {
		t.Fatalf("expected _method=PUT&ids[]=first-contact-id&ids[]=second-contact-id, got %s", mbtest.Request.URL.RawQuery)
	}
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
		if err := AddContacts(client, "group-id", tc.contactIDs); err == nil {
			t.Fatalf("expected error, got nil")
		}
	}
}

func TestAddContactsWithTooManyContacts(t *testing.T) {
	client := mbtest.Client(t)

	// Only 50 contacts are allowed at a time.
	contactIDs := make([]string, 51)

	if err := AddContacts(client, "group-id", contactIDs); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestListContacts(t *testing.T) {
	mbtest.WillReturnTestdata(t, "groupContactListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := ListContacts(client, "group-id", DefaultListOptions)
	if err != nil {
		t.Fatalf("unexpected error listing Contacts: %s", err)
	}

	if list.Offset != 0 {
		t.Fatalf("expected 0, got %d", list.Offset)
	}

	if list.Limit != 20 {
		t.Fatalf("expected 20, got %d", list.Limit)
	}

	if list.Count != 3 {
		t.Fatalf("expected 3, got %d", list.Count)
	}

	if list.TotalCount != 3 {
		t.Fatalf("expected 3, got %d", list.TotalCount)
	}

	if list.Items[0].ID != "first-contact-id" {
		t.Fatalf("expected first-contact-id, got %s", list.Items[0].ID)
	}

	if list.Items[1].ID != "second-contact-id" {
		t.Fatalf("expected second-contact-id, got %s", list.Items[1].ID)
	}

	if list.Items[2].ID != "third-contact-id" {
		t.Fatalf("expected third-contact-id, got %s", list.Items[2].ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/groups/group-id/contacts")
}

func TestRemoveContact(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	if err := RemoveContact(client, "group-id", "contact-id"); err != nil {
		t.Fatalf("unexpected error deleting Group: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/groups/group-id/contacts/contact-id")
}
