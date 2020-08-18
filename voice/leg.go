package voice

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v6"
)

// LegStatus enumerates all valid values for a leg status.
type LegStatus string

const (
	// LegStatusStarting indicates that a leg is currently starting.
	LegStatusStarting LegStatus = "starting"
	// LegStatusRinging indicates that a leg is connected and ringing.
	LegStatusRinging LegStatus = "ringing"
	// LegStatusOngoing indicates a healthy leg that is currently participating
	// in a call.
	LegStatusOngoing LegStatus = "ongoing"
	// LegStatusBusy indicates that a leg could not be established because the
	// other side is busy.
	LegStatusBusy LegStatus = "busy"
	// LegStatusNoAnswer indicates that a leg could not be established because
	// the other side did not pick up.
	LegStatusNoAnswer LegStatus = "no_answer"
	// LegStatusFailed indicates some kind of failure of a leg.
	LegStatusFailed LegStatus = "failed"
	// LegStatusHangup indicates that a leg has been hung up.
	LegStatusHangup LegStatus = "hangup"
)

// LegDirection indicates the direction of some leg in a call.
type LegDirection string

const (
	// LegDirectionOutgoing is the direction of a leg that are created when a
	// call is transferred.
	LegDirectionOutgoing LegDirection = "outgoing"
	// LegDirectionIncoming is the direction of a leg that is created when a
	// number is called.
	LegDirectionIncoming LegDirection = "incoming"
)

// A Leg describes a leg object (inbound or outbound) that belongs to a call.
//
// At least one leg exists per call. Inbound legs are being created when an
// incoming call to a Number is being initiated. Outgoing legs are created when
// a call is transferred or when a call is being originated from the API.
type Leg struct {
	// The unique ID of the leg.
	ID string
	// The unique ID of the call that this leg belongs to.
	CallID string
	// The number/SIP URL that is making the connection.
	Source string
	// The number/SIP URL that a connection is made to.
	Destination string
	// The status of the leg. Possible values: starting, ringing, ongoing,
	// busy, no_answer, failed and hangup.
	Status LegStatus
	// The direction of the leg, indicating if it's an incoming connection or
	// outgoing (e.g. for transferring a call). Possible values: incoming,
	// outgoing.
	Direction LegDirection
	// The cost of the leg. The amount relates to the currency parameter.
	Cost float64
	// The three-letter currency code (ISO 4217) related to the cost of the
	// leg.
	Currency string
	// The duration of the leg.
	//
	// Truncated to seconds.
	Duration time.Duration
	// The date-time the leg was created.
	CreatedAt time.Time
	// The date-time the leg was last updated.
	UpdatedAt time.Time
	// The date-time the leg was answered.
	AnsweredAt *time.Time
	// The date-time the leg ended.
	EndedAt *time.Time
}

type jsonLeg struct {
	ID          string  `json:"id"`
	CallID      string  `json:"callID"`
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	Status      string  `json:"status"`
	Direction   string  `json:"direction"`
	Cost        float64 `json:"cost"`
	Currency    string  `json:"currency"`
	Duration    int     `json:"duration"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	AnsweredAt  string  `json:"answeredAt,omitempty"`
	EndedAt     string  `json:"endedAt,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (leg *Leg) UnmarshalJSON(data []byte) error {
	var raw jsonLeg
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	createdAt, err := time.Parse(time.RFC3339, raw.CreatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse Leg CreatedAt: %v", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, raw.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse Leg UpdatedAt: %v", err)
	}
	var answeredAt *time.Time
	if raw.EndedAt != "" {
		aat, err := time.Parse(time.RFC3339, raw.EndedAt)
		if err != nil {
			return fmt.Errorf("unable to parse Leg AnsweredAt: %v", err)
		}
		answeredAt = &aat
	}
	var endedAt *time.Time
	if raw.EndedAt != "" {
		eat, err := time.Parse(time.RFC3339, raw.EndedAt)
		if err != nil {
			return fmt.Errorf("unable to parse Leg EndedAt: %v", err)
		}
		endedAt = &eat
	}
	*leg = Leg{
		ID:          raw.ID,
		CallID:      raw.CallID,
		Source:      raw.Source,
		Destination: raw.Destination,
		Status:      LegStatus(raw.Status),
		Direction:   LegDirection(raw.Direction),
		Cost:        raw.Cost,
		Currency:    raw.Currency,
		Duration:    time.Second * time.Duration(raw.Duration),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		AnsweredAt:  answeredAt,
		EndedAt:     endedAt,
	}
	return nil
}

// Recordings retrieves the Recording objects associated with a leg.
func (leg *Leg) Recordings(client *messagebird.Client) *Paginator {
	return newPaginator(client, fmt.Sprintf("%s/calls/%s/legs/%s/recordings", apiRoot, leg.CallID, leg.ID), reflect.TypeOf(Recording{}))
}
