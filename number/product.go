package number

import (
	"fmt"
	messagebird "github.com/messagebird/go-rest-api/v9"
	"net/http"
	"net/url"
	"strconv"
)

type Product struct {
	Country                     string
	NumberType                  string
	BackOrderLeadTime           string
	ReachableFromNationalFixed  bool
	ReachableFromNationalMobile bool
	ReachableFromPayPhone       bool
	VerificationRequired        bool
	InitialContractDuration     string
	Prefixes                    []*Prefix
	Remarks                     []string
	Conditions                  []string
	EndUserData                 []string
	ForbiddenContent            []string
}

type Prefix struct {
	Prefix     string
	City       string
	StateProv  string
	PrefixType string
}

type ShortProduct struct {
	NumberType           string
	VerificationRequired bool
	Country              string
	ID                   int
	Currency             string
	Price                int
}

type Products struct {
	Items []*ShortProduct
	Count int
	Limit int
}

// ProductsRequest can be used to set query params in SearchProducts().
type ProductsRequest struct {
	CountryCode string   `json:"countryCode"`
	Limit       int      `json:"limit"`
	Features    []string `json:"features"`
	Type        string   `json:"type"`
	Prefix      string   `json:"prefix"`
}

func (req *ProductsRequest) QueryParams() string {
	if req == nil {
		return ""
	}

	query := url.Values{}

	if len(req.Features) > 0 {
		paramsForArrays("features", req.Features, &query)
	}

	if req.Limit > 0 {
		query.Set("limit", strconv.Itoa(req.Limit))
	}

	if len(req.Type) > 0 {
		query.Set("type", req.Type)
	}

	if len(req.Prefix) > 0 {
		query.Set("prefix", req.Prefix)
	}

	return query.Encode()
}

// SearchProducts searches for unified communication phone numbers that are available for you to back order.
func SearchProducts(c messagebird.MessageBirdClient, params *ProductsRequest) (*Products, error) {
	uri := fmt.Sprintf("%s?%s", pathProducts, params.QueryParams())

	pr := &Products{}
	if err := request(c, pr, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return pr, nil
}

// ReadProduct get a purchased phone number
func ReadProduct(c messagebird.MessageBirdClient, productID string) (*Product, error) {
	uri := fmt.Sprintf("%s/%s", pathProducts, productID)

	pr := &Product{}
	if err := request(c, pr, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return pr, nil
}
