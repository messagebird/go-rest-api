package messagebird

import "time"

// Recipient struct holds information for a single msisdn with status details.
type Recipient struct {
	Recipient              int64
	Status                 string
	StatusDatetime         *time.Time
	RecipientCountry       *string
	RecipientCountryPrefix *int
	RecipientOperator      *string
	MessageLength          *int
	StatusErrorCode        *int
	StatusReason           *string
	Price                  *Price
	Mccmnc                 *string
	Mcc                    *string
	Mnc                    *string
	MessagePartCount       int
}

type Price struct {
	Amount   float64
	Currency string
}

// Recipients holds a collection of Recepient structs along with send stats.
type Recipients struct {
	TotalCount               int
	TotalSentCount           int
	TotalDeliveredCount      int
	TotalDeliveryFailedCount int
	Items                    []Recipient
}
