/*
Package signature implements signature verification for MessageBird webhooks.

To use define a new validator using your MessageBird Signing key. You can use the
ValidRequest method, just pass the request and base url as parameters:

    validator := signature.NewValidator([]byte("your signing key"))
	baseUrl := "https://messagebird.io"
    if err := validator.ValidRequest(r, baseUrl); err != nil {
        // handle error
    }

Or use the handler as a middleware for your server:

	http.Handle("/path", validator.Validate(YourHandler, baseUrl))

It will reject the requests that contain invalid signatures.
*/
package signature

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const signatureHeader = "MessageBird-Signature-JWT"

// TimeFunc provides the current time same as time.Now but can be overridden for testing.
var TimeFunc = time.Now

// allowedMethods lists the signing methods that we accept.  We only allow symmetric-key
// algorithms as our customer signing keys are currently all simple byte strings.  HMAC is
// also the only symkey signature method that is required by the RFC7518 Section 3.1 and
// thus should be supported by all JWT implementations.
var allowedMethods = []string{
	jwt.SigningMethodHS256.Name,
	jwt.SigningMethodHS384.Name,
	jwt.SigningMethodHS512.Name,
}

// Validator type represents a MessageBird signature validator.
type Validator struct {
	SigningKey []byte // Signing Key provided by MessageBird.
}

// NewValidator returns a signature validator object.
func NewValidator(signingKey []byte) *Validator {
	return &Validator{
		SigningKey: signingKey,
	}
}

// ValidSignature is a method that takes care of the signature validation of
// incoming requests.
func (v *Validator) ValidSignature(signature, url string, payload []byte) error {
	parser := jwt.Parser{ValidMethods: allowedMethods}
	keyFn := func(*jwt.Token) (interface{}, error) { return v.SigningKey, nil }

	claims := Claims{
		receivedTime:   TimeFunc(),
		correctURLHash: sha256Hash([]byte(url)),
	}
	if payload != nil && len(payload) != 0 {
		claims.correctPayloadHash = sha256Hash(payload)
	}

	if _, err := parser.ParseWithClaims(signature, &claims, keyFn); err != nil {
		return fmt.Errorf("invalid jwt: %w", err)
	}

	return nil
}

// ValidRequest is a method that takes care of the signature validation of
// incoming requests.
func (v *Validator) ValidRequest(r *http.Request, baseUrl string) error {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return fmt.Errorf("error parsing base url: %v", err)
	}
	signature := r.Header.Get(signatureHeader)
	if signature == "" {
		return fmt.Errorf("signature not found")
	}
	b, _ := ioutil.ReadAll(r.Body)
	if err := v.ValidSignature(signature, base.ResolveReference(r.URL).String(), b); err != nil {
		return fmt.Errorf("invalid signature: %s", err.Error())
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return nil
}

// Validate is a handler wrapper that takes care of the signature validation of
// incoming requests and rejects them if invalid or pass them on to your handler
// otherwise.
func (v *Validator) Validate(h http.Handler, baseUrl string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := v.ValidRequest(r, baseUrl); err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func sha256Hash(data []byte) string {
	if data == nil {
		return ""
	}

	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
