package voice

import (
	"net/http"
	"reflect"
	"testing"
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
		if v, ok := val.(myStruct); !ok || v.Val != i+1 {
			t.Fatalf("unexpected item at index %d: %#v", i, val)
		}
		i++
	}
	if i != 3 {
		t.Fatalf("unexpected number of elements: %d", i)
	}
}
