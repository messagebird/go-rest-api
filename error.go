package messagebird

import "fmt"

// Error holds details including error code, human readable description and optional parameter that is related to the error.
type Error struct {
	Code        int
	Description string
	Parameter   string
}

// ErrorResponse represents errored API response.
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error implements error interface.
func (r ErrorResponse) Error() string {
	eString := "API returned an error: "
	for i, e := range r.Errors {
		eString = eString + fmt.Sprintf("code: %d, description: %s, parameter: %s", e.Code, e.Description, e.Parameter)
		if i < len(r.Errors)-1 {
			eString = eString + ", "
		}
	}

	return eString
}
