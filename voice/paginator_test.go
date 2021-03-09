package voice

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginatorStream(t *testing.T) {
	type myStruct struct {
		Val int
	}
	mbClient, stop := testRequest(http.StatusOK, []byte(`{
		"data": [
			{ "Val": 1 },
			{ "Val": 2 },
			{ "Val": 3 }
		],
		"pagination": {
			"totalCount": 3,
			"pageCount": 1,
			"currentPage": 1,
			"perPage": 10
		}
	}`))
	defer stop()

	pag := newPaginator(mbClient, "", reflect.TypeOf(myStruct{}))

	i := 0
	for val := range pag.Stream() {
		t.Logf("%d, %v", i+1, val)
		v, ok := val.(myStruct)
		assert.True(t, ok)
		assert.Equal(t, i+1, v.Val)
		i++
	}
	assert.Equal(t, 3, i)
}
