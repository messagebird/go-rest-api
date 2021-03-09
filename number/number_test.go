package number

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
	"github.com/stretchr/testify/assert"
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
		Type:          "mobile",
		SearchPattern: NumberPatternEnd,
	})
	assert.NoError(t, err)
	assert.Equal(t, "NL", numLis.Items[0].Country)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/available-phone-numbers/NL")

	query := mbtest.Request.URL.RawQuery
	assert.Equal(t, "features=sms&features=voice&limit=10&search_pattern=end&type=mobile", query)
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "numberList.json", http.StatusOK)
	client := mbtest.Client(t)

	numLis, err := List(client, &NumberListParams{Limit: 10})
	assert.NoError(t, err)
	assert.Equal(t, "NL", numLis.Items[0].Country)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/phone-numbers")

	query := mbtest.Request.URL.RawQuery
	assert.Equal(t, "limit=10", query)
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "numberRead.json", http.StatusOK)
	client := mbtest.Client(t)

	num, err := Read(client, "31612345670")
	assert.NoError(t, err)
	assert.Equal(t, "31612345670", num.Number)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/phone-numbers/31612345670")
}

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)

	err := Delete(client, "31612345670")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/phone-numbers/31612345670")
}

func TestUpdate(t *testing.T) {

	mbtest.WillReturnTestdata(t, "numberUpdatedObject.json", http.StatusOK)
	client := mbtest.Client(t)

	number, err := Update(client, "31612345670", &NumberUpdateRequest{
		Tags: []string{"tag1", "tag2", "tag3"},
	})
	assert.NoError(t, err)

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
		Number:                "31971234567",
		Country:               "NL",
		BillingIntervalMonths: 1,
	})
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/phone-numbers")
	mbtest.AssertTestdata(t, "numberCreateRequestObject.json", mbtest.Request.Body)
	assert.Equal(t, "31971234567", number.Number)
	assert.Equal(t, "NL", number.Country)
}
