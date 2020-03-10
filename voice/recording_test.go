package voice

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/internal/mbtest"
)

func TestRecordingGetFile(t *testing.T) {
	fileContents := []byte("this is not really a WAV file")
	mbClient, stop := testRequest(http.StatusOK, []byte(fileContents))
	defer stop()

	rec := &Recording{
		ID: "1337",
		Links: map[string]string{
			"file": "/yolo/swag.wav",
		},
	}
	r, err := rec.DownloadFile(mbClient)
	if err != nil {
		t.Fatal(err)
	}
	wav, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(wav) != string(fileContents) {
		t.Logf("exp: %q", string(fileContents))
		t.Logf("got: %q", string(wav))
		t.Fatalf("mismatched downloaded contents")
	}
}

func TestReadRecording(t *testing.T) {
	mbtest.WillReturnTestdata(t, "recordingObject.json", http.StatusOK)
	client := mbtest.Client(t)

	recording, err := ReadRecording(client, "callid", "legid", "recid")
	if err != nil {
		t.Fatalf("unexpected error read recording: %s", err)
	}

	if recording.ID != "recid" {
		t.Fatalf("expect %s got %s", "recid", recording.ID)
	}
	if recording.LegID != "legid" {
		t.Fatalf("expect %s got %s", "legid", recording.LegID)
	}
	if recording.Status != RecordingStatusDone {
		t.Fatalf("expect %s got %s", RecordingStatusDone, recording.Status)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/calls/callid/legs/legid/recordings/recid")
}
