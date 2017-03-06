package messagebird

import "time"

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
