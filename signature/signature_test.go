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
		sKey string
		ts   string
		qp   string
		b    string
		es   string
		e    bool
	}{
		{
			sKey: testKey,
			ts:   testTs,
			qp:   testQp,
			b:    testBody,
			es:   testSignature,
			e:    true,
		},
		{
			sKey: testKey,
			ts:   testTs,
			qp:   testQp,
			b:    testBody,
			es:   "LISw4Je7n0/MkYDgVSzTJm8dW6BkytKTXMZZk1IElMs=",
			e:    false,
		},
		{
			sKey: "secret",
			ts:   testTs,
			qp:   "",
			b:    "",
			es:   "LISw4Je7n0/MkYDgVSzTJm8dW6BkytKTXMZZk1IElMs=",
			e:    true,
		},
		{
			sKey: "secret",
			ts:   testTs,
			qp:   "",
			b:    testBody,
			es:   "p2e20OtAg39DEmz1ORHpjQ556U4o1ZaH4NWbM9Q8Qjk=",
			e:    true,
		},
		{
			sKey: "secret",
			ts:   testTs,
			qp:   testQp,
			b:    "",
			es:   "Tfn+nRUBsn6lQgf6IpxBMS1j9lm7XsGjt5xh47M3jCk=",
			e:    true,
		},
	}
	for i, tt := range cases {
		v := NewValidator(tt.sKey, nil, nil, nil)
		s, err := v.CalculateSignature(tt.ts, tt.qp, []byte(tt.b))
		if err != nil {
			t.Errorf("Error calculating signature: %s, expected: %s", s, tt.es)
		}
		drs, _ := base64.StdEncoding.DecodeString(tt.es)
		e := bool(bytes.Compare(s, drs) == 0)
		if e != tt.e {
			t.Errorf("Unexpected signature: %s, test case: %d", s, i)
		}
	}
}
func TestValidTimestamp(t *testing.T) {
	var p float64 = 2
	now := time.Now()
	nowts := fmt.Sprintf("%d", now.Unix())
	var cases = []struct {
		ts string
		p  ValidityPeriod
		e  bool
	}{
		{
			ts: nowts,
			p:  nil,
			e:  true,
		},
		{
			ts: "",
			p:  nil,
			e:  false,
		},
		{
			ts: "wrongTs",
			p:  nil,
			e:  false,
		},
		{
			ts: nowts,
			p:  &p,
			e:  true,
		},
		{
			ts: fmt.Sprintf("%d", now.AddDate(0, 0, 1).Unix()),
			p:  &p,
			e:  false,
		},
		{
			ts: fmt.Sprintf("%d", now.AddDate(0, 0, -1).Unix()),
			p:  &p,
			e:  false,
		},
	}

	for i, tt := range cases {
		v := NewValidator(testKey, tt.p, nil, nil)
		r := v.ValidTimestamp(tt.ts)
		if r != tt.e {
			t.Errorf("Unexpected error validating ts: %s, test case: %d", tt.ts, i)
		}
	}
}

func TestValidSignature(t *testing.T) {
	var cases = []struct {
		ts string
		qp string
		b  string
		s  string
		e  bool
	}{
		{
			ts: testTs,
			qp: testQp,
			b:  testBody,
			s:  testSignature,
			e:  true,
		},
		{
			ts: testTs,
			qp: "def=bar&abc=foo",
			b:  testBody,
			s:  testSignature,
			e:  true,
		},
		{
			ts: testTs,
			qp: testQp,
			b:  testBody,
			s:  "wrong signature",
			e:  false,
		},
	}

	for i, tt := range cases {
		v := NewValidator(testKey, nil, nil, nil)
		r := v.ValidSignature(tt.ts, tt.qp, []byte(tt.b), tt.s)
		if r != tt.e {
			t.Errorf("Unexpected error validating signature: %s, test case: %d", tt.s, i)
		}
	}
}
func TestValidate(t *testing.T) {
	var cases = []struct {
		k   string
		ts  string
		s   string
		sh  string
		tsh string
		e   int
	}{
		{
			k:   testKey,
			ts:  testTs,
			s:   testSignature,
			sh:  sHeader,
			tsh: tsHeader,
			e:   http.StatusOK,
		},
		{
			k:   "",
			ts:  testTs,
			s:   testSignature,
			sh:  sHeader,
			tsh: tsHeader,
			e:   http.StatusUnauthorized,
		},
		{
			k:   testKey,
			ts:  "",
			s:   testSignature,
			sh:  sHeader,
			tsh: tsHeader,
			e:   http.StatusUnauthorized,
		},
		{
			k:   testKey,
			ts:  testTs,
			s:   "",
			sh:  sHeader,
			tsh: tsHeader,
			e:   http.StatusUnauthorized,
		},
		{
			k:   testKey,
			ts:  testTs,
			s:   testSignature,
			sh:  "wrong-header",
			tsh: tsHeader,
			e:   http.StatusUnauthorized,
		},
		{
			k:   testKey,
			ts:  testTs,
			s:   testSignature,
			sh:  sHeader,
			tsh: "wrong-header",
			e:   http.StatusUnauthorized,
		},
	}

	for i, tt := range cases {
		v := NewValidator(tt.k, nil, nil, nil)
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
			t.Errorf("Unexpected response code: %s, test case: %d", res.Status, i)
		}
	}

}
