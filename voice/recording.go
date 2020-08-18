package voice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v6"
)

// RecordingStatus enumerates all valid values for the status of a recording.
type RecordingStatus string

const (
	// RecordingStatusInitialised indicates that the recording has been created
	// but has not begun just yet.
	RecordingStatusInitialised RecordingStatus = "initialised"
	// RecordingStatusRecording indicates that recording is currently in progress.
	RecordingStatusRecording RecordingStatus = "recording"
	// RecordingStatusDone indicates that a recording is completed and may be downloaded.
	RecordingStatusDone RecordingStatus = "done"
	// RecordingStatusFailed indicates that something went wrong while
	// recording a leg.
	RecordingStatusFailed RecordingStatus = "failed"
)

// A Recording describes a voice recording of a leg.
//
// You can initiate a recording of a leg by having a step in your callflow with
// the record action set.
type Recording struct {
	// The unique ID of the recording.
	ID string
	// The format of the recording. Supported formats are: wav.
	Format string
	// The ID of the leg that the recording belongs to.
	LegID string
	// The status of the recording. Available statuses are: initialised, recording, done and failed
	Status RecordingStatus
	// The duration of the recording.
	//
	// Truncated to seconds.
	Duration time.Duration
	// The date-time the call was created.
	CreatedAt time.Time
	// The date-time the call was last updated.
	UpdatedAt time.Time

	// A hash with HATEOAS links related to the object. This includes the file
	// link that has the URI for downloading the wave file of the recording.
	Links map[string]string
}

type jsonRecording struct {
	ID        string            `json:"id"`
	Format    string            `json:"format"`
	LegID     string            `json:"legID"`
	Status    string            `json:"status"`
	Duration  int               `json:"duration"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
	Links     map[string]string `json:"_links"`
}

// ReadRecording fetches a single Recording based on its call ID, leg ID and the recording ID.
func ReadRecording(c *messagebird.Client, callID, legID, id string) (*Recording, error) {
	json := new(struct {
		Data []*Recording `json:"data"`
	})

	if err := c.Request(json, http.MethodGet, fmt.Sprintf("%s/calls/%s/legs/%s/recordings/%s",
		apiRoot, callID, legID, id), nil); err != nil {
		return nil, err
	}

	return json.Data[0], nil
}

// Recordings returns a Paginator which iterates over Recordings.
func Recordings(c *messagebird.Client, callID, legID string) *Paginator {
	return newPaginator(c, fmt.Sprintf("%s/calls/%s/legs/%s/recordings", apiRoot, callID,
		legID), reflect.TypeOf(Recording{}))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (rec *Recording) UnmarshalJSON(data []byte) error {
	recording := new(jsonRecording)

	if err := json.Unmarshal(data, recording); err != nil {
		return err
	}

	r, err := parseJSON(recording)
	if err != nil {
		return err
	}

	*rec = *r
	return nil
}

// Transcriptions returns a paginator for retrieving all Transcription objects.
func (rec *Recording) Transcriptions(client *messagebird.Client, callID string) *Paginator {
	path := apiRoot + rec.Links["self"] + "/transcriptions"
	return newPaginator(client, path, reflect.TypeOf(Transcription{}))
}

// Delete deletes a recording.
func Delete(client *messagebird.Client, callID, legID, recordingID string) error {
	return client.Request(nil, http.MethodDelete, fmt.Sprintf("%s/calls/%s/legs/%s/recordings/%s", apiRoot, callID, legID, recordingID), nil)
}

// DownloadFile streams the recorded WAV file.
func (rec *Recording) DownloadFile(client *messagebird.Client) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, apiRoot+rec.Links["file"], nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "audio/*")
	req.Header.Set("Authorization", "AccessKey "+client.AccessKey)
	req.Header.Set("User-Agent", "MessageBird/ApiClient/"+messagebird.ClientVersion+" Go/"+runtime.Version())

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad HTTP status: %d", resp.StatusCode)
	}
	return resp.Body, nil
}

func parseJSON(recording *jsonRecording) (*Recording, error) {
	createdAt, err := time.Parse(time.RFC3339, recording.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Recording CreatedAt: %v", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, recording.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Recording UpdatedAt: %v", err)
	}

	return &Recording{
		ID:        recording.ID,
		Format:    recording.Format,
		LegID:     recording.LegID,
		Status:    RecordingStatus(recording.Status),
		Duration:  time.Second * time.Duration(recording.Duration),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Links:     recording.Links,
	}, nil
}
