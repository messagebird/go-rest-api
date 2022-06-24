package conversation

import (
	"encoding/json"
	"time"
)

type Contact struct {
	ID            string
	Href          string
	MSISDN        string
	FirstName     string
	LastName      string
	CustomDetails map[string]interface{}
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// UnmarshalJSON is used to unmarshal the MSISDN to a string rather than an
// int64. The API returns integers, but this client always uses strings.
// Exposing a json.Number doesn't seem nice.
func (c *Contact) UnmarshalJSON(data []byte) error {
	target := struct {
		ID            string
		Href          string
		MSISDN        json.Number
		FirstName     string
		LastName      string
		CustomDetails map[string]interface{}
		CreatedAt     *time.Time
		UpdatedAt     *time.Time
	}{}

	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}

	// In many cases, the CustomDetails will contain the user ID. As
	// CustomDetails has interface{} values, these are unmarshalled as floats.
	// Convert them to int64.
	// Map key is not a typo: API returns userId and not userID.
	if val, ok := target.CustomDetails["userId"]; ok {
		var userID float64
		if userID, ok = val.(float64); ok {
			target.CustomDetails["userId"] = int64(userID)
		}
	}

	*c = Contact{
		target.ID,
		target.Href,
		target.MSISDN.String(),
		target.FirstName,
		target.LastName,
		target.CustomDetails,
		target.CreatedAt,
		target.UpdatedAt,
	}

	return nil
}
