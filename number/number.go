package number

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	messagebird "github.com/messagebird/go-rest-api/v7"
)

const (
	// apiRoot is the absolute URL of the Numbers API.
	apiRoot = "https://numbers.messagebird.com/v1"

	// pathNumbers is the path for the Numbers resource, relative to apiRoot.
	// and path.
	pathNumbers = "phone-numbers"

	// pathNumbersAvailable is the path for the Search Number resource, relative to apiRoot.
	pathNumbersAvailable = "available-phone-numbers"
)

// Number represents a specific phone number.
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

// NumberList provide a list of all purchased phone numbers.
type NumberList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []*Number
}

// NumberSearchingList provide a list of all phone numbers.
// that are available for purchase.
type NumberSearchingList struct {
	Items []*Number
	Limit int
	Count int
}

// NumberListParams can be used to set query params in List().
type NumberListParams struct {
	Limit         int
	Offset        int
	Number        string
	Country       string
	Region        string
	Locality      string
	Features      []string
	Type          string
	Status        string
	SearchPattern NumberPattern
}

// NumberUpdateRequest can be used to set tags update.
type NumberUpdateRequest struct {
	Tags []string `json:"tags"`
}

// NumberPurchaseRequest can be used to purchase a number.
type NumberPurchaseRequest struct {
	Number                string `json:"number"`
	Country               string `json:"countryCode"`
	BillingIntervalMonths int    `json:"billingIntervalMonths"`
}

type NumberPattern string

const (
	// NumberPatternStart force phone numbers to start with the provided fragment.
	NumberPatternStart NumberPattern = "start"

	// NumberPatternEnd phone numbers can be somewhere within the provided fragment.
	NumberPatternEnd NumberPattern = "end"

	// NumberPatternAnyWhere force phone numbers to end with the provided fragment.
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

// Delete a purchased phone number
func Delete(c *messagebird.Client, phoneNumber string) error {
	uri := fmt.Sprintf("%s/%s", pathNumbers, phoneNumber)
	return request(c, nil, http.MethodDelete, uri, nil)
}

// Update updates a purchased phone number.
// Only updating *tags* is supported at the moment.
func Update(c *messagebird.Client, phoneNumber string, numberUpdateRequest *NumberUpdateRequest) (*Number, error) {
	uri := fmt.Sprintf("%s/%s", pathNumbers, phoneNumber)

	number := &Number{}
	if err := request(c, number, http.MethodPatch, uri, numberUpdateRequest); err != nil {
		return nil, err
	}

	return number, nil
}

// Purchases purchases a phone number.
func Purchase(c *messagebird.Client, numberPurchaseRequest *NumberPurchaseRequest) (*Number, error) {

	number := &Number{}
	if err := request(c, number, http.MethodPost, pathNumbers, numberPurchaseRequest); err != nil {
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

	if params.Type != "" {
		urlParams.Set("type", params.Type)
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
