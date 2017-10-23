package messagebird

import "time"

type HLR struct {
	Id              string
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
