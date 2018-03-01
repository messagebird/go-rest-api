package voice

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRecordingGetFile(t *testing.T) {
	fileContents := []byte("this is not really a WAV file")
	mbClient, stop := testClient(http.StatusOK, []byte(fileContents))
	defer stop()

	rec := &Recording{
		ID: "1337",
		links: map[string]string{
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
