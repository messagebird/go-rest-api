package voice

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleCallFlow() {
	callflow := CallFlow{
		Steps: []CallFlowStep{
			&CallFlowSayStep{
				Payload:  "Hello",
				Voice:    "male",
				Language: "en-US",
			},
			&CallFlowPauseStep{
				Length: time.Second * 4,
			},
			&CallFlowTransferStep{
				Destination: "31600000000",
			},
			&CallFlowRecordStep{
				CallFlowStepBase:   CallFlowStepBase{},
				MaxLength:          10,
				Timeout:            5,
				FinishOnKey:        "#",
				Transcribe:         true,
				TranscribeLanguage: "en-US",
			},
		},
	}
	_ = callflow
}

func TestCallFlowJSONMarshal(t *testing.T) {
	refCreatedAt, _ := time.Parse(time.RFC3339, "2018-01-29T13:46:06Z")
	refUpdatedAt, _ := time.Parse(time.RFC3339, "2018-01-30T16:00:34Z")
	referenceCallflow := &CallFlow{
		ID: "id",
		Steps: []CallFlowStep{
			&CallFlowSayStep{
				CallFlowStepBase: CallFlowStepBase{
					ID: "1",
				},
				Payload:  "Hello",
				Language: "en-US",
				Voice:    "male",
			},
			&CallFlowPauseStep{
				CallFlowStepBase: CallFlowStepBase{
					ID: "2",
				},
				Length: time.Second * 10,
			},
			&CallFlowRecordStep{
				CallFlowStepBase: CallFlowStepBase{
					ID: "3",
				},
				MaxLength:          10,
				Timeout:            5,
				FinishOnKey:        "#",
				Transcribe:         true,
				TranscribeLanguage: "en-US",
			},
		},
		CreatedAt: refCreatedAt,
		UpdatedAt: refUpdatedAt,
	}

	jsonData, err := json.Marshal(referenceCallflow)
	assert.NoError(t, err)
	var callflow CallFlow
	unmarshallErr := json.Unmarshal(jsonData, &callflow)
	assert.NoError(t, unmarshallErr)
	if !reflect.DeepEqual(*referenceCallflow, callflow) {
		t.Logf("exp: %#v", *referenceCallflow)
		t.Logf("got: %#v", callflow)
		t.Fatalf("mismatched call flows")
	}
}

func TestCallFlowJSONUnmarshal(t *testing.T) {
	referenceJSON := `{
		 "id": "id",
		 "steps": [
			 {
				 "id": "1",
				 "action": "say",
				 "options": {
					 "payload": "Hello",
					 "language": "en-US",
					 "voice": "male"
				 }
			 },
			 {
				 "id": "2",
				 "action": "pause",
				 "options": {
					 "length": 10
				 }
			 },
			 {
				"id": "3",
				"action": "record",
				"options": {
					"maxLength": 10,
					"timeout": 5,
					"finishOnKey": "#",
					"transcribe": true,
					"transcribeLanguage": "en-US"
				}
			 }
		 ],
		 "createdAt": "2018-01-29T13:46:06Z",
		 "updatedAt": "2018-01-30T16:00:34Z"
	}`
	refCreatedAt, _ := time.Parse(time.RFC3339, "2018-01-29T13:46:06Z")
	refUpdatedAt, _ := time.Parse(time.RFC3339, "2018-01-30T16:00:34Z")
	referenceCallflow := &CallFlow{
		ID: "id",
		Steps: []CallFlowStep{
			&CallFlowSayStep{
				CallFlowStepBase: CallFlowStepBase{
					ID: "1",
				},
				Payload:  "Hello",
				Language: "en-US",
				Voice:    "male",
			},
			&CallFlowPauseStep{
				CallFlowStepBase: CallFlowStepBase{
					ID: "2",
				},
				Length: time.Second * 10,
			},
			&CallFlowRecordStep{
				CallFlowStepBase: CallFlowStepBase{
					ID: "3",
				},
				MaxLength:          10,
				Timeout:            5,
				FinishOnKey:        "#",
				Transcribe:         true,
				TranscribeLanguage: "en-US",
			},
		},
		CreatedAt: refCreatedAt,
		UpdatedAt: refUpdatedAt,
	}

	var callflow CallFlow
	err := json.Unmarshal([]byte(referenceJSON), &callflow)
	assert.NoError(t, err)
	if !reflect.DeepEqual(*referenceCallflow, callflow) {
		t.Logf("exp: %#v", *referenceCallflow)
		t.Logf("got: %#v", callflow)
		t.Fatalf("mismatched call flows")
	}
}

func TestCreateCallFlow(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	newCf := &CallFlow{
		Steps: []CallFlowStep{
			&CallFlowSayStep{
				Payload:  "Hello",
				Language: "en-US",
				Voice:    "male",
			},
			&CallFlowPauseStep{
				Length: time.Second,
			},
			&CallFlowRecordStep{
				CallFlowStepBase: CallFlowStepBase{
					ID: "3",
				},
				MaxLength:          10,
				Timeout:            5,
				FinishOnKey:        "#",
				Transcribe:         true,
				TranscribeLanguage: "en-US",
			},
		},
	}
	err := newCf.Create(mbClient)
	assert.NoError(t, err)
	assert.Len(t, newCf.Steps, 3)
}

func TestCallFlowByID(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	newCf := &CallFlow{
		Steps: []CallFlowStep{
			&CallFlowHangupStep{},
		},
	}
	err := newCf.Create(mbClient)
	assert.NoError(t, err)

	fetchedCf, err := CallFlowByID(mbClient, newCf.ID)
	assert.NoError(t, err)
	assert.Equal(t, newCf.ID, fetchedCf.ID)
}

func TestCallFlowList(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	for _, c := range []string{"foo", "bar", ""} {
		newCf := &CallFlow{
			Title: c,
			Steps: []CallFlowStep{
				&CallFlowHangupStep{},
			},
		}
		err := newCf.Create(mbClient)
		assert.NoError(t, err)
	}

	i := 0
	for cf := range CallFlows(mbClient).Stream() {
		_, ok := cf.(error)
		assert.False(t, ok)
		i++
	}
	assert.NotEqual(t, 0, i)
}
