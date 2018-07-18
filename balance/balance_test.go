package balance

import (
	"net/http"
	"testing"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/internal/messagebirdtest"
)

var balanceObject = []byte(`{
  		"payment":"prepaid",
  		"type":"credits",
  		"amount":9.2
	}`)

const Epsilon float32 = 0.001

func cmpFloat32(a, b float32) bool {
	return (a-b) < Epsilon && (b-a) < Epsilon
}

func TestMain(m *testing.M) {
	messagebirdtest.EnableServer(m)
}

func TestRead(t *testing.T) {
	messagebirdtest.WillReturn(balanceObject, http.StatusOK)
	client := messagebirdtest.Client(t)

	balance, err := Read(client)
	if err != nil {
		t.Fatalf("Didn't expect error while fetching the balance: %s", err)
	}

	if balance.Payment != "prepaid" {
		t.Errorf("Unexpected balance payment: %s", balance.Payment)
	}

	if balance.Type != "credits" {
		t.Errorf("Unexpected balance type: %s", balance.Type)
	}

	if !cmpFloat32(balance.Amount, 9.2) {
		t.Errorf("Unexpected balance amount: %.2f", balance.Amount)
	}
}

func TestReadError(t *testing.T) {
	messagebirdtest.WillReturnAccessKeyError()
	client := messagebirdtest.Client(t)

	_, err := Read(client)

	errorResponse, ok := err.(messagebird.ErrorResponse)
	if !ok {
		t.Fatalf("Expected ErrorResponse to be returned, instead I got %s", err)
	}

	if len(errorResponse.Errors) != 1 {
		t.Fatalf("Unexpected number of errors: %d, expected: 1", len(errorResponse.Errors))
	}

	if errorResponse.Errors[0].Code != 2 {
		t.Errorf("Unexpected error code: %d, expected: 2", errorResponse.Errors[0].Code)
	}

	if errorResponse.Errors[0].Parameter != "access_key" {
		t.Errorf("Unexpected error parameter: %s, expected: access_key", errorResponse.Errors[0].Parameter)
	}
}
