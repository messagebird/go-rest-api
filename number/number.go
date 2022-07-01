package number

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	messagebird "github.com/messagebird/go-rest-api/v9"
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

// Numbers provide a list of all purchased phone numbers.
type Numbers struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Items      []*Number
}

// NumbersSearching provide a list of all phone numbers.
// that are available for purchase.
type NumbersSearching struct {
	Items []*Number
	Limit int
	Count int
}

// ListRequest can be used to set query params in List().
type ListRequest struct {
	Limit                             int
	Offset                            int
	Number                            string
	Country                           string
	Region                            string
	Locality                          string
	Features                          []string // Possible values: sms, voice, mms.
	Type                              string   // Possible values: landline, mobile, premium_rate and toll_free.
	Status                            string
	SearchPattern                     SearchPattern
	ExcludeNumbersRequireVerification bool // exclude_numbers_require_verification
	Prices                            bool // exclude_numbers_require_verification
}

func (lr *ListRequest) QueryParams() string {
	if lr == nil {
		return ""
	}

	query := url.Values{}

	if len(lr.Features) > 0 {
		paramsForArrays("features", lr.Features, &query)
	}

	if len(lr.Type) > 0 {
		query.Set("type", lr.Type)
	}

	if len(lr.Number) > 0 {
		query.Set("number", lr.Number)
	}
	if len(lr.Country) > 0 {
		query.Set("country", lr.Country)
	}
	if lr.Limit > 0 {
		query.Set("limit", strconv.Itoa(lr.Limit))
	}
	if lr.Offset > 0 {
		query.Set("offset", strconv.Itoa(lr.Offset))
	}
	if lr.SearchPattern != "" {
		query.Set("search_pattern", string(lr.SearchPattern))
	}

	if lr.Offset != 0 {
		query.Set("offset", strconv.Itoa(lr.Offset))
	}

	return query.Encode()
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

type SearchPattern string

const (
	// SearchPatternStart force phone numbers to start with the provided fragment.
	SearchPatternStart SearchPattern = "start"

	// SearchPatternEnd phone numbers can be somewhere within the provided fragment.
	SearchPatternEnd SearchPattern = "end"

	// SearchPatternAnyWhere force phone numbers to end with the provided fragment.
	SearchPatternAnyWhere SearchPattern = "anywhere"
)

// List fetch all purchased phone numbers
func List(c messagebird.MessageBirdClient, listParams *ListRequest) (*Numbers, error) {
	uri := getpath(listParams, pathNumbers)

	numberList := &Numbers{}
	if err := request(c, numberList, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}
	return numberList, nil
}

// Search for phone numbers available for purchase, countryCode needs to be in Alpha-2 country code (example: NL)
func Search(c messagebird.MessageBirdClient, countryCode string, listParams *ListRequest) (*NumbersSearching, error) {
	uri := getpath(listParams, pathNumbersAvailable+"/"+countryCode)

	numberList := &NumbersSearching{}
	if err := request(c, numberList, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return numberList, nil
}

// Read get a purchased phone number
func Read(c messagebird.MessageBirdClient, phoneNumber string) (*Number, error) {
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
func Delete(c messagebird.MessageBirdClient, phoneNumber string) error {
	uri := fmt.Sprintf("%s/%s", pathNumbers, phoneNumber)
	return request(c, nil, http.MethodDelete, uri, nil)
}

// Update updates a purchased phone number.
// Only updating *tags* is supported at the moment.
func Update(c messagebird.MessageBirdClient, phoneNumber string, numberUpdateRequest *NumberUpdateRequest) (*Number, error) {
	uri := fmt.Sprintf("%s/%s", pathNumbers, phoneNumber)

	number := &Number{}
	if err := request(c, number, http.MethodPatch, uri, numberUpdateRequest); err != nil {
		return nil, err
	}

	return number, nil
}

// Purchases purchases a phone number.
func Purchase(c messagebird.MessageBirdClient, numberPurchaseRequest *NumberPurchaseRequest) (*Number, error) {

	number := &Number{}
	if err := request(c, number, http.MethodPost, pathNumbers, numberPurchaseRequest); err != nil {
		return nil, err
	}

	return number, nil
}

// GetPath get the full path for the request
func getpath(listParams *ListRequest, path string) string {
	params := paramsForMessageList(listParams)
	return fmt.Sprintf("%s?%s", path, params.Encode())
}

// paramsForMessageList build query params
func paramsForMessageList(params *ListRequest) *url.Values {
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

// request does the exact same thing as Client.Request. It does, however,
// prefix the path with the Numbers API's root. This ensures the client
// doesn't "handle" this for us: by default, it uses the REST API.
func request(c messagebird.MessageBirdClient, v interface{}, method, path string, data interface{}) error {
	return c.Request(v, method, fmt.Sprintf("%s/%s", apiRoot, path), data)
}
