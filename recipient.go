package messagebird

import "time"

// Recipient struct holds information for a single msisdn with status details.
type Recipient struct {
	Recipient      		int64
	Status         		string
	StatusDatetime 		*time.Time
	MessagePartCount 	int

}

// Recipients holds a collection of Recepient structs along with send stats.
type Recipients struct {
	TotalCount               int
	TotalSentCount           int
	TotalDeliveredCount      int
	TotalDeliveryFailedCount int
	Items                    []Recipient
}
