package signature

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	receivedTime       time.Time `json:"-"`
	correctPayloadHash string    `json:"-"`
	correctURLHash     string    `json:"-"`

	Issuer         string `json:"iss"`
	IssuedAt       int64  `json:"iat"`
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
		errs = append(errs, "wrong iss")
	}

	if iat := time.Unix(c.IssuedAt, int64(c.receivedTime.Nanosecond())).Add(-maxSkew); c.receivedTime.Before(iat) {
		errs = append(errs, "claim iat is in the future")
	}

	if exp := time.Unix(c.ExpirationTime, int64(c.receivedTime.Nanosecond())).Add(maxSkew); c.receivedTime.After(exp) {
		errs = append(errs, "claim exp is in the past")
	}

	if c.JWTID == "" {
		errs = append(errs, "claim jti is empty or missing")
	}

	if c.correctURLHash != c.URLHash {
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
