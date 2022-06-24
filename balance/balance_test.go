package balance

import (
	"net/http"
	"testing"

	messagebird "github.com/messagebird/go-rest-api/v8"
	"github.com/messagebird/go-rest-api/v8/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "balance.json", http.StatusOK)
	client := mbtest.Client(t)

	balance, err := Read(client)

	assert.NoError(t, err)

	assert.Equal(t, "prepaid", balance.Payment)

	assert.Equal(t, "credits", balance.Type)

	assert.EqualValuesf(t, 9.2, balance.Amount, "Unexpected balance amount: %.2f", balance.Amount)
}

func TestReadError(t *testing.T) {
	mbtest.WillReturnAccessKeyError()
	client := mbtest.Client(t)

	_, err := Read(client)

	errorResponse, ok := err.(messagebird.ErrorResponse)

	assert.True(t, ok)

	assert.Len(t, errorResponse.Errors, 1)

	assert.Equal(t, 2, errorResponse.Errors[0].Code)

	assert.Equal(t, "access_key", errorResponse.Errors[0].Parameter)
}
