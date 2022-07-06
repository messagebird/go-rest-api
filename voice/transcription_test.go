package voice

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/v9/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestTranscriptionGetContents(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog"
	mbClient, stop := testRequest(http.StatusOK, []byte(text))
	defer stop()

	trans := &Transcription{
		ID: "1337",
		links: map[string]string{
			"file": "/yolo/swag.txt",
		},
	}
	contents, err := trans.Contents(mbClient)
	assert.NoError(t, err)
	if contents != text {
		t.Logf("exp: %q", text)
		t.Logf("got: %q", contents)
		t.Fatalf("mismatched downloaded contents")
	}
}

func TestCreateTranscription(t *testing.T) {
	mbtest.WillReturnTestdata(t, "transcriptObject.json", http.StatusOK)
	client := mbtest.Client(t)

	callID, legID, recordingID := "7777777", "88888888", "999999999"
	_, err := CreateTranscription(client, callID, legID, recordingID)
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodPost, fmt.Sprintf("/v1/calls/%s/legs/%s/recordings/%s/transcriptions", callID, legID, recordingID))
}
