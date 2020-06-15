/*
Package signature implements signature verification for MessageBird webhooks.

To use define a new validator using your MessageBird Signing key.  You can use the
ValidRequest method, just pass the request as a parameter:

    validator := signature.NewValidator("your signing key")
    if err := validator.ValidRequest(r); err != nil {
        // handle error
    }

Or use the handler as a middleware for your server:

	http.Handle("/path", validator.Validate(YourHandler))

It will reject the requests that contain invalid signatures.
The validator uses a 5ms seconds window to accept requests as valid, to change
this value, set the ValidityWindow to the disired duration.
Take into account that the validity window works around the current time:
	[now - ValidityWindow/2, now + ValidityWindow/2]
*/
package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	tsHeader = "MessageBird-Request-Timestamp"
	sHeader  = "MessageBird-Signature"
)

// ValidityWindow defines the time window in which to validate a request.
var ValidityWindow = 5 * time.Second

// StringToTime converts from Unicode Epoch encoded timestamps to the time.Time type.
func stringToTime(s string) (time.Time, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

// Validator type represents a MessageBird signature validator.
type Validator struct {
	SigningKey string // Signing Key provided by MessageBird.
}

// NewValidator returns a signature validator object.
func NewValidator(signingKey string) *Validator {
	return &Validator{
		SigningKey: signingKey,
	}
}

// validTimestamp validates if the MessageBird-Request-Timestamp is a valid
// date and if the request is older than the validator Period.
func (v *Validator) validTimestamp(ts string) bool {
	t, err := stringToTime(ts)
	if err != nil {
		return false
	}
	diff := time.Now().Add(ValidityWindow / 2).Sub(t)
	return diff < ValidityWindow && diff > 0
}

// calculateSignature calculates the MessageBird-Signature using HMAC_SHA_256
// encoding and the timestamp, query params and body from the request:
// signature = HMAC_SHA_256(
//	TIMESTAMP + \n + QUERY_PARAMS + \n + SHA_256_SUM(BODY),
//	signing_key)
func (v *Validator) calculateSignature(ts, qp string, b []byte) ([]byte, error) {
	var m bytes.Buffer
	bh := sha256.Sum256(b)
	fmt.Fprintf(&m, "%s\n%s\n%s", ts, qp, bh[:])
	mac := hmac.New(sha256.New, []byte(v.SigningKey))
	if _, err := mac.Write(m.Bytes()); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

// validSignature takes the timestamp, query params and body from the request,
// calculates the expected signature and compares it to the one sent by MessageBird.
func (v *Validator) validSignature(ts, rqp string, b []byte, rs string) bool {
	uqp, err := url.Parse("?" + rqp)
	if err != nil {
		return false
	}
	es, err := v.calculateSignature(ts, uqp.Query().Encode(), b)
	if err != nil {
		return false
	}
	drs, err := base64.StdEncoding.DecodeString(rs)
	if err != nil {
		return false
	}
	return hmac.Equal(drs, es)
}

// ValidRequest is a method that takes care of the signature validation of
// incoming requests.
func (v *Validator) ValidRequest(r *http.Request) error {
	ts := r.Header.Get(tsHeader)
	rs := r.Header.Get(sHeader)
	if ts == "" || rs == "" {
		return fmt.Errorf("Unknown host: %s", r.Host)
	}
	b, _ := ioutil.ReadAll(r.Body)
	if !v.validTimestamp(ts) || !v.validSignature(ts, r.URL.RawQuery, b, rs) {
		return fmt.Errorf("Unknown host: %s", r.Host)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return nil
}

// Validate is a handler wrapper that takes care of the signature validation of
// incoming requests and rejects them if invalid or pass them on to your handler
// otherwise.
func (v *Validator) Validate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := v.ValidRequest(r); err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
