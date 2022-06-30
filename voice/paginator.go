package voice

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	messagebird "github.com/messagebird/go-rest-api/v8"
)

// A Paginator is used to stream the contents of a collection of some type from
// the MessageBird API.
//
// Paginators are single use and can therefore not be reset.
type Paginator struct {
	endpoint   string
	nextPage   int
	structType reflect.Type
	client     messagebird.MessageBirdClient
}

// newPaginator creates a new paginator.
//
// endpoint is called with the `page` query parameter until no more pages are
// available.
//
// typ is the non-pointer type of a single element returned by a page.
func newPaginator(client messagebird.MessageBirdClient, endpoint string, typ reflect.Type) *Paginator {
	return &Paginator{
		endpoint:   endpoint,
		nextPage:   1, // Page indices start at 1.
		structType: typ,
		client:     client,
	}
}

// NextPage queries the next page from the MessageBird API.
//
// The interface{} contains a slice of the type this paginator handles.
//
// When no more items are available, an empty slice and io.EOF are returned.
// If another kind of error occurs, nil and and the error are returned.
func (pag *Paginator) NextPage() (interface{}, error) {
	type pagination struct {
		TotalCount  int `json:"totalCount"`
		PageCount   int `json:"pageCount"`
		CurrentPage int `json:"currentPage"`
		PerPage     int `json:"perPage"`
	}
	rawType := reflect.StructOf([]reflect.StructField{
		{
			Name: "Data",
			Type: reflect.SliceOf(pag.structType),
			Tag:  "json:\"data\"",
		},
		{
			Name: "Pagination",
			Type: reflect.TypeOf(pagination{}),
			Tag:  "json:\"pagination\"",
		},
	})
	rawVal := reflect.New(rawType)

	if err := pag.client.Request(rawVal.Interface(), http.MethodGet, fmt.Sprintf("%s?page=%d", pag.endpoint, pag.nextPage), nil); err != nil {
		return nil, err
	}

	data := rawVal.Elem().FieldByName("Data").Interface()
	pageInfo := rawVal.Elem().FieldByName("Pagination").Interface().(pagination)

	// If no more items are available, a page with 0 elements is returned.
	if pag.nextPage > pageInfo.PageCount {
		return data, io.EOF
	}

	pag.nextPage++
	return data, nil
}

// Stream creates a channel which streams the contents of all remaining pages
// ony by one.
//
// The Paginator is consumed in the process, meaning that after elements have
// been received, NextPage will return EOF. It is invalid to mix calls to
// NextPage an Stream, even after the stream channel was closed.
//
// If an error occurs, the next item sent over the channel will be an error
// instead of a regular value. The channel is closed directly after this.
func (pag *Paginator) Stream() <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			page, err := pag.NextPage()
			if err != nil {
				if err != io.EOF {
					out <- err
				}
				break
			}
			v := reflect.ValueOf(page)
			for i, l := 0, v.Len(); i < l; i++ {
				out <- v.Index(i).Interface()
			}
		}
	}()
	return out
}
