package voice

import (
	"testing"
	"time"
)

func TestInitiateCall(t *testing.T) {
	client, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	source, destination := "31000000000", "31000000000"
	callflow := CallFlow{
		Title: "Say test",
		Steps: []CallFlowStep{
			&CallFlowSayStep{
				Voice:    "male",
				Payload:  "You are about to experience a great adventure which reaches from the inner mind to the outer limits",
				Language: "en-US",
			},
			&CallFlowPauseStep{
				Length: time.Second,
			},
			&CallFlowHangupStep{},
		},
	}
	_, err := InitiateCall(client, source, destination, callflow, nil)
	if err != nil {
		t.Fatal(err)
	}
}
