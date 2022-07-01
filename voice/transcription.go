package voice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	messagebird "github.com/messagebird/go-rest-api/v9"
)

// A Transcription is a textual representation of a recording as text.
//
// You can request an automated transcription for a recording by doing a POST
// request to the API.
type Transcription struct {
	// The unique ID of the transcription.
	ID string
	// The ID of the recording that the transcription belongs to.
	RecordingID string
	// The status of the transcription. Possible values: created, transcribing, done, failed.
	Status string
	// The date-time the transcription was created/requested.
	CreatedAt time.Time
	// The date-time the transcription was last updated.
	UpdatedAt time.Time

	// A hash with HATEOAS links related to the object. This includes the file
	// link that has the URI for downloading the text transcription of a
	// recording.
	links map[string]string
}

type jsonTranscription struct {
	ID          string            `json:"id"`
	RecordingID string            `json:"recordingID"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
	Links       map[string]string `json:"_links"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (trans *Transcription) UnmarshalJSON(data []byte) error {
	var raw jsonTranscription
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	createdAt, err := time.Parse(time.RFC3339, raw.CreatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse Transcription CreatedAt: %v", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, raw.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse Transcription UpdatedAt: %v", err)
	}
	*trans = Transcription{
		ID:          raw.ID,
		RecordingID: raw.RecordingID,
		Status:      raw.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		links:       raw.Links,
	}
	return nil
}

// Contents gets the transcription file.
//
// This is a plain text file.
func (trans *Transcription) Contents(client *messagebird.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, apiRoot+trans.links["file"], nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Authorization", "AccessKey "+client.AccessKey)
	req.Header.Set("User-Agent", "MessageBird/ApiClient/"+messagebird.ClientVersion+" Go/"+runtime.Version())

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad HTTP status: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return string(b), err
}

// CreateTranscription creates a transcription request for an existing recording
func CreateTranscription(client messagebird.MessageBirdClient, callID string, legID string, recordingID string) (trans *Transcription, err error) {
	var body struct{}
	path := fmt.Sprintf("/calls/%s/legs/%s/recordings/%s/transcriptions", callID, legID, recordingID)
	var resp struct {
		Data []Transcription `json:"data"`
	}
	if err := client.Request(&resp, http.MethodPost, apiRoot+path, body); err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	return &resp.Data[0], nil
}
