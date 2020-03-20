package messagebird

import (
	"fmt"
	"strings"
)

// Error holds details including error code, human readable description and optional parameter that is related to the error.
type Error struct {
	Code        int
	Description string
	Parameter   string
}

// Error implements error interface.
func (e Error) Error() string {
	return e.Description
}

// ErrorResponse represents errored API response.
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error implements error interface.
func (r ErrorResponse) Error() string {
	var inners []string
	for _, inner := range r.Errors {
		inners = append(inners, inner.Error())
	}
	return fmt.Sprintf("API errors: %s", strings.Join(inners, ", "))
}
