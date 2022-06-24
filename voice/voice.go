package voice

import (
	"encoding/json"
	"fmt"
	"strings"

	messagebird "github.com/messagebird/go-rest-api/v8"
)

const apiRoot = "https://voice.messagebird.com/v1"

type ErrorResponse struct {
	Errors []Error
}

type Error struct {
	Code    int
	Message string
}

func init() {
	// The Voice API returns errors in a format that slightly differs from other
	// APIs. Here we instruct package messagebird to use our custom
	// voice.errorReader func, which has access to voice.ErrorResponse, to
	// unmarshal those. Package messagebird must not import the voice package to
	// safeguard against import cycles, so it can not use voice.ErrorResponse
	// directly.
	messagebird.SetVoiceErrorReader(errorReader)
}

// errorReader takes a []byte representation of a Voice API JSON error and
// parses it to a voice.ErrorResponse.
func errorReader(b []byte) error {
	var er ErrorResponse
	if err := json.Unmarshal(b, &er); err != nil {
		return fmt.Errorf("encoding/json: Unmarshal: %v", err)
	}
	return er
}

func (e ErrorResponse) Error() string {
	errStrings := make([]string, len(e.Errors))
	for i, v := range e.Errors {
		errStrings[i] = v.Error()
	}
	return strings.Join(errStrings, "; ")
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, message: %q", e.Code, e.Message)
}
