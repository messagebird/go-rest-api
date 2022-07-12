package voice

import (
	"testing"

	"github.com/messagebird/go-rest-api/v9/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestErrorReader(t *testing.T) {
	t.Run("Single error", func(t *testing.T) {
		b := mbtest.Testdata(t, "error.json")
		err := errorReader(b).(ErrorResponse)

		assert.Len(t, err.Errors, 1)
		assert.Equal(t, 13, err.Errors[0].Code)
		assert.Equal(t, "some-error", err.Errors[0].Message)
	})

	t.Run("Multiple errors", func(t *testing.T) {
		b := mbtest.Testdata(t, "errors.json")
		err := errorReader(b).(ErrorResponse)

		assert.Len(t, err.Errors, 2)
		assert.Equal(t, 11, err.Errors[0].Code)
		assert.Equal(t, "some-error", err.Errors[0].Message)
		assert.Equal(t, 15, err.Errors[1].Code)
		assert.Equal(t, "other-error", err.Errors[1].Message)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		b := []byte("clearly not json")
		_, ok := errorReader(b).(ErrorResponse)

		assert.False(t, ok)
	})
}

func TestErrorResponseError(t *testing.T) {
	err := ErrorResponse{
		[]Error{
			{
				Code:    1,
				Message: "foo",
			},
			{
				Code:    2,
				Message: "bar",
			},
		},
	}

	expect := `code: 1, message: "foo"; code: 2, message: "bar"`
	actual := err.Error()
	assert.Equal(t, expect, actual)
}
