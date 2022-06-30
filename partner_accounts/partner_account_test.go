package partner_accounts

import (
	"github.com/messagebird/go-rest-api/v8/internal/mbtest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestCreateChildAccount(t *testing.T) {
	mbtest.WillReturnTestdata(t, "createChildAccountResponse.json", http.StatusCreated)
	client := mbtest.Client(t)

	acc, err := CreateChildAccount(client, "Partner Account 3 Sub 1")
	assert.NoError(t, err)
	assert.Equal(t, 6249799, acc.ID)
	assert.Equal(t, "Partner Account 3 Sub 1", acc.Name)
	assert.Equal(t, "7qxJg4lsDKLAEBXAdxyarcwwvDn7YB00", acc.SigningKey)
	assert.Len(t, acc.AccessKeys, 2)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/child-accounts")
}

func TestUpdateChildAccount(t *testing.T) {
	mbtest.WillReturnTestdata(t, "updateChildAccountResponse.json", http.StatusCreated)
	client := mbtest.Client(t)

	acc, err := UpdateChildAccount(client, "6249609", "Partner Account 1 Sub 2")
	assert.NoError(t, err)
	assert.Equal(t, 6249609, acc.ID)
	assert.Equal(t, "Partner Account 1 Sub 2", acc.Name)
	assert.Equal(t, "", acc.SigningKey)
	assert.Len(t, acc.AccessKeys, 0)

	mbtest.AssertEndpointCalled(t, http.MethodPatch, "/v1/child-accounts/6249609")
}

func TestReadChildAccount(t *testing.T) {
	mbtest.WillReturnTestdata(t, "readChildAccountResponse.json", http.StatusCreated)
	client := mbtest.Client(t)

	acc, err := ReadChildAccount(client, "6249609")
	assert.NoError(t, err)
	assert.Equal(t, 6249609, acc.ID)
	assert.Equal(t, "Partner Account 1 Sub 1", acc.Name)
	assert.Equal(t, "", acc.SigningKey)
	assert.True(t, acc.InvoiceAggregation)
	assert.Len(t, acc.AccessKeys, 0)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/child-accounts/6249609")
}

func TestListChildAccount(t *testing.T) {
	mbtest.WillReturnTestdata(t, "listChildAccountResponse.json", http.StatusCreated)
	client := mbtest.Client(t)

	acc, err := ListChildAccount(client)
	assert.NoError(t, err)
	assert.Len(t, *acc, 3)

	expected := []struct {
		id   int
		name string
	}{{6249623, "Partner Account 1 Sub 1"}, {6249654, "Partner Account 1 Sub 2"}, {62496654, "Partner Account 1 Sub 3"}}

	for k, v := range expected {
		assert.Equal(t, v.id, (*acc)[k].ID)
		assert.Equal(t, v.name, (*acc)[k].Name)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/child-accounts")
}

func TestDeleteChildAccount(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNoContent)
	client := mbtest.Client(t)

	err := DeleteChildAccount(client, "6249633")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/child-accounts/6249633")
}

func TestDeleteChildAccountError(t *testing.T) {
	mbtest.WillReturnTestdata(t, "accountNotFound.json", http.StatusNotFound)
	client := mbtest.Client(t)

	err := DeleteChildAccount(client, "6249633")
	//assert.Error(t, err)
	assert.EqualError(t, err, "An error occurred: Account Not found")

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/child-accounts/6249633")
}
