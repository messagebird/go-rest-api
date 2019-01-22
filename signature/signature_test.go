package signature

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const testTs = "1544544948"
const testQp = "abc=foo&def=bar"
const testBody = `{"a key":"some value"}`
const testSignature = "orb0adPhRCYND1WCAvPBr+qjm4STGtyvNDIDNBZ4Ir4="
const testKey = "other-secret"

func TestCalculateSignature(t *testing.T) {
	var cases = []struct {
		name string
		sKey string
		ts   string
		qp   string
		b    string
		es   string
		e    bool
	}{
		{
			name: "Succesful",
			sKey: testKey,
			ts:   testTs,
			qp:   testQp,
			b:    testBody,
			es:   testSignature,
			e:    true,
		},
		{
			name: "Wrong signature",
			sKey: testKey,
			ts:   testTs,
			qp:   testQp,
			b:    testBody,
			es:   "LISw4Je7n0/MkYDgVSzTJm8dW6BkytKTXMZZk1IElMs=",
			e:    false,
		},
		{
			name: "Empty query params and body",
			sKey: "secret",
			ts:   testTs,
			qp:   "",
			b:    "",
			es:   "LISw4Je7n0/MkYDgVSzTJm8dW6BkytKTXMZZk1IElMs=",
			e:    true,
		},
		{
			name: "Empty query params",
			sKey: "secret",
			ts:   testTs,
			qp:   "",
			b:    testBody,
			es:   "p2e20OtAg39DEmz1ORHpjQ556U4o1ZaH4NWbM9Q8Qjk=",
			e:    true,
		},
		{
			name: "Empty body",
			sKey: "secret",
			ts:   testTs,
			qp:   testQp,
			b:    "",
			es:   "Tfn+nRUBsn6lQgf6IpxBMS1j9lm7XsGjt5xh47M3jCk=",
			e:    true,
		},
	}
	for _, tt := range cases {
		v := NewValidator(tt.sKey)
		s, err := v.calculateSignature(tt.ts, tt.qp, []byte(tt.b))
		if err != nil {
			t.Errorf("Error calculating signature: %s, expected: %s", s, tt.es)
		}
		drs, _ := base64.StdEncoding.DecodeString(tt.es)
		e := bool(bytes.Compare(s, drs) == 0)
		if e != tt.e {
			t.Errorf("Unexpected signature: %s, test case: %s", s, tt.name)
		}
	}
}
func TestValidTimestamp(t *testing.T) {
	now := time.Now()
	nowts := fmt.Sprintf("%d", now.Unix())
	var cases = []struct {
		name string
		ts   string
		e    bool
	}{
		{
			name: "Succesful",
			ts:   nowts,
			e:    true,
		},
		{
			name: "Empty time stamp",
			ts:   "",
			e:    false,
		},
		{
			name: "Invalid time stamp",
			ts:   "wrongTs",
			e:    false,
		},
		{
			name: "Time stamp 24 hours in the futute",
			ts:   fmt.Sprintf("%d", now.AddDate(0, 0, 1).Unix()),
			e:    false,
		},
		{
			name: "Time stamp 24 hours in the past",
			ts:   fmt.Sprintf("%d", now.AddDate(0, 0, -1).Unix()),
			e:    false,
		},
	}

	for _, tt := range cases {
		v := NewValidator(testKey)
		r := v.validTimestamp(tt.ts)
		if r != tt.e {
			t.Errorf("Unexpected error validating ts: %s, test case: %s", tt.ts, tt.name)
		}
	}
}

func TestValidSignature(t *testing.T) {
	var cases = []struct {
		name string
		ts   string
		qp   string
		b    string
		s    string
		e    bool
	}{
		{
			name: "succesful",
			ts:   testTs,
			qp:   testQp,
			b:    testBody,
			s:    testSignature,
			e:    true,
		},
		{
			name: "Unorganized query params",
			ts:   testTs,
			qp:   "def=bar&abc=foo",
			b:    testBody,
			s:    testSignature,
			e:    true,
		},
		{
			name: "Wrong signature received",
			ts:   testTs,
			qp:   testQp,
			b:    testBody,
			s:    "wrong signature",
			e:    false,
		},
	}

	for _, tt := range cases {
		v := NewValidator(testKey)
		ValidityWindow = time.Hour * 100000
		r := v.validSignature(tt.ts, tt.qp, []byte(tt.b), tt.s)
		if r != tt.e {
			t.Errorf("Unexpected error validating signature: %s, test case: %s", tt.s, tt.name)
		}
	}
}
func TestValidate(t *testing.T) {
	var cases = []struct {
		name string
		k    string
		ts   string
		s    string
		sh   string
		tsh  string
		e    int
	}{
		{
			name: "Succesful",
			k:    testKey,
			ts:   testTs,
			s:    testSignature,
			sh:   sHeader,
			tsh:  tsHeader,
			e:    http.StatusOK,
		},
		{
			name: "NO Access key",
			k:    "",
			ts:   testTs,
			s:    testSignature,
			sh:   sHeader,
			tsh:  tsHeader,
			e:    http.StatusUnauthorized,
		},
		{
			name: "Request with empty time stamp",
			k:    testKey,
			ts:   "",
			s:    testSignature,
			sh:   sHeader,
			tsh:  tsHeader,
			e:    http.StatusUnauthorized,
		},
		{
			name: "Request with empty signature",
			k:    testKey,
			ts:   testTs,
			s:    "",
			sh:   sHeader,
			tsh:  tsHeader,
			e:    http.StatusUnauthorized,
		},
		{
			name: "Request with wrong signature header",
			k:    testKey,
			ts:   testTs,
			s:    testSignature,
			sh:   "wrong-header",
			tsh:  tsHeader,
			e:    http.StatusUnauthorized,
		},
		{
			name: "Request with wrong timestamp header",
			k:    testKey,
			ts:   testTs,
			s:    testSignature,
			sh:   sHeader,
			tsh:  "wrong-header",
			e:    http.StatusUnauthorized,
		},
	}

	for _, tt := range cases {
		v := NewValidator(tt.k)
		testTime, _ := stringToTime(testTs)
		ValidityWindow = time.Now().Add(time.Second*1).Sub(testTime) * 2
		ts := httptest.NewServer(v.Validate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})))
		defer ts.Close()

		client := &http.Client{}
		req, _ := http.NewRequest("GET", ts.URL+"?"+testQp, strings.NewReader(testBody))
		req.Header.Set(tt.sh, tt.s)
		req.Header.Set(tt.tsh, tt.ts)
		res, _ := client.Do(req)
		if res.StatusCode != tt.e {
			t.Errorf("Unexpected response code: %s, test case: %s", res.Status, tt.name)
		}
	}

}
