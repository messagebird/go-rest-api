package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StringToTime(s string) (time.Time, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

func HMACSHA256(message, signingKey []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(signingKey))
	if _, err := mac.Write(message); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

type Validator struct {
	SigningKey  string
	Period      float64
	Log         *log.Logger
	LogMesssage string
}

func (v *Validator) ValidTimestamp(ts time.Time) bool {
	if v.Period != 0 {
		now := time.Now()
		if diff := now.Sub(ts); diff.Hours() > v.Period {
			return false
		}
	}
	return true
}

func (v *Validator) CalculateSignature(ts, qp string, b []byte) ([]byte, error) {
	var m bytes.Buffer
	bh := sha256.Sum256(b)
	fmt.Fprintf(&m, "%s\n%s\n%s", ts, qp, bh[:])
	s, err := HMACSHA256(m.Bytes(), []byte(v.SigningKey))
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (v *Validator) ValidSignature(ts, qp, rs string, b []byte) bool {
	es, _ := v.CalculateSignature(ts, qp, b)
	drs, _ := base64.StdEncoding.DecodeString(rs)
	return hmac.Equal(drs, es)
}

func (v *Validator) Wrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts := r.Header.Get("MessageBird-Request-Timestamp")
		rs := r.Header.Get("MessageBird-Request-Signature")
		if ts == "" || rs == "" {
			http.Error(w, "Request not allowed", http.StatusUnauthorized)
			return
		}
		t, _ := StringToTime(ts)
		b, _ := ioutil.ReadAll(r.Body)
		if v.ValidTimestamp(t) == false || v.ValidSignature(ts, r.URL.RawQuery, rs, b) == false {
			http.Error(w, "Request not allowed", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
