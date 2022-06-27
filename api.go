package messagebird

import (
	"net/url"
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
