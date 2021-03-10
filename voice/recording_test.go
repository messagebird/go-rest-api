package voice

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestRecordingGetFile(t *testing.T) {
	fileContents := []byte("this is not really a WAV file")
	mbClient, stop := testRequest(http.StatusOK, fileContents)
	defer stop()

	rec := &Recording{
		ID: "1337",
		Links: map[string]string{
			"file": "/yolo/swag.wav",
		},
	}
	r, err := rec.DownloadFile(mbClient)
	assert.NoError(t, err)
	wav, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	if string(wav) != string(fileContents) {
		t.Logf("exp: %q", string(fileContents))
		t.Logf("got: %q", string(wav))
		t.Fatalf("mismatched downloaded contents")
	}
}

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)
	err := Delete(client, "callid", "legid", "recid")
	assert.NoError(t, err)

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/calls/callid/legs/legid/recordings/recid")
}

func TestReadRecording(t *testing.T) {
	mbtest.WillReturnTestdata(t, "recordingObject.json", http.StatusOK)
	client := mbtest.Client(t)

	recording, err := ReadRecording(client, "callid", "legid", "recid")
	assert.NoError(t, err)
	assert.Equal(t, "recid", recording.ID)
	assert.Equal(t, "legid", recording.LegID)
	assert.Equal(t, RecordingStatusDone, recording.Status)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/calls/callid/legs/legid/recordings/recid")
}

func TestRecordings(t *testing.T) {
	mbtest.WillReturnTestdata(t, "recordingPaginatorObject.json", http.StatusOK)
	client := mbtest.Client(t)

	paginator := Recordings(client, "callid", "legid")

	data, err := paginator.NextPage()
	assert.NoError(t, err)

	recordings := data.([]Recording)
	assert.Len(t, recordings, 2)
	assert.Equal(t, "recid", recordings[0].ID)
	assert.Equal(t, "legid", recordings[0].LegID)
	assert.Equal(t, RecordingStatusDone, recordings[0].Status)

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/calls/callid/legs/legid/recordings")
}
