package number

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v9"
)

const (
	// apiRoot is the absolute URL of the Numbers API.
	apiRoot = "https://numbers.messagebird.com/v1"

	// pathPhoneNumbers is the path for the Numbers resource, relative to apiRoot.
	// and path.
	pathPhoneNumbers = "phone-numbers"

	pathNumbers = "numbers"

	pathProducts = "products"

	pathBackorders = "backorders"

	pathDocuments = "documents"

	pathPools = "pools"

	pathEndUserDetails = "end-user-details"

	// pathNumbersAvailable is the path for the Search Number resource, relative to apiRoot.
	pathNumbersAvailable = "available-phone-numbers"
)

type SearchPattern string

const (
	// SearchPatternStart force phone numbers to start with the provided fragment.
	SearchPatternStart SearchPattern = "start"

	// SearchPatternEnd phone numbers can be somewhere within the provided fragment.
	SearchPatternEnd SearchPattern = "end"

	// SearchPatternAnyWhere force phone numbers to end with the provided fragment.
	SearchPatternAnyWhere SearchPattern = "anywhere"
)

type Type string

const (
	TypeLandline    Type = "landline"
	TypeMobile      Type = "mobile"
	TypePremiumRate Type = "premium_rate"
	TypeTollFree    Type = "toll_free"
)

type Feature string

const (
	FeatureSMS   Feature = "sms"
	FeatureVoice Feature = "voice"
	FeatureMMS   Feature = "mms"
)

// Number represents a specific phone number.
type Number struct {
	Number                  string
	Country                 string
	Region                  string
	Locality                string
	Features                []Feature
	Tags                    []string
	Type                    Type
	Status                  string
	VerificationRequired    bool
	InitialContractDuration int
	InboundCallsOnly        bool
	MonthlyPrice            float64
	Currency                string
	Conditions              []string
	CreatedAt               *time.Time
	RenewalAt               *time.Time
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
	Limit    int
	Offset   int
	Features []string // Possible values: sms, voice, mms.
	Tags     []string
	Number   string
	Region   string
	Locality string
	Type     string // Possible values: landline, mobile, premium_rate and toll_free.
}

func (lr *ListRequest) QueryParams() string {
	if lr == nil {
		return ""
	}

	query := url.Values{}

	if len(lr.Features) > 0 {
		paramsForArrays("features", lr.Features, &query)
	}

	if len(lr.Tags) > 0 {
		paramsForArrays("tags", lr.Tags, &query)
	}

	if lr.Limit != 0 {
		query.Set("limit", strconv.Itoa(lr.Limit))
	}

	if lr.Offset != 0 {
		query.Set("offset", strconv.Itoa(lr.Offset))
	}

	if lr.Type != "" {
		query.Set("type", lr.Type)
	}

	if lr.Locality != "" {
		query.Set("locality", lr.Locality)
	}

	if lr.Number != "" {
		query.Set("number", lr.Number)
	}

	if lr.Region != "" {
		query.Set("region", lr.Region)
	}

	return query.Encode()
}

// SearchRequest can be used to set query params in Search().
type SearchRequest struct {
	Limit                             int
	Offset                            int
	Number                            string
	Country                           string
	Region                            string
	Locality                          string
	Features                          []string // Possible values: sms, voice, mms.
	Tags                              []string
	Type                              string // Possible values: landline, mobile, premium_rate and toll_free.
	Status                            string
	SearchPattern                     SearchPattern
	ExcludeNumbersRequireVerification bool // exclude_numbers_require_verification
	Prices                            bool
}

func (sr *SearchRequest) QueryParams() string {
	if sr == nil {
		return ""
	}

	query := url.Values{}

	if len(sr.Features) > 0 {
		paramsForArrays("features", sr.Features, &query)
	}

	if len(sr.Tags) > 0 {
		paramsForArrays("tags", sr.Tags, &query)
	}

	if sr.Limit > 0 {
		query.Set("limit", strconv.Itoa(sr.Limit))
	}
	if sr.Offset > 0 {
		query.Set("offset", strconv.Itoa(sr.Offset))
	}

	if len(sr.Type) > 0 {
		query.Set("type", sr.Type)
	}

	if len(sr.Number) > 0 {
		query.Set("number", sr.Number)
	}
	if len(sr.Country) > 0 {
		query.Set("country", sr.Country)
	}
	if len(sr.Region) > 0 {
		query.Set("region", sr.Region)
	}
	if len(sr.Locality) > 0 {
		query.Set("locality", sr.Locality)
	}
	if len(sr.Status) > 0 {
		query.Set("status", sr.Status)
	}
	query.Set("exclude_numbers_require_verification", strconv.FormatBool(sr.ExcludeNumbersRequireVerification))
	query.Set("prices", strconv.FormatBool(sr.Prices))

	if sr.SearchPattern != "" {
		query.Set("search_pattern", string(sr.SearchPattern))
	}

	return query.Encode()
}

// UpdateRequest can be used to set tags update.
type UpdateRequest struct {
	Tags []string `json:"tags"`
}

// PurchaseRequest can be used to purchase a number.
type PurchaseRequest struct {
	Number                string `json:"number"`
	Country               string `json:"countryCode"`
	BillingIntervalMonths int    `json:"billingIntervalMonths"`
}

// List fetches all purchased phone numbers
func List(c messagebird.Client, params *ListRequest) (*Numbers, error) {
	uri := fmt.Sprintf("%s?%s", pathPhoneNumbers, params.QueryParams())

	numberList := &Numbers{}
	if err := request(c, numberList, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return numberList, nil
}

// Search for phone numbers available for purchase, countryCode needs to be in Alpha-2 country code (example: NL)
func Search(c messagebird.Client, countryCode string, params *SearchRequest) (*NumbersSearching, error) {
	uri := fmt.Sprintf("%s/%s?%s", pathNumbersAvailable, countryCode, params.QueryParams())

	numberList := &NumbersSearching{}
	if err := request(c, numberList, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return numberList, nil
}

// Read get a purchased phone number
func Read(c messagebird.Client, phoneNumber string) (*Number, error) {
	if len(phoneNumber) < 5 {
		return nil, fmt.Errorf("a phoneNumber is too short")
	}

	uri := fmt.Sprintf("%s/%s", pathPhoneNumbers, phoneNumber)

	number := &Number{}
	if err := request(c, number, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return number, nil
}

// Delete a purchased phone number
func Delete(c messagebird.Client, phoneNumber string) error {
	uri := fmt.Sprintf("%s/%s", pathPhoneNumbers, phoneNumber)
	return request(c, nil, http.MethodDelete, uri, nil)
}

// Update updates a purchased phone number.
// Only updating *tags* is supported at the moment.
func Update(c messagebird.Client, phoneNumber string, req *UpdateRequest) (*Number, error) {
	uri := fmt.Sprintf("%s/%s", pathPhoneNumbers, phoneNumber)

	number := &Number{}
	if err := request(c, number, http.MethodPatch, uri, req); err != nil {
		return nil, err
	}

	return number, nil
}

// Purchase purchases a phone number.
func Purchase(c messagebird.Client, numberPurchaseRequest *PurchaseRequest) (*Number, error) {
	number := &Number{}
	if err := request(c, number, http.MethodPost, pathPhoneNumbers, numberPurchaseRequest); err != nil {
		return nil, err
	}

	return number, nil
}

// paramsForArrays build query for array params
func paramsForArrays(field string, values []string, urlParams *url.Values) {
	for _, value := range values {
		urlParams.Add(field, value)
	}
}

// request does the exact same thing as DefaultClient.Request. It does, however,
// prefix the path with the Numbers API's root. This ensures the client
// doesn't "handle" this for us: by default, it uses the REST API.
func request(c messagebird.Client, v interface{}, method, path string, data interface{}) error {
	return c.Request(v, method, fmt.Sprintf("%s/%s", apiRoot, path), data)
}
