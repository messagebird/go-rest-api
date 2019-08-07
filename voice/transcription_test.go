package voice

import (
	"net/http"
	"testing"
)

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
	if err != nil {
		t.Fatal(err)
	}
	if contents != text {
		t.Logf("exp: %q", text)
		t.Logf("got: %q", contents)
		t.Fatalf("mismatched downloaded contents")
	}
}

func TestCreateTranscription(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	callID, legID, recordingID := "7777777", "88888888", "999999999"

	_, err := CreateTranscription(mbClient, callID, legID, recordingID)
	if err != nil {
		t.Fatal(err)
	}
}
