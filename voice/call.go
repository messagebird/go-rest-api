package voice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v9"
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
		return fmt.Errorf("unable to parse Call CreatedAt: %v", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, raw.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse Call UpdatedAt: %v", err)
	}
	var endedAt *time.Time
	if raw.EndedAt != "" {
		eat, err := time.Parse(time.RFC3339, raw.EndedAt)
		if err != nil {
			return fmt.Errorf("unable to parse Call EndedAt: %v", err)
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

type createCallRequest struct {
	Source      string       `json:"source"`
	Destination string       `json:"destination"`
	CallFlow    CallFlow     `json:"callflow"`
	Webhook     *callWebhook `json:"webhook,omitempty"`
}

type callWebhook struct {
	URL   string `json:"url,omitempty"`
	Token string `json:"token,omitempty"`
}

type response struct {
	Data []Call `json:"data"`
}

// CallByID fetches a call by it's ID.
//
// An error is returned if no such call flow exists or is accessible.
func CallByID(client messagebird.MessageBirdClient, id string) (*Call, error) {
	var resp response

	if err := client.Request(&resp, http.MethodGet, apiRoot+"/calls/"+id, nil); err != nil {
		return nil, err
	}

	return &resp.Data[0], nil
}

// Calls returns a Paginator which iterates over all Calls.
func Calls(client messagebird.MessageBirdClient) *Paginator {
	return newPaginator(client, apiRoot+"/calls/", reflect.TypeOf(Call{}))
}

// InitiateCall initiates an outbound call.
//
// When placing a call, you pass the source (the caller ID), the destination
// (the number/address that will be called), and the callFlow (the call flow to
// execute when the call is answered).
func InitiateCall(client messagebird.MessageBirdClient, source, destination string, callflow CallFlow, webhook *Webhook) (*Call, error) {
	req := createCallRequest{
		Source:      source,
		Destination: destination,
		CallFlow:    callflow,
	}

	if webhook != nil {
		req.Webhook = &callWebhook{webhook.URL, webhook.Token}
	}

	var resp response

	if err := client.Request(&resp, http.MethodPost, fmt.Sprintf("%s/%s", apiRoot, callsPath), req); err != nil {
		return nil, err
	}
	return &resp.Data[0], nil
}

// Delete deletes the Call.
//
// If the call is in progress, it hangs up all legs.
func (call *Call) Delete(client messagebird.MessageBirdClient) error {
	return client.Request(nil, http.MethodDelete, fmt.Sprintf("%s/%s/%s", apiRoot, callsPath, call.ID), nil)
}

// Legs returns a paginator over all Legs associated with a call.
func (call *Call) Legs(client messagebird.MessageBirdClient) *Paginator {
	return newPaginator(client, fmt.Sprintf("%s/%s/%s/%s", apiRoot, callsPath, call.ID, legsPath), reflect.TypeOf(Leg{}))
}
