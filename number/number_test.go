package number

import (
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/internal/mbtest"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestSearch(t *testing.T) {
	mbtest.WillReturnTestdata(t, "numberSearch.json", http.StatusOK)
	client := mbtest.Client(t)

	numLis, err := Search(client, "NL", &NumberListParams{Limit: 10, Features: []string{"sms", "voice"}, Type: []string{"mobile"}})
	if err != nil {
		t.Fatalf("unexpected error searching Numbers: %s", err)
	}

	if numLis.Limit > 100 {
		t.Fatalf("got %d, expected <= 100", numLis.Limit)
	}

	if numLis.Items[0].Country != "NL" {
		t.Fatalf("got %s, expected NL", numLis.Items[0].Country)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/available-phone-numbers/NL")

	if query := mbtest.Request.URL.RawQuery; query != "features=sms&features=voice&limit=10&type=mobile" {
		t.Fatalf("got %s, expected features=sms&features=voice&limit=10", query)
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
