package voice

import (
	"testing"

	"github.com/messagebird/go-rest-api/internal/mbtest"
)

func TestErrorReader(t *testing.T) {
	t.Run("Single error", func(t *testing.T) {
		b := mbtest.Testdata(t, "error.json")
		err := errorReader(b).(ErrorResponse)

		if count := len(err.Errors); count != 1 {
			t.Fatalf("Got %d, expected 1", count)
		}

		if err.Errors[0].Code != 13 {
			t.Errorf("Got %d, expected 13", err.Errors[0].Code)
		}
		if err.Errors[0].Message != "some-error" {
			t.Errorf("Got %q, expected some-error", err.Errors[0].Message)
		}
	})

	t.Run("Multiple errors", func(t *testing.T) {
		b := mbtest.Testdata(t, "errors.json")
		err := errorReader(b).(ErrorResponse)

		if count := len(err.Errors); count != 2 {
			t.Fatalf("Got %d, expected 2", count)
		}

		if err.Errors[0].Code != 11 {
			t.Errorf("Got %d, expected 11", err.Errors[0].Code)
		}
		if err.Errors[0].Message != "some-error" {
			t.Errorf("Got %q, expected some-error", err.Errors[0].Message)
		}
		if err.Errors[1].Code != 15 {
			t.Errorf("Got %d, expected 15", err.Errors[1].Code)
		}
		if err.Errors[1].Message != "other-error" {
			t.Errorf("Got %q, expected other-error", err.Errors[1].Message)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		b := []byte("clearly not json")
		_, ok := errorReader(b).(ErrorResponse)

		if ok {
			// If the data b is not JSON, we expect a "generic" errorString
			// (from fmt.Errorf), but we somehow got our own ErrorResponse back.
			t.Fatalf("Got ErrorResponse, expected errorString")
		}
	})
}

func TestErrorResponseError(t *testing.T) {
	err := ErrorResponse{
		[]Error{
			{
				Code: 1,
				Message: "foo",
			},
			{
				Code: 2,
				Message: "bar",
			},
		},
	}

	expect := `code: 1, message: "foo"; code: 2, message: "bar"`
	if actual := err.Error(); actual != expect {
		t.Fatalf("Got %q, expected %q", actual, expect)
	}
}
