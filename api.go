package messagebird

import (
	"errors"
	"log"
	"net/url"
	"reflect"
	"strconv"
)

type PaginationRequest interface {
	QueryParams() string
}

// CommonPaginationRequest can be used to set pagination options in List().
type CommonPaginationRequest struct {
	Limit, Offset int
}

func (cpr *CommonPaginationRequest) QueryParams() string {
	if cpr == nil {
		return ""
	}

	query := url.Values{}
	query.Set("limit", strconv.Itoa(cpr.Limit))
	query.Set("offset", strconv.Itoa(cpr.Offset))

	return query.Encode()
}

// DefaultPagination provides reasonable values for List requests.
var DefaultPagination = &CommonPaginationRequest{
	Limit:  20,
	Offset: 0,
}

func MakeQueryParams(s interface{}) (string, error) {
	rt := reflect.TypeOf(s)

	if rt.Kind() != reflect.Struct && rt.Elem().Kind() != reflect.Struct {
		return "", errors.New("unexpected kind of value, expected strcut")
	}

	//query := url.Values{}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		log.Fatalln(field.Tag)
		//v := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		//if v == tag {
		//	return f.Name
		//}
	}
	return "", nil
}
