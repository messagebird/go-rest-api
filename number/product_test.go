package number

import (
	"github.com/messagebird/go-rest-api/v9/internal/mbtest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSearchProducts(t *testing.T) {
	mbtest.WillReturnTestdata(t, "searchProducts.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := SearchProducts(client, &ProductsRequest{
		CountryCode: "GB",
		Limit:       10,
		Features:    []string{"sms", "voice"},
		Type:        "mobile",
		Prefix:      "380",
	})
	assert.NoError(t, err)
	assert.Equal(t, 10, list.Limit)
	assert.Equal(t, 1, list.Count)
	assert.Len(t, list.Items, 1)
	assert.Equal(t, "bv", list.Items[0].NumberType)
	assert.True(t, list.Items[0].VerificationRequired)
	assert.Equal(t, "GB", list.Items[0].Country)
	assert.Equal(t, 1993, list.Items[0].ID)
	assert.Equal(t, "EUR", list.Items[0].Currency)
	assert.Equal(t, 1, list.Items[0].Price)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/products")
}

func TestReadProduct(t *testing.T) {
	mbtest.WillReturnTestdata(t, "readProduct.json", http.StatusOK)
	client := mbtest.Client(t)

	product, err := ReadProduct(client, "1993")
	assert.NoError(t, err)
	assert.Equal(t, "GB", product.Country)
	assert.Equal(t, "bv", product.NumberType)
	assert.Equal(t, "5 - 10 Business days", product.BackOrderLeadTime)
	assert.True(t, product.ReachableFromNationalFixed)
	assert.True(t, product.ReachableFromNationalMobile)
	assert.True(t, product.ReachableFromPayPhone)
	assert.True(t, product.VerificationRequired)
	assert.Equal(t, "1 year", product.InitialContractDuration)
	assert.Len(t, product.Prefixes, 3)
	assert.Len(t, product.Remarks, 2)
	assert.Len(t, product.Conditions, 9)
	assert.Len(t, product.EndUserData, 6)
	assert.Len(t, product.ForbiddenContent, 1)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/products/1993")
}
