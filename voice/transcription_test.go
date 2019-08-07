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
	mbClient, stop := testRequest(http.StatusOK, []byte(`{
		"data": [
			{
				"id": "00000000-1111-2222-3333-444444444444",
				"recordingId": "55555555-6666-7777-8888-999999999999",
				"error": null,
				"createdAt": "2011-01-01T02:03:04Z",
				"updatedAt": "2011-01-02T03:04:05Z"
			}
		]
	}`))
	defer stop()

	callID, legID, recordingID := "7777777", "88888888", "999999999"

	_, err := CreateTranscription(mbClient, callID, legID, recordingID)
	if err != nil {
		t.Fatal(err)
	}
}
