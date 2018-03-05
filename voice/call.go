package voice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

// CallStatus enumerates all valid values for a call status.
type CallStatus string

const (
	// CallStatusStarting is the status of a call that is currently being set up.
	CallStatusStarting CallStatus = "starting"
	// CallStatusOngoing indicates that a call is active.
	CallStatusOngoing CallStatus = "ongoing"
	// CallStatusEnded indicates that a call has been terminated.
	CallStatusEnded CallStatus = "ended"
)

// A Call describes a voice call which is  made to a number.
//
// A call has legs which are incoming or outgoing voice connections. An
// incoming leg is created when somebody calls a number. Outgoing legs are
// created when a call is transferred.
type Call struct {
	ID          string
	Status      CallStatus
	Source      string
	Destination string
	NumberID    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	EndedAt     *time.Time
}

type jsonCall struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	NumberID    string `json:"numberId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	EndedAt     string `json:"endedAt,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
func (call Call) MarshalJSON() ([]byte, error) {
	endedAt := ""
	if call.EndedAt != nil {
		endedAt = call.EndedAt.Format(time.RFC3339)
	}
	data := jsonCall{
		ID:          call.ID,
		Status:      string(call.Status),
		Source:      call.Source,
		Destination: call.Destination,
		NumberID:    call.NumberID,
		CreatedAt:   call.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   call.UpdatedAt.Format(time.RFC3339),
		EndedAt:     endedAt,
	}
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (call *Call) UnmarshalJSON(data []byte) error {
	var raw jsonCall
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	createdAt, err := time.Parse(time.RFC3339, raw.CreatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse CallFlow CreatedAt: %v", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, raw.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse CallFlow UpdatedAt: %v", err)
	}
	var endedAt *time.Time
	if raw.EndedAt != "" {
		eat, err := time.Parse(time.RFC3339, raw.EndedAt)
		if err != nil {
			return fmt.Errorf("unable to parse CallFlow EndedAt: %v", err)
		}
		endedAt = &eat
	}
	*call = Call{
		ID:          raw.ID,
		Status:      CallStatus(raw.Status),
		Source:      raw.Source,
		Destination: raw.Destination,
		NumberID:    raw.NumberID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		EndedAt:     endedAt,
	}
	return nil
}

// CallByID fetches a call by it's ID.
//
// An error is returned if no such call flow exists or is accessible.
func CallByID(client *messagebird.Client, id string) (*Call, error) {
	call := &Call{}
	err := client.Request(call, http.MethodGet, "calls/"+id, nil)
	return call, err
}

// Calls returns a Paginator which iterates over all Calls.
func Calls(client *messagebird.Client) *Paginator {
	return newPaginator(client, "calls/", reflect.TypeOf(Call{}))
}

// InitiateCall initiates an outbound call.
//
// When placing a call, you pass the source (the caller ID), the destination
// (the number/address that will be called), and the callFlow (the call flow to
// execute when the call is answered).
func InitiateCall(client *messagebird.Client, source, destination string, callflow CallFlow, webhook *Webhook) (*Call, error) {
	body := struct {
		Source      string   `json:"source"`
		Destination string   `json:"destination"`
		Callflow    CallFlow `json:"callflow"`
		Webhook     struct {
			URL   string `json:"url,omitempty"`
			Token string `json:"token,omitempty"`
		}
	}{
		Source:      source,
		Destination: destination,
		Callflow:    callflow,
	}
	if webhook != nil {
		body.Webhook.URL = webhook.URL
		body.Webhook.Token = webhook.Token
	}
	call := &Call{}
	err := client.Request(call, http.MethodPost, "calls/", body)
	return call, err
}

// Delete deletes the Call.
//
// If the call is in progress, it hangs up all legs.
func (call *Call) Delete(client *messagebird.Client) error {
	return client.Request(nil, http.MethodDelete, "calls/"+call.ID, nil)
}

// Legs returns a paginator over all Legs associated with a call.
func (call *Call) Legs(client *messagebird.Client) *Paginator {
	return newPaginator(client, fmt.Sprintf("calls/%s/legs", call.ID), reflect.TypeOf(Leg{}))
}
