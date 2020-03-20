package messagebird

import (
	"testing"
)

func TestError(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		errRes := ErrorResponse{
			Errors: []Error{
				Error{
					Code:        42,
					Description: "something bad",
					Parameter:   "foo",
				},
			},
		}
		if s := errRes.Error(); s != "API errors: something bad" {
			t.Errorf("Got %q, expected API response: something bad", s)
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		errRes := ErrorResponse{
			Errors: []Error{
				Error{
					Code:        42,
					Description: "something bad",
					Parameter:   "foo",
				},
				Error{
					Code:        42,
					Description: "something else",
					Parameter:   "foo",
				},
			},
		}
		if s := errRes.Error(); s != "API errors: something bad, something else" {
			t.Errorf("Got %q, expected API response: something bad, something else", s)
		}
	})
}
