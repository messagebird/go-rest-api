package number

import (
	"github.com/messagebird/go-rest-api/v9/internal/mbtest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCreatePool(t *testing.T) {
	mbtest.WillReturnTestdata(t, "createPool.json", http.StatusOK)
	client := mbtest.Client(t)

	pool, err := CreatePool(client, &CreatePoolRequest{
		PoolName:      "myPool",
		Service:       "randomcli",
		Configuration: &PoolConfiguration{false},
	})
	assert.NoError(t, err)
	assert.Equal(t, "726db88c-0fd6-44e8-8f0d-9ce2b2a5e16f", pool.ID)
	assert.Equal(t, "myPool", pool.Name)
	assert.Equal(t, "randomcli", pool.Service)
	assert.Equal(t, "2021-12-15T13:40:38Z", pool.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, "2021-12-16T13:40:38Z", pool.UpdatedAt.Format(time.RFC3339))
	assert.Equal(t, 0, pool.NumbersCount)
	assert.Equal(t, false, pool.Configuration.ByCountry)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/pools")
}

func TestReadPool(t *testing.T) {
	mbtest.WillReturnTestdata(t, "createPool.json", http.StatusOK)
	client := mbtest.Client(t)

	pool, err := ReadPool(client, "qwerty")
	assert.NoError(t, err)
	assert.Equal(t, "726db88c-0fd6-44e8-8f0d-9ce2b2a5e16f", pool.ID)
	assert.Equal(t, "myPool", pool.Name)
	assert.Equal(t, "randomcli", pool.Service)
	assert.Equal(t, "2021-12-15T13:40:38Z", pool.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, "2021-12-16T13:40:38Z", pool.UpdatedAt.Format(time.RFC3339))
	assert.Equal(t, 0, pool.NumbersCount)
	assert.Equal(t, false, pool.Configuration.ByCountry)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/pools/qwerty")
}

func TestUpdatePool(t *testing.T) {
	mbtest.WillReturnTestdata(t, "createPool.json", http.StatusOK)
	client := mbtest.Client(t)

	pool, err := UpdatePool(client, "qwerty2", &UpdatePoolRequest{
		PoolName:      "myPool",
		Configuration: &PoolConfiguration{false},
	})
	assert.NoError(t, err)
	assert.Equal(t, "726db88c-0fd6-44e8-8f0d-9ce2b2a5e16f", pool.ID)
	assert.Equal(t, "myPool", pool.Name)
	assert.Equal(t, "randomcli", pool.Service)
	assert.Equal(t, "2021-12-15T13:40:38Z", pool.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, "2021-12-16T13:40:38Z", pool.UpdatedAt.Format(time.RFC3339))
	assert.Equal(t, 0, pool.NumbersCount)
	assert.Equal(t, false, pool.Configuration.ByCountry)

	mbtest.AssertEndpointCalled(t, http.MethodPut, "/v1/pools/qwerty2")
}

func TestDeletePool(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNoContent)
	client := mbtest.Client(t)

	err := DeletePool(client, "qwerty3")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/pools/qwerty3")
}

func TestDeletePoolError(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNotAcceptable)
	client := mbtest.Client(t)

	err := DeletePool(client, "qwerty3")
	assert.Error(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/pools/qwerty3")
}

func TestListPool(t *testing.T) {
	mbtest.WillReturnTestdata(t, "listPool.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := ListPool(client, &ListPoolRequest{
		PoolName: "name",
		Service:  "randomcli",
		Limit:    1,
		Offset:   0,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, list.Limit)
	assert.Equal(t, 0, list.Offset)
	assert.Equal(t, 1, list.Count)
	assert.Equal(t, 5, list.TotalCount)
	assert.Len(t, list.Items, 1)
	assert.Equal(t, "myPool", list.Items[0].Name)
	assert.Equal(t, "randomcli", list.Items[0].Service)
	assert.Equal(t, "2021-12-16T15:27:04Z", list.Items[0].CreatedAt.Format(time.RFC3339))
	assert.Equal(t, 10, list.Items[0].NumbersCount)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/pools")
}

func TestListPoolNumbers(t *testing.T) {
	mbtest.WillReturnTestdata(t, "listPoolNumbers.json", http.StatusOK)
	client := mbtest.Client(t)

	list, err := ListPoolNumbers(client, "pool-name", &ListPoolNumbersRequest{
		Limit:  20,
		Offset: 0,
		Number: "316",
	})
	assert.NoError(t, err)
	assert.Equal(t, 20, list.Limit)
	assert.Equal(t, 0, list.Offset)
	assert.Equal(t, 3, list.Count)
	assert.Equal(t, 3, list.TotalCount)
	assert.Len(t, list.Numbers, 3)
	assert.Equal(t, []string{"31612345678", "31612345679", "31612345670"}, list.Numbers)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/pools/pool-name/numbers")
}

func TestAddNumberToPool(t *testing.T) {
	mbtest.WillReturnTestdata(t, "addNumberToPool.json", http.StatusOK)
	client := mbtest.Client(t)

	num, err := AddNumberToPool(client, "pool-name", []string{"31612345678", "31612345679", "31612345670"})
	assert.NoError(t, err)
	assert.Len(t, num.Success, 2)
	assert.Equal(t, []string{"31612345678", "31612345679"}, num.Success)
	assert.Len(t, num.Fail, 1)
	assert.Equal(t, "31612345670", num.Fail[0].Number)
	assert.Equal(t, "number is not verified", num.Fail[0].Error)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/pools/pool-name/numbers")
}

func TestDeleteNumberFromPool(t *testing.T) {
	mbtest.WillReturnOnlyStatus(http.StatusNoContent)
	client := mbtest.Client(t)

	err := DeleteNumberFromPool(client, "pool-name", []string{"31612345678", "31612345679", "31612345670"})
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/pools/pool-name/numbers")
}
