package signature

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testBaseUrl = "https://example.com"

var testSecret = []byte("hunter2")

func TestValidate(t *testing.T) {
	var cases = []struct {
		name            string
		signature       string
		signatureHeader string
		signatureKey    []byte
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
			signatureKey:    []byte(""),
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
	var cases = []struct {
		name           string
		requestParams  string
		requestPayload string
		receivedAt     string
		signature      string
		wantErr        string
	}{
		{
			name:       "valid with no params/body",
			receivedAt: "2021-07-05T12:00:00+02:00",
			signature:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiI1OWEyNDRkYy1lOWFkLTRlMjMtOTc3OC0zNzFmYWEyMzhmNzIiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDcifQ.SrhlKJ-ES4Dg8BBXKtop3u92Z_k4L4VjHKsyHWpweGE",
		},
		{
			name:          "valid with params and without body",
			requestParams: "/path?bar=1&foo=2",
			receivedAt:    "2021-07-05T12:00:00+02:00",
			signature:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiJjOTQ2YWY3Ny1lMTgyLTRlYWEtYjJmZi0xYTU0NWI1ZTk5MWEiLCJ1cmxfaGFzaCI6IjQxZjA1ZjBkZGQwYTIyYWIyMDlhYzQ2ZjQ3YzQ1NzJkOWNlZmEyNTdlZDc0YjI0MDA0YmFlNzUzZWNlNmMyNjAifQ.wUeGukU50HcPIr8d-zcCpttlGnPE-W57ujVb36AbAYw",
		},
		{
			name:           "valid with params and body",
			requestParams:  "/path?bar=1&foo=2",
			requestPayload: "Hello, World!",
			receivedAt:     "2021-07-05T12:00:00+02:00",
			signature:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiI5M2U1NTAwNi1hMGU4LTQ1MjYtYTE5MC1mYTVmZjAwZWExMTYiLCJ1cmxfaGFzaCI6IjQxZjA1ZjBkZGQwYTIyYWIyMDlhYzQ2ZjQ3YzQ1NzJkOWNlZmEyNTdlZDc0YjI0MDA0YmFlNzUzZWNlNmMyNjAiLCJwYXlsb2FkX2hhc2giOiJkZmZkNjAyMWJiMmJkNWIwYWY2NzYyOTA4MDllYzNhNTMxOTFkZDgxYzdmNzBhNGIyODY4OGEzNjIxODI5ODZmIn0.K6HyLDRdYgQBKN2tBcu0dOSxsfb_lOLaWby3un4rxIc",
		},
		{
			name:       "invalid token received before it is issued",
			receivedAt: "2021-07-05T12:00:00+02:00",
			signature:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ4MjgwMCwiZXhwIjoxNjI1NDgyODYwLCJqdGkiOiJmOWY4YzM4Mi0yNDQ5LTQzMTEtYjcyYi0xZGY3MTY4NzkzMWUiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDcifQ._59NNTg0j5YVXCRHgyeJAj8n6rTg1gwTh_I_coe7RDQ",
			wantErr:    "invalid jwt: iat is in the future",
		},
		{
			name:       "invalid token received after it is expired",
			receivedAt: "2021-07-05T12:00:00+02:00",
			signature:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3NTYwMCwiZXhwIjoxNjI1NDc1NjYwLCJqdGkiOiI1ZjAyZjUyMi02MDMwLTQ2YzgtYjVhMy0wMTI0NjQ3OGQ4YmMiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDcifQ.iGUCLsYVQG4iYWe2MkRoLQBBMzq7p_bLy4u0mhC3Jfc",
			wantErr:    "invalid jwt: exp is in the past",
		},
		{
			name:          "invalid token received on different URL",
			requestParams: "/path?bar=1&foo=2",
			receivedAt:    "2021-07-05T12:00:00+02:00",
			signature:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiJhNzVjOTA5Ni1lODIzLTQ0MmItODVmMi03ZDNjOWQ5YjcyNmIiLCJ1cmxfaGFzaCI6IjlmZGExZmNkYzc0YjEwMzUzNjhlNWY2NjhmNTdjOTFlOTk0MTJmZjU5Y2YwM2E0NmNlYjk1YWVhNWU2YjU4ZmQifQ.G4lpxrDOxZs75G1vIJ6J1jVbYS19tx2yq-lkIE-oETY",
			wantErr:       "invalid jwt: url_hash is invalid",
		},

		{
			name:           "invalid payload not match",
			requestPayload: "Hello, World!",
			receivedAt:     "2021-07-05T12:00:00+02:00",
			signature:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiIxNDUwMTUzMi05NmYyLTQ2ODQtOTgzMi02OGYwOTUxYWUzNDIiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDciLCJwYXlsb2FkX2hhc2giOiIzMjRjYzA2N2IyNTdlZGEwYmNiZDljOGQ4MTgwNzdhMDlhOTU2OGMwZDRjYTA2MDM4ZGVkOGZhZGRmODEzZmQ2In0.rQqiANogDOMafgg_B6p362PuhInAro9lMm2j_vruBA0",
			wantErr:        "invalid jwt: payload_hash is invalid",
		},

		{
			name:       "invalid signature key",
			receivedAt: "2021-07-05T12:00:00+02:00",
			signature:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiIyNDNjMjdhZS0yZjAyLTQ2YTAtODg1Mi1jNjZmMzdlYTlmNDYiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDcifQ._Uwf4HMtfAT6jvbBbh85Q9TunX0QlsXoaLGKX0I4VDg",
			wantErr:    "invalid jwt: signature is invalid",
		},

		{
			name:       "invalid missing payload",
			receivedAt: "2021-07-05T12:00:00+02:00",
			signature:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiIxNDUwMTUzMi05NmYyLTQ2ODQtOTgzMi02OGYwOTUxYWUzNDIiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDciLCJwYXlsb2FkX2hhc2giOiIzMjRjYzA2N2IyNTdlZGEwYmNiZDljOGQ4MTgwNzdhMDlhOTU2OGMwZDRjYTA2MDM4ZGVkOGZhZGRmODEzZmQ2In0.rQqiANogDOMafgg_B6p362PuhInAro9lMm2j_vruBA0",
			wantErr:    "invalid jwt: payload_hash was set; expected no payload value",
		},

		{
			name:           "invalid unexpected payload",
			requestPayload: "Hello, World!",
			receivedAt:     "2021-07-05T12:00:00+02:00",
			signature:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNZXNzYWdlQmlyZCIsImlhdCI6MTYyNTQ3OTIwMCwiZXhwIjoxNjI1NDc5MjYwLCJqdGkiOiI1OWEyNDRkYy1lOWFkLTRlMjMtOTc3OC0zNzFmYWEyMzhmNzIiLCJ1cmxfaGFzaCI6IjBmMTE1ZGIwNjJiN2MwZGQwMzBiMTY4NzhjOTlkZWE1YzM1NGI0OWRjMzdiMzhlYjg4NDYxNzljNzc4M2U5ZDcifQ.SrhlKJ-ES4Dg8BBXKtop3u92Z_k4L4VjHKsyHWpweGE",
			wantErr:        "invalid jwt: payload_hash is invalid",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			TimeFunc = func() time.Time {
				r, _ := time.Parse(time.RFC3339, test.receivedAt)
				return r
			}

			v := NewValidator(testSecret)
			reqUrl := testBaseUrl + test.requestParams
			if test.requestParams == "" {
				reqUrl += "/"
			}
			err := v.ValidSignature(test.signature, reqUrl, []byte(test.requestPayload))
			if test.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, test.wantErr)
		})
	}
}
