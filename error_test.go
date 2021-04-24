package messagebird

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		errRes := ErrorResponse{
			Errors: []Error{
				{
					Code:        42,
					Description: "something bad",
					Parameter:   "foo",
				},
			},
		}
		assert.Error(t, errRes)
	})

	t.Run("Multiple", func(t *testing.T) {
		errRes := ErrorResponse{
			Errors: []Error{
				{
					Code:        42,
					Description: "something bad",
					Parameter:   "foo",
				},
				{
					Code:        42,
					Description: "something else",
					Parameter:   "foo",
				},
			},
		}
		assert.Error(t, errRes)
	})
}

func TestError_ErrorMessageContainsOnlyDescription(t *testing.T) {
	err := Error{
		Code:        42,
		Description: "something bad",
		Parameter:   "foo",
	}

	assert.Equal(t, "something bad", err.Error())
}

func TestErrorResponse_ErrorMessageFormat(t *testing.T) {
	err1 := Error{
		Code:        42,
		Description: "something bad",
		Parameter:   "foo",
	}

	err2 := Error{
		Code:        43,
		Description: "something else",
		Parameter:   "bar",
	}

	cases := map[string]struct {
		errors []Error
		exp    string
	}{
		"nil errors slice": {
			errors: nil,
			exp:    "API errors: ",
		},
		"one error": {
			errors: []Error{
				err1,
			},
			exp: "API errors: " + err1.Error(),
		},
		"two errors separated with coma": {
			errors: []Error{
				err1,
				err2,
			},
			exp: "API errors: " + err1.Error() + ", " + err2.Error(),
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			err := ErrorResponse{
				Errors: c.errors,
			}

			assert.Equal(t, c.exp, err.Error())
		})
	}

}
