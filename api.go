package messagebird

import (
	"net/url"
	"strconv"
)

// PaginationRequest can be used to set pagination options in List().
type PaginationRequest struct {
	Limit, Offset int
}

func (cpr *PaginationRequest) QueryParams() string {
	if cpr == nil {
		return ""
	}

	query := url.Values{}
	query.Set("limit", strconv.Itoa(cpr.Limit))
	query.Set("offset", strconv.Itoa(cpr.Offset))

	return query.Encode()
}

// DefaultPagination provides reasonable values for List requests.
var DefaultPagination = &PaginationRequest{
	Limit:  20,
	Offset: 0,
}
