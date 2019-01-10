package signature

import (
	"bytes"
	"encoding/base64"
	"testing"
)

const testTs = "1544544948"
const testQp = "abc=foo&def=bar"
const testBody = `{"a key":"some value"}`
const testSignature = "orb0adPhRCYND1WCAvPBr+qjm4STGtyvNDIDNBZ4Ir4="

func TestCalculateSignature(t *testing.T) {
	v := NewValidator("other-secret", 2, nil, nil)
	s, err := v.CalculateSignature(testTs, testQp, []byte(testBody))
	if err != nil {
		t.Errorf("Error calculating signature: %s, expected: orb0adPhRCYND1WCAvPBr+qjm4STGtyvNDIDNBZ4Ir4=", s)
	}
	drs, _ := base64.StdEncoding.DecodeString(testSignature)
	if bytes.Compare(s, drs) != 0 {
		t.Errorf("Unexpected signature: %s, expected: orb0adPhRCYND1WCAvPBr+qjm4STGtyvNDIDNBZ4Ir4=", s)
	}
}
