package number

import (
	"github.com/messagebird/go-rest-api/v9/internal/mbtest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPlaceBackorder(t *testing.T) {
	mbtest.WillReturnTestdata(t, "placeBackorder.json", http.StatusOK)
	client := mbtest.Client(t)

	backorderID, err := PlaceBackorder(client, &PlaceBackorderRequest{
		ProductID: 1993,
		Prefix:    "31114",
		Quantity:  3,
	})
	assert.NoError(t, err)
	assert.Equal(t, "48f6057c21de42d4bbf73fb86caaf361", string(backorderID))

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/backorders")
}

func TestReadBackorder(t *testing.T) {
	mbtest.WillReturnTestdata(t, "readBackorder.json", http.StatusOK)
	client := mbtest.Client(t)

	backorder, err := ReadBackorder(client, "vn4oor3c21de42d4bbf73fb86caaf361")
	assert.NoError(t, err)
	assert.Equal(t, "48f6057c21de42d4bbf73fb86caaf361", backorder.ID)
	assert.Equal(t, 1993, backorder.ProductID)
	assert.Equal(t, "GB", backorder.Country)
	assert.Equal(t, "44113", backorder.Prefix)
	assert.Equal(t, "blocked", backorder.Status)
	assert.Equal(t, []string{"MISSING_KYC", "MISSING_EUD"}, backorder.ReasonCodes)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/backorders/vn4oor3c21de42d4bbf73fb86caaf361")
}

func TestListBackorderDocuments(t *testing.T) {
	mbtest.WillReturnTestdata(t, "listBackorderDocuments.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := ListBackorderDocuments(client, "vn4oor3c21de42d4bbf73fb86caaf361")
	assert.NoError(t, err)
	assert.Equal(t, 20, list.Limit)
	assert.Equal(t, 1, list.Count)
	assert.Len(t, list.Items, 1)
	assert.Equal(t, 62, list.Items[0].ID)
	assert.Equal(t, "Proof of in region address (max 3 months old)", list.Items[0].Name)
	assert.Equal(t, "The number-user must provide proof that he has an address within the region of the prefix.", list.Items[0].Description)
	assert.Equal(t, "missing", list.Items[0].Status)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/backorders/vn4oor3c21de42d4bbf73fb86caaf361/documents")
}

func TestCreateBackorderDocument(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNoContent)
	client := mbtest.Client(t)

	err := CreateBackorderDocument(client, "vn4oor3c21de42d4bbf73fb86caaf361", &CreateBackorderDocumentRequest{
		ID:       62,
		Name:     "messagebird-kyc-upload-test.txt",
		MimeType: "text/plain",
		Content:  "aHR0cHM6Ly93d3cueW91dHViZS5jb20vd2F0Y2g/dj1kUXc0dzlXZ1hjUQ==",
	})
	assert.NoError(t, err)
	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/backorders/vn4oor3c21de42d4bbf73fb86caaf361/documents")
}

func TestCreateBackorderDocumentError(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNotAcceptable)
	client := mbtest.Client(t)

	err := CreateBackorderDocument(client, "vn4oor3c21de42d4bbf73fb86caaf361", &CreateBackorderDocumentRequest{
		ID:       62,
		Name:     "messagebird-kyc-upload-test.txt",
		MimeType: "text/plain",
		Content:  "aHR0cHM6Ly93d3cueW91dHViZS5jb20vd2F0Y2g/dj1kUXc0dzlXZ1hjUQ==",
	})
	assert.Error(t, err)
	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/backorders/vn4oor3c21de42d4bbf73fb86caaf361/documents")
}

func TestListBackorderEndUserDetails(t *testing.T) {
	mbtest.WillReturnTestdata(t, "listBackorderEndUserDetails.json", http.StatusOK)
	client := mbtest.Client(t)

	eud, err := ListBackorderEndUserDetails(client, "vn4oor3c21de42d4bbf73fb86caaf361")
	assert.NoError(t, err)
	assert.Len(t, eud.Items, 6)
	assert.Equal(t, []*EndUserDetail{
		{"CompanyName", "Company name"},
		{"Street", "Street"},
		{"StreetNumber", "Street number"},
		{"ZipCode", "Zip code"},
		{"City", "City"},
		{"Country", "Country"},
	}, eud.Items)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/backorders/vn4oor3c21de42d4bbf73fb86caaf361/end-user-details")
}

func TestCreateBackorderEndUserDetail(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNoContent)
	client := mbtest.Client(t)

	err := CreateBackorderEndUserDetail(client, "vn4oor3c21de42d4bbf73fb86caaf361", &CreateBackorderEndUserDetailRequest{
		CompanyName:  "Messagebird",
		Street:       "Burgwal",
		StreetNumber: "35",
		ZipCode:      "2011",
		City:         "Amsterdam",
		Country:      "Netherlands",
	})
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/backorders/vn4oor3c21de42d4bbf73fb86caaf361/end-user-details")
}
