/*
Package signature_jwt implements signature verification for MessageBird webhooks.

To use define a new validator using your MessageBird Signing key. Can be
retrieved through https://dashboard.messagebird.com/developers/settings.
This is NOT your API key.

You can use the ValidateRequest method, just pass the request and base url as parameters:

    validator := signature_jwt.NewValidator([]byte("your signing key"))
	baseUrl := "https://yourdomain.com"
    if err := validator.ValidateRequest(r, baseUrl); err != nil {
        // handle error
    }

Or use the handler as a middleware for your server:

	http.Handle("/path", validator.Validate(YourHandler, baseUrl))

It will reject the requests that contain invalid signatures.
*/
package signature_jwt

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
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
	parser jwt.Parser
	keyFn  jwt.Keyfunc
}

// NewValidator returns a signature validator object.
func NewValidator(signingKey string) *Validator {
	return &Validator{
		parser: jwt.Parser{
			ValidMethods: allowedMethods,
		},
		keyFn: func(*jwt.Token) (interface{}, error) { return []byte(signingKey), nil },
	}
}

// ValidateSignature is a method that takes care of the signature validation of
// incoming requests.
func (v *Validator) ValidateSignature(signature, url string, payload []byte) (*jwt.Token, error) {
	claims := Claims{
		receivedTime: TimeFunc(),
	}

	if url != "" {
		claims.correctURLHash = sha256Hash([]byte(url))
	}
	if payload != nil && len(payload) != 0 {
		claims.correctPayloadHash = sha256Hash(payload)
	}

	if token, err := v.parser.ParseWithClaims(signature, &claims, v.keyFn); err != nil {
		return nil, fmt.Errorf("invalid jwt: %w", err)
	} else {
		return token, nil
	}
}

// ValidateRequest is a method that takes care of the signature validation of
// incoming requests.
func (v *Validator) ValidateRequest(r *http.Request, baseURL string) error {
	signature := r.Header.Get(signatureHeader)
	if signature == "" {
		return fmt.Errorf("signature not found")
	}

	var fullURL string
	if baseURL != "" {
		base, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("error parsing base url: %v", err)
		}
		fullURL = base.ResolveReference(r.URL).String()
	}

	b, _ := ioutil.ReadAll(r.Body)
	if _, err := v.ValidateSignature(signature, fullURL, b); err != nil {
		return fmt.Errorf("invalid signature: %s", err.Error())
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return nil
}

// Validate is a handler wrapper that takes care of the signature validation of
// incoming requests and rejects them if invalid or pass them on to your handler
// otherwise.
func (v *Validator) Validate(h http.Handler, baseURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := v.ValidateRequest(r, baseURL); err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func sha256Hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
