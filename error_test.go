package messagebird

import "testing"

func TestErrorResponseError(t *testing.T) {
	tests := []struct {
		name   string
		errors []Error
		expect string
	}{
		{
			name: "single error",
			errors: []Error{
				Error{
					Code:        2,
					Description: "Request not allowed (incorrect access_key)",
					Parameter:   "access_key",
				},
			},
			expect: "API returned an error: code: 2, description: Request not allowed (incorrect access_key), parameter: access_key",
		},
		{
			name: "multiple errors",
			errors: []Error{
				Error{
					Code:        2,
					Description: "Request not allowed (incorrect access_key)",
					Parameter:   "access_key",
				},
				Error{
					Code:        2,
					Description: "Request not allowed (incorrect access_key)",
					Parameter:   "access_key",
				},
			},
			expect: "API returned an error: code: 2, description: Request not allowed (incorrect access_key), parameter: access_key, code: 2, description: Request not allowed (incorrect access_key), parameter: access_key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := ErrorResponse{tt.errors}
			if er.Error() != tt.expect {
				t.Errorf("expected error message to be:\n%s\ngot:\n%s", tt.expect, er.Error())
			}
		})
	}
}
