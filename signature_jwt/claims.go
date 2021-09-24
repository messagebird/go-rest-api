package signature_jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// maxSkew is the maximum time skew that we accept.  Sometimes the Internet is *so* fast
// that messages are received before they are sent, or the clocks of two servers are not
// in-sync, whichever cause seems more likely to you.
const maxSkew = 1 * time.Second

// Claims replaces jwt.StandardClaims as it checks all aspects of the the JWT token that
// have been specified by the MessageBird RFC.
type Claims struct {
	// The following 3 fields are added to Claims before JWT is parsed so that the
	// immediately following call to Valid() by jwt-go has *all* necessary information to
	// determine whether JWT is valid.  These fields should not be overwritten by JSON
	// unmarshal.
	receivedTime       time.Time
	correctPayloadHash string
	correctURLHash     string
	skipURLValidation  bool

	Issuer         string `json:"iss"`
	NotBefore      int64  `json:"nbf"`
	ExpirationTime int64  `json:"exp"`
	JWTID          string `json:"jti"`
	URLHash        string `json:"url_hash"`
	PayloadHash    string `json:"payload_hash,omitempty"`
}

// Valid is called by jwt-go after the Claims struct has been filled.  If an error is
// returned, it means that the JWT should not be trusted.
func (c Claims) Valid() error {
	var errs []string

	if c.Issuer != "MessageBird" {
		errs = append(errs, "claim iss has wrong value")
	}

	if iat := time.Unix(c.NotBefore, int64(c.receivedTime.Nanosecond())).Add(-maxSkew); c.receivedTime.Before(iat) {
		errs = append(errs, "claim nbf is in the future")
	}

	if exp := time.Unix(c.ExpirationTime, int64(c.receivedTime.Nanosecond())).Add(maxSkew); c.receivedTime.After(exp) {
		errs = append(errs, "claim exp is in the past")
	}

	if c.JWTID == "" {
		errs = append(errs, "claim jti is empty or missing")
	}

	if !c.skipURLValidation && c.correctURLHash != c.URLHash {
		errs = append(errs, "claim url_hash is invalid")
	}

	switch {
	case c.correctPayloadHash == "" && c.PayloadHash != "":
		errs = append(errs, "claim payload_hash is set but actual payload is missing")
	case c.correctPayloadHash != "" && c.PayloadHash == "":
		errs = append(errs, "claim payload_hash is not set but payload is present")
	case c.correctPayloadHash != c.PayloadHash:
		errs = append(errs, "claim payload_hash is invalid")
	}

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errs, "; "))
}

// Claims satisfies jwt.Claims.
var _ jwt.Claims = Claims{}
