package number

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	messagebird "github.com/messagebird/go-rest-api"
)

const (
	apiRoot              = "https://numbers.messagebird.com/v1"
	pathNumbers          = "phone-numbers"
	pathNumbersAvailable = "available-phone-numbers"
)

type Number struct {
	Number   string
	Country  string
	Region   string
	Locality string
	Features []string
	Tags     []string
	Type     string
	Status   string
}

type NumberList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []*Number
}

type NumberSearchingList struct {
	Items []*Number
	Limit int
	Count int
}

type NumberListParams struct {
	Limit         int
	Offset        int
	Number        string
	Country       string
	Region        string
	Locality      string
	Features      []string
	Type          []string
	Status        string
	SearchPattern NumberPattern
}

type NumberPattern string

const (
	NumberPatternStart    NumberPattern = "start"
	NumberPatternEnd      NumberPattern = "end"
	NumberPatternAnyWhere NumberPattern = "anywhere"
)

// request does the exact same thing as Client.Request. It does, however,
// prefix the path with the Numbers API's root. This ensures the client
// doesn't "handle" this for us: by default, it uses the REST API.
func request(c *messagebird.Client, v interface{}, method, path string, data interface{}) error {
	return c.Request(v, method, fmt.Sprintf("%s/%s", apiRoot, path), data)
}

// List get all purchased phone numbers
func List(c *messagebird.Client, listParams *NumberListParams) (*NumberList, error) {
	uri := getpath(listParams, pathNumbers)

	numberList := &NumberList{}
	if err := request(c, numberList, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}
	return numberList, nil
}

// Search for phone numbers available for purchase, countryCode needs to be in Alpha-2 country code (example: NL)
func Search(c *messagebird.Client, countryCode string, listParams *NumberListParams) (*NumberSearchingList, error) {
	uri := getpath(listParams, pathNumbersAvailable+"/"+countryCode)

	numberList := &NumberSearchingList{}
	if err := request(c, numberList, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return numberList, nil
}

// Read get a purchased phone number
func Read(c *messagebird.Client, phoneNumber string) (*Number, error) {
	if len(phoneNumber) < 5 {
		return nil, fmt.Errorf("a phoneNumber is too short")
	}

	uri := fmt.Sprintf("%s/%s", pathNumbers, phoneNumber)

	number := &Number{}
	if err := request(c, number, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return number, nil
}

// GetPath get the full path for the request
func getpath(listParams *NumberListParams, path string) string {
	params := paramsForMessageList(listParams)
	return fmt.Sprintf("%s?%s", path, params.Encode())
}

// paramsForMessageList build query params
func paramsForMessageList(params *NumberListParams) *url.Values {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams
	}

	if len(params.Features) > 0 {
		paramsForArrays("features", params.Features, urlParams)
	}

	if len(params.Type) > 0 {
		paramsForArrays("type", params.Type, urlParams)
	}

	if params.Number != "" {
		urlParams.Set("number", params.Number)
	}
	if params.Country != "" {
		urlParams.Set("country", params.Country)
	}
	if params.Limit != 0 {
		urlParams.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.SearchPattern != "" {
		urlParams.Set("search_pattern", string(params.SearchPattern))
	}

	if params.Offset != 0 {
		urlParams.Set("offset", strconv.Itoa(params.Offset))
	}

	return urlParams
}

// paramsForArrays build query for array params
func paramsForArrays(field string, values []string, urlParams *url.Values) {
	for _, value := range values {
		urlParams.Add(field, value)
	}
}
