package messagebird

import (
	"fmt"
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

func BenchmarkErrorResponse_Error(b *testing.B) {
	for n := 1; n <= 1024; n *= 2 {
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			errors := make([]Error, n)

			for i := 0; i < b.N; i++ {
				err := ErrorResponse{Errors: errors}
				_ = err.Error()
			}
		})
	}
}
