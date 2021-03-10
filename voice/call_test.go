package voice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
}

func TestCallByID(t *testing.T) {
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
	call, err := InitiateCall(client, source, destination, callflow, nil)
	assert.NoError(t, err)

	time.Sleep(time.Second)
	fetchedCall, err := CallByID(client, call.ID)
	assert.NoError(t, err)
	assert.Equal(t, call.Source, fetchedCall.Source)
}
