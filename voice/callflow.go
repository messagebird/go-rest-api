package voice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/messagebird/go-rest-api"
)

// A CallFlow describes the flow of operations (steps) to be executed when
// handling an incoming call.
type CallFlow struct {
	ID    string
	Title string

	// Each object in the steps array in a call flow describes an operation
	// that executes during a call, e.g. transferring a call or playing back an
	// audio file.
	Steps []CallFlowStep

	// Record instructs the Voice system to record the entire call.
	//
	// Note that this is distinct from the Record CallFlow step which records
	// only a single message from the callee.
	Record bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type jsonCallFlow struct {
	ID        string         `json:"id,omitempty"`
	Title     string         `json:"title"`
	Steps     []CallFlowStep `json:"steps"`
	Record    bool           `json:"record"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
}

// MarshalJSON implements the json.Marshaler interface.
func (callflow CallFlow) MarshalJSON() ([]byte, error) {
	data := jsonCallFlow{
		ID:    callflow.ID,
		Title: callflow.Title,
		// Steps are able to serialize themselves to JSON.
		Steps:     callflow.Steps,
		Record:    callflow.Record,
		CreatedAt: callflow.CreatedAt.Format(time.RFC3339),
		UpdatedAt: callflow.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (callflow *CallFlow) UnmarshalJSON(data []byte) error {
	var stepTypeLookahead struct {
		Steps []struct {
			Action string `json:"action"`
		} `json:"steps"`
	}
	if err := json.Unmarshal(data, &stepTypeLookahead); err != nil {
		return err
	}
	raw := jsonCallFlow{
		Steps: make([]CallFlowStep, len(stepTypeLookahead.Steps)),
	}
	for i, s := range stepTypeLookahead.Steps {
		switch s.Action {
		case "transfer":
			raw.Steps[i] = &CallFlowTransferStep{}
		case "say":
			raw.Steps[i] = &CallFlowSayStep{}
		case "play":
			raw.Steps[i] = &CallFlowPlayStep{}
		case "pause":
			raw.Steps[i] = &CallFlowPauseStep{}
		case "record":
			raw.Steps[i] = &CallFlowRecordStep{}
		case "fetchCallFlow":
			raw.Steps[i] = &CallFlowFetchStep{}
		case "hangup":
			raw.Steps[i] = &CallFlowHangupStep{}
		default:
			return fmt.Errorf("unknown step action: %q", s.Action)
		}
	}

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
	*callflow = CallFlow{
		ID:        raw.ID,
		Title:     raw.Title,
		Steps:     raw.Steps,
		Record:    raw.Record,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return nil
}

// CallFlowByID fetches a callflow by it's ID.
//
// An error is returned if no such call flow exists or is accessible.
func CallFlowByID(client *messagebird.Client, id string) (*CallFlow, error) {
	var data struct {
		Data []CallFlow `json:"data"`
	}
	if err := client.Request(&data, http.MethodGet, apiRoot+"/call-flows/"+id, nil); err != nil {
		return nil, err
	}
	return &data.Data[0], nil
}

// CallFlows returns a Paginator which iterates over all CallFlows.
func CallFlows(client *messagebird.Client) *Paginator {
	return newPaginator(client, apiRoot+"/call-flows/", reflect.TypeOf(CallFlow{}))
}

// Create creates the callflow remotely.
//
// The callflow is updated in-place.
func (callflow *CallFlow) Create(client *messagebird.Client) error {
	var data struct {
		Data []CallFlow `json:"data"`
	}
	if err := client.Request(&data, http.MethodPost, apiRoot+"/call-flows/", callflow); err != nil {
		return err
	}
	*callflow = data.Data[0]
	return nil
}

// Update updates the call flow by overwriting it.
//
// An error is returned if no such call flow exists or is accessible.
func (callflow *CallFlow) Update(client *messagebird.Client) error {
	var data struct {
		Data []CallFlow `json:"data"`
	}
	if err := client.Request(&data, http.MethodPut, apiRoot+"/call-flows/"+callflow.ID, callflow); err != nil {
		return err
	}
	*callflow = data.Data[0]
	return nil
}

// Delete deletes the CallFlow.
func (callflow *CallFlow) Delete(client *messagebird.Client) error {
	return client.Request(nil, http.MethodDelete, apiRoot+"/call-flows/"+callflow.ID, nil)
}

// A CallFlowStep is a single step that can be taken in a callflow.
//
// It can be any of CallflowTransferStep, CallFlowSayStep, CallFlowPlayStep,
// CallFlowPauseStep, CallFlowRecordStep, CallFlowFetchStep,
// CallFlowHangupStep.
//
// This interface is provided for clarity and not meant to be implemented by
// other (external) types.
type CallFlowStep interface {
	json.Marshaler
	json.Unmarshaler
}

// CallFlowStepBase contains all common properties for call flow steps.
type CallFlowStepBase struct {
	ID string `json:"id,omitempty"`

	OnKeypressGoto string `json:"onKeypressGoto,omitempty"`
	OnKeypressVar  string `json:"onKeypressVar,omitempty"`

	Conditions []struct {
		Variable string `json:"variable"`
		Operator string `json:"operator"`
		Value    string `json:"value"`
	} `json:"conditions,omitempty"`
}

// A CallFlowTransferStep transfers the call to a different phone/server.
type CallFlowTransferStep struct {
	CallFlowStepBase

	// The destination (E.164 formatted number or SIP URI) to transfer a call
	// to. E.g. 31612345678 or sip:foobar@example.com.
	Destination string

	// Record sets the the side of the call that should be recorded.
	//
	// Optional. Available options are "in", "out" and "both". "in" can be used
	// to record the voice of the destination, "out" records the source. "both"
	// records both source and destination individually.
	Record string
}

type jsonCallFlowTransferStep struct {
	CallFlowStepBase
	Action  string `json:"action"`
	Options struct {
		Destination string `json:"destination"`
		Record      string `json:"record,omitempty"`
	} `json:"options"`
}

// MarshalJSON implements the json.Marshaler interface.
func (step *CallFlowTransferStep) MarshalJSON() ([]byte, error) {
	data := jsonCallFlowTransferStep{}
	data.CallFlowStepBase = step.CallFlowStepBase
	data.Action = "transfer"
	data.Options.Destination = step.Destination
	data.Options.Record = step.Record
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (step *CallFlowTransferStep) UnmarshalJSON(data []byte) error {
	var raw jsonCallFlowTransferStep
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*step = CallFlowTransferStep{
		CallFlowStepBase: raw.CallFlowStepBase,
		Destination:      raw.Options.Destination,
		Record:           raw.Options.Record,
	}
	return nil
}

// A CallFlowSayStep pronounces a text message with a given voice and language.
type CallFlowSayStep struct {
	CallFlowStepBase

	// The text to pronounce.
	Payload string

	// The voice to use for pronouncing text.
	//
	// Allowed values: male, female.
	Voice string

	// The language of the text that is to be pronounced in <2 letter lang>-<2
	// letter country> format.
	//
	// Allowed values: cy-GB, da-DK, de-DE, en-AU, en-GB, en-GB-WLS, en-IN,
	// en-US, es-ES, es-US, fr-CA, fr-FR, is-IS, it-IT, ja-JP, nb-NO, nl-NL,
	// pl-PL, pt-BR, pt-PT, ro-RO, ru-RU, sv-SE, tr-TR.
	//
	// Please refer to the online documentation for an up to date list.
	Language string

	// The amount of times to repeat the payload.
	//
	// Allowed values: Between 1 and 10.
	Repeat int

	// Determines what happens if a machine picks up the phone.
	//
	// Possible values are: "continue": do not check, just play the message.
	// "delay": (default) if a machine answers, wait until the machine stops
	// talking. "hangup": Hangup when a machine answers.
	IfMachine string

	// The time to analyze if a machine has picked up the phone.
	//
	// Used in combination with the delay and hangup values of the ifMachine
	// attribute. Minimum: 400ms, maximum: 10s.
	//
	// Optional. The default is 7 seconds.
	//
	// Truncated to milliseconds
	MachineTimeout time.Duration
}

type jsonCallFlowSayStep struct {
	CallFlowStepBase
	Action  string `json:"action"`
	Options struct {
		Payload        string `json:"payload"`
		Voice          string `json:"voice,omitempty"`
		Language       string `json:"language,omitempty"`
		Repeat         int    `json:"Repeat,omitempty"`
		IfMachine      string `json:"ifMachine,omitempty"`
		MachineTimeout int    `json:"machineTimeout,omitempty"`
	} `json:"options"`
}

// MarshalJSON implements the json.Marshaler interface.
func (step *CallFlowSayStep) MarshalJSON() ([]byte, error) {
	data := jsonCallFlowSayStep{}
	data.CallFlowStepBase = step.CallFlowStepBase
	data.Action = "say"
	data.Options.Payload = step.Payload
	data.Options.Voice = step.Voice
	data.Options.Language = step.Language
	data.Options.Repeat = step.Repeat
	data.Options.IfMachine = step.IfMachine
	data.Options.MachineTimeout = int(step.MachineTimeout / time.Millisecond)
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (step *CallFlowSayStep) UnmarshalJSON(data []byte) error {
	var raw jsonCallFlowSayStep
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*step = CallFlowSayStep{
		CallFlowStepBase: raw.CallFlowStepBase,
		Payload:          raw.Options.Payload,
		Voice:            raw.Options.Voice,
		Language:         raw.Options.Language,
		Repeat:           raw.Options.Repeat,
		IfMachine:        raw.Options.IfMachine,
		MachineTimeout:   time.Duration(raw.Options.MachineTimeout) * time.Millisecond,
	}
	return nil
}

// A CallFlowPlayStep plays back an audio file.
type CallFlowPlayStep struct {
	CallFlowStepBase

	// The URL of the media file to play. The media file should be a WAV file
	// (8 kHz, 16 bit).
	Media string
}

type jsonCallFlowPlayStep struct {
	CallFlowStepBase
	Action  string `json:"action"`
	Options struct {
		Media string `json:"media"`
	} `json:"options"`
}

// MarshalJSON implements the json.Marshaler interface.
func (step *CallFlowPlayStep) MarshalJSON() ([]byte, error) {
	data := jsonCallFlowPlayStep{}
	data.CallFlowStepBase = step.CallFlowStepBase
	data.Action = "play"
	data.Options.Media = step.Media
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (step *CallFlowPlayStep) UnmarshalJSON(data []byte) error {
	var raw jsonCallFlowPlayStep
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*step = CallFlowPlayStep{
		CallFlowStepBase: raw.CallFlowStepBase,
		Media:            raw.Options.Media,
	}
	return nil
}

// A CallFlowPauseStep Pauses (silently) for a given duration.
type CallFlowPauseStep struct {
	CallFlowStepBase

	// The length of the pause.
	//
	// Truncated to seconds.
	Length time.Duration
}

type jsonCallFlowPauseStep struct {
	CallFlowStepBase
	Action  string `json:"action"`
	Options struct {
		Length int `json:"length"`
	} `json:"options"`
}

// MarshalJSON implements the json.Marshaler interface.
func (step *CallFlowPauseStep) MarshalJSON() ([]byte, error) {
	data := jsonCallFlowPauseStep{}
	data.CallFlowStepBase = step.CallFlowStepBase
	data.Action = "pause"
	data.Options.Length = int(step.Length / time.Second)
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (step *CallFlowPauseStep) UnmarshalJSON(data []byte) error {
	var raw jsonCallFlowPauseStep
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*step = CallFlowPauseStep{
		CallFlowStepBase: raw.CallFlowStepBase,
		Length:           time.Duration(raw.Options.Length) * time.Second,
	}
	return nil
}

// A CallFlowRecordStep initiates an audio recording during a call, e.g. for
// capturing a user response.
//
// Note that this is distinct from CallFlow.Record which records the entire
// call.
type CallFlowRecordStep struct {
	CallFlowStepBase

	// Setting the MaxLength limits the duration of the recording.
	//
	// Optional with the default value being 0 which imposes no limit.
	//
	// Truncated to seconds.
	MaxLength time.Duration

	// The duration of a moment of silence allowed before a recording is
	// stopped.
	//
	// If you omit this parameter, silence detection is disabled.
	//
	// Truncated to seconds.
	Timeout time.Duration

	// Key DTMF input to terminate recording.
	//
	// Values allowed are: "any", "#", "*", "none" (default).
	FinishOnKey string

	// Set this to get a transcription of a recording in the specified language
	// after the recording has finished.
	//
	// Allowed values: de-DE, en-AU, en-UK, en-US, es-ES, es-LA, fr-FR, it-IT, nl-NL, pt-BR.
	TranscribeLanguage string

	// (Optional) OnFinish contains the URL to get a new CallFlow from when the recording terminates and this CallFlowRecordStep ends.
	//
	// The URL must contain a schema e.g. http://... or https://...
	// This attribute is used for chaining call flows. When the current step ends,
	// a POST request containing information about the recording is sent to the URL specified.
	// This gets a new callflow from the URL specified, but re-uses the original Call ID and Leg ID i.e. it's the same Call.
	//
	// To get at the recording information from the POST request body, you must call (instead of relying on req.Form):
	// ```go
	// body,_ := ioutil.ReadAll(req.Body)
	// recordingInfo := string(body[:])
	// ```
	OnFinish string
}

type jsonCallFlowRecordStep struct {
	CallFlowStepBase
	Action  string `json:"action"`
	Options struct {
		MaxLength          int    `json:"maxLength"`
		Timeout            int    `json:"timeout"`
		FinishOnKey        string `json:"finishOnKey"`
		TranscribeLanguage string `json:"transcribeLanguage"`
		OnFinish           string `json:"onFinish"`
	} `json:"options"`
}

// MarshalJSON implements the json.Marshaler interface.
func (step *CallFlowRecordStep) MarshalJSON() ([]byte, error) {
	data := jsonCallFlowRecordStep{}
	data.CallFlowStepBase = step.CallFlowStepBase
	data.Action = "record"
	data.Options.MaxLength = int(step.MaxLength / time.Second)
	data.Options.Timeout = int(step.Timeout / time.Second)
	data.Options.FinishOnKey = step.FinishOnKey
	data.Options.TranscribeLanguage = step.TranscribeLanguage
	data.Options.OnFinish = step.OnFinish
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (step *CallFlowRecordStep) UnmarshalJSON(data []byte) error {
	var raw jsonCallFlowRecordStep
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*step = CallFlowRecordStep{
		CallFlowStepBase:   raw.CallFlowStepBase,
		MaxLength:          time.Duration(raw.Options.MaxLength) * time.Second,
		Timeout:            time.Duration(raw.Options.Timeout) * time.Second,
		FinishOnKey:        raw.Options.FinishOnKey,
		TranscribeLanguage: raw.Options.TranscribeLanguage,
		OnFinish:           raw.Options.OnFinish,
	}
	return nil
}

// A CallFlowFetchStep fetches a call flow from a remote URL.
//
// For more information on dynamic call flows, see:
// https://developers.messagebird.com/docs/voice-calling#dynamic-call-flows
//
// Any steps following this tag are ignored.
type CallFlowFetchStep struct {
	CallFlowStepBase

	// The URL to fetch the call flow from.
	URL string
}

type jsonCallFlowFetchStep struct {
	CallFlowStepBase
	Action  string `json:"action"`
	Options struct {
		URL string `json:"url"`
	} `json:"options"`
}

// MarshalJSON implements the json.Marshaler interface.
func (step *CallFlowFetchStep) MarshalJSON() ([]byte, error) {
	data := jsonCallFlowFetchStep{}
	data.CallFlowStepBase = step.CallFlowStepBase
	data.Action = "fetchCallFlow"
	data.Options.URL = step.URL
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (step *CallFlowFetchStep) UnmarshalJSON(data []byte) error {
	var raw jsonCallFlowFetchStep
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*step = CallFlowFetchStep{
		CallFlowStepBase: raw.CallFlowStepBase,
		URL:              raw.Options.URL,
	}
	return nil
}

// A CallFlowHangupStep ends the call.
type CallFlowHangupStep struct {
	CallFlowStepBase
}

type jsonCallFlowHangupStep struct {
	CallFlowStepBase
	Action string `json:"action"`
}

// MarshalJSON implements the json.Marshaler interface.
func (step *CallFlowHangupStep) MarshalJSON() ([]byte, error) {
	data := jsonCallFlowHangupStep{}
	data.CallFlowStepBase = step.CallFlowStepBase
	data.Action = "hangup"
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (step *CallFlowHangupStep) UnmarshalJSON(data []byte) error {
	var raw jsonCallFlowHangupStep
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*step = CallFlowHangupStep{
		CallFlowStepBase: raw.CallFlowStepBase,
	}
	return nil
}
