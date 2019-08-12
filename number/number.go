package number

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
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
	SearchPattern string
}

type NumberUpdateRequest struct {
	Tags []string `json:"tags"`
}

type NumberPurchaseRequest struct {
	Number                string `json:"number"`
	Country               string `json:"countryCode"`
	BillingIntervalMonths int    `json:"billingIntervalMonths"`
}

// request does the exact same thing as Client.Request. It does, however,
// prefix the path with the Numbers API's root. This ensures the client
// doesn't "handle" this for us: by default, it uses the REST API.
func request(c *messagebird.Client, v interface{}, method, path string, data interface{}) error {
	fmt.Println(path)
	return c.Request(v, method, fmt.Sprintf("%s/%s", apiRoot, path), data)
}

// List get all purchased phone numbers
func List(c *messagebird.Client, listParams *NumberListParams) (*NumberList, error) {
	uri := GetPath(listParams, pathNumbers)

	numberList := &NumberList{}
	if err := request(c, numberList, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}
	return numberList, nil
}

// Search for phone numbers available for purchase
func Search(c *messagebird.Client, cc string, listParams *NumberListParams) (*NumberSearchingList, error) {
	uri := GetPath(listParams, pathNumbersAvailable+"/"+cc)

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

// Delete cancels a purchased phone number
func Delete(c *messagebird.Client, phoneNumber string) (error) {
	uri := fmt.Sprintf("%s/%s", pathNumbers, phoneNumber)
	return request(c, nil, http.MethodDelete, uri, nil)
}

// Update updates a purchased phone number
// Only updating *tags* is supported at the moment
func Update(c *messagebird.Client, phoneNumber string, numberUpdateRequest *NumberUpdateRequest) (*Number, error) {
	uri := fmt.Sprintf("%s/%s", pathNumbers, phoneNumber)

	number := &Number{}
	if err := request(c, number, http.MethodPatch, uri, numberUpdateRequest); err != nil {
		return nil, err
	}

	return number, nil
}

// Create purchases a phone number
func Create(c *messagebird.Client, numberPurchaseRequest *NumberPurchaseRequest) (*Number, error) {

	number := &Number{}
	if err := request(c, number, http.MethodPost, pathNumbers, numberPurchaseRequest); err != nil {
		return nil, err
	}

	return number, nil
}

// GetPath get the full path for the request
func GetPath(listParams *NumberListParams, path string) string {
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
		paramsForArrays("features", "^(sms|voice|mms)$", params.Features, urlParams)
	}

	if len(params.Type) > 0 {
		paramsForArrays("type", "^(mobile|mobile|premium_rate)$", params.Type, urlParams)
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
		urlParams.Set("search_pattern", params.SearchPattern)
	}

	if params.Offset != 0 {
		urlParams.Set("offset", strconv.Itoa(params.Offset))
	}

	return urlParams
}

// paramsForArrays build query for array params
func paramsForArrays(field string, pattern string, array []string, urlParams *url.Values) {
	r, _ := regexp.Compile(pattern)

	for i := 0; i < len(array); i++ {
		if match := r.MatchString(array[i]); match {
			urlParams.Add(field, array[i])
		}
	}
}
