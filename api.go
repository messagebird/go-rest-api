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
	if cpr.Limit > 0 {
		query.Set("limit", strconv.Itoa(cpr.Limit))
	}
	if cpr.Offset >= 0 {
		query.Set("offset", strconv.Itoa(cpr.Offset))
	}

	return query.Encode()
}

// DefaultPagination provides reasonable values for List requests.
var DefaultPagination = &PaginationRequest{
	Limit:  20,
	Offset: 0,
}
