package number

import (
	"net/http"
	"testing"
	"reflect"

	"../internal/mbtest"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestSearch(t *testing.T) {
	mbtest.WillReturnTestdata(t, "numberSearch.json", http.StatusOK)
	client := mbtest.Client(t)

	numLis, err := Search(client, "NL", &NumberListParams{
		Limit:         10,
		Features:      []string{"sms", "voice"},
		Type:          []string{"mobile"},
		SearchPattern: NumberPatternEnd,
	})
	if err != nil {
		t.Fatalf("unexpected error searching Numbers: %s", err)
	}

	if numLis.Items[0].Country != "NL" {
		t.Fatalf("got %s, expected NL", numLis.Items[0].Country)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/available-phone-numbers/NL")

	if query := mbtest.Request.URL.RawQuery; query != "features=sms&features=voice&limit=10&search_pattern=end&type=mobile" {
		t.Fatalf("got %s, expected features=sms&features=voice&limit=10&search_pattern=end&type=mobile", query)
	}
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "numberList.json", http.StatusOK)
	client := mbtest.Client(t)

	numLis, err := List(client, &NumberListParams{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error searching Numbers: %s", err)
	}

	if numLis.Items[0].Country != "NL" {
		t.Fatalf("got %s, expected NL", numLis.Items[0].Country)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/phone-numbers")

	if query := mbtest.Request.URL.RawQuery; query != "limit=10" {
		t.Fatalf("got %s, expected limit=10", query)
	}
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "numberRead.json", http.StatusOK)
	client := mbtest.Client(t)

	num, err := Read(client, "31612345670")
	if err != nil {
		t.Fatalf("unexpected error searching Numbers: %s", err)
	}

	if num.Number != "31612345670" {
		t.Fatalf("got %s, expected 31612345670", num.Number)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/phone-numbers/31612345670")
}

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	if err := Delete(client, "31612345670"); err != nil {
		t.Errorf("unexpected error canceling Number: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/phone-numbers/31612345670")
}

func TestUpdate(t *testing.T) {

	mbtest.WillReturnTestdata(t, "numberUpdatedObject.json", http.StatusOK)
	client := mbtest.Client(t)

	number, err := Update(client, "31612345670", &NumberUpdateRequest{
		Tags: []string{"tag1", "tag2", "tag3"},
	})

	if err != nil {
		t.Errorf("unexpected error updating Number: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPatch, "/v1/phone-numbers/31612345670")
	mbtest.AssertTestdata(t, "numberUpdateRequestObject.json", mbtest.Request.Body)
	
	if !reflect.DeepEqual(number.Tags, []string{"tag1", "tag2", "tag3"}) {
		t.Errorf("Unexpected number tags: %s, expected: ['tag1', 'tag2', 'tag3']", number.Tags)
	}
}

func TestPurchase(t *testing.T) {
	mbtest.WillReturnTestdata(t, "numberCreateObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	number, err := Purchase(client, &NumberPurchaseRequest{
			Number: "31971234567",
			Country: "NL",
			BillingIntervalMonths: 1, 
	})
	if err != nil {
		t.Errorf("unexpected error creating Number: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/phone-numbers")
	mbtest.AssertTestdata(t, "numberCreateRequestObject.json", mbtest.Request.Body)

	if number.Number != "31971234567" {
		t.Errorf("Unexpected number message id: %s, expected: 31971234567", number.Number)
	}

	if number.Country != "NL" {
		t.Errorf("Unexpected number country: %s, expected: NL", number.Country)
	}
}
