package signature_jwt

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/v9/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	var testBaseUrl = "https://example.com"
	var testSecret = "hunter2"

	var cases = []struct {
		name            string
		signature       string
		signatureHeader string
		signatureKey    string
		receivedAt      string
		wantCode        int
	}{
		{
			name:            "valid request",
			signature:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiI1OWEyNDRkYy1lOWFkLTRlMjMtOTc3OC0zNzFmYWEyMzhmNzIiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDcifQ.SrhlKJ-ES4Dg8BBXKtop3u92Z_k4L4VjHKsyHWpweGE",
			signatureHeader: signatureHeader,
			signatureKey:    testSecret,
			receivedAt:      "2021-07-05T12:00:00+02:00",
			wantCode:        http.StatusOK,
		},
		{
			name:            "empty signature key",
			signature:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiI1OWEyNDRkYy1lOWFkLTRlMjMtOTc3OC0zNzFmYWEyMzhmNzIiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDcifQ.SrhlKJ-ES4Dg8BBXKtop3u92Z_k4L4VjHKsyHWpweGE",
			signatureHeader: signatureHeader,
			signatureKey:    "",
			receivedAt:      "2021-07-05T12:00:00+02:00",
			wantCode:        http.StatusUnauthorized,
		},
		{
			name:            "empty signature",
			signature:       "",
			signatureHeader: signatureHeader,
			signatureKey:    testSecret,
			receivedAt:      "2021-07-05T12:00:00+02:00",
			wantCode:        http.StatusUnauthorized,
		},

		{
			name:            "wrong signature header",
			signature:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiI1OWEyNDRkYy1lOWFkLTRlMjMtOTc3OC0zNzFmYWEyMzhmNzIiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDcifQ.SrhlKJ-ES4Dg8BBXKtop3u92Z_k4L4VjHKsyHWpweGE",
			signatureHeader: "Wrong-Header",
			signatureKey:    testSecret,
			receivedAt:      "2021-07-05T12:00:00+02:00",
			wantCode:        http.StatusUnauthorized,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			TimeFunc = func() time.Time {
				r, _ := time.Parse(time.RFC3339, test.receivedAt)
				return r
			}

			v := NewValidator(test.signatureKey)
			ts := httptest.NewServer(v.Validate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}), testBaseUrl))
			defer ts.Close()

			client := &http.Client{}
			req, _ := http.NewRequest("GET", ts.URL, nil)
			req.Header.Set(test.signatureHeader, test.signature)

			res, _ := client.Do(req)

			assert.Equal(t, test.wantCode, res.StatusCode)

			TimeFunc = time.Now
		})
	}
}

func TestValidSignature(t *testing.T) {
	testData := mbtest.Testdata(t, "reference.json")

	var tcs []struct {
		Name      string `json:"name"`
		Method    string `json:"method"`
		Secret    string `json:"secret"`
		Url       string `json:"url"`
		Payload   string `json:"payload"`
		Timestamp string `json:"timestamp"`
		Token     string `json:"token"`
		Valid     bool   `json:"valid"`
		Reason    string `json:"reason"`
	}
	if err := json.Unmarshal(testData, &tcs); err != nil {
		assert.NoError(t, err)
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			TimeFunc = func() time.Time {
				r, _ := time.Parse(time.RFC3339, tc.Timestamp)
				return r
			}

			v := NewValidator(tc.Secret)
			claims, err := v.ValidateSignature(tc.Token, tc.Url, []byte(tc.Payload))
			if tc.Valid {
				assert.NoError(t, err)
				assert.NotEmpty(t, claims)
				return
			}
			assert.EqualError(t, err, tc.Reason)
		})
	}
}
