package voice

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/messagebird/go-rest-api/v6/internal/mbtest"
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

func TestDelete(t *testing.T) {
	mbtest.WillReturn([]byte(""), http.StatusNoContent)
	client := mbtest.Client(t)
	if err := Delete(client, "callid", "legid", "recid"); err != nil {
		t.Errorf("unexpected error while deleting recording: %s", err)
	}

	mbtest.AssertEndpointCalled(t, http.MethodDelete, "/v1/calls/callid/legs/legid/recordings/recid")
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

func TestRecordings(t *testing.T) {
	mbtest.WillReturnTestdata(t, "recordingPaginatorObject.json", http.StatusOK)
	client := mbtest.Client(t)

	paginator := Recordings(client, "callid", "legid")

	data, err := paginator.NextPage()
	if err != nil {
		t.Fatalf("unexpected error read recording: %s", err)
	}

	recordings := data.([]Recording)

	if len(recordings) != 2 {
		t.Errorf("got %d recordings expect 1", len(recordings))
	}

	if recordings[0].ID != "recid" {
		t.Errorf("expect %s got %s", "recid", recordings[0].ID)
	}
	if recordings[0].LegID != "legid" {
		t.Errorf("expect %s got %s", "legid", recordings[0].LegID)
	}
	if recordings[0].Status != RecordingStatusDone {
		t.Errorf("expect %s got %s", RecordingStatusDone, recordings[0].Status)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/calls/callid/legs/legid/recordings")
}
