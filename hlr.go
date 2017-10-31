package messagebird

import "time"

// HLR stands for Home Location Register.
// Contains information about the subscribers identity, telephone number, the associated services and general information about the location of the subscriber
type HLR struct {
	ID              string
	HRef            string
	MSISDN          int
	Network         int
	Reference       string
	Status          string
	Details         map[string]interface{}
	CreatedDatetime *time.Time
	StatusDatetime  *time.Time
	Errors          []Error
}

// HLRList represents a list of HLR requests.
type HLRList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      map[string]*string
	Items      []HLR
}
