package voice

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func ExampleCallFlow() {
	callflow := CallFlow{
		Title: "My CallFlow",
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
		ID:    "id",
		Title: "title",
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
				CallFlowStepBase:   CallFlowStepBase{
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
	if err != nil {
		t.Fatal(err)
	}
	var callflow CallFlow
	if err := json.Unmarshal(jsonData, &callflow); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*referenceCallflow, callflow) {
		t.Logf("exp: %#v", *referenceCallflow)
		t.Logf("got: %#v", callflow)
		t.Fatalf("mismatched call flows")
	}
}

func TestCallFlowJSONUnmarshal(t *testing.T) {
	referenceJSON := `{
		 "id": "id",
		 "title": "title",
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
		ID:    "id",
		Title: "title",
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
				CallFlowStepBase:   CallFlowStepBase{
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
	if err := json.Unmarshal([]byte(referenceJSON), &callflow); err != nil {
		t.Fatal(err)
	}
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
		Title: "the-title",
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
				CallFlowStepBase:   CallFlowStepBase{
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
	if err := newCf.Create(mbClient); err != nil {
		t.Fatal(err)
	}
	if newCf.Title != "the-title" {
		t.Fatalf("Unexpected Title: %q", newCf.Title)
	}
	if len(newCf.Steps) != 3 {
		t.Fatalf("Unexpected number of steps: %q", len(newCf.Steps))
	}
}

func TestCallFlowByID(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	newCf := &CallFlow{
		Title: "the-title",
		Steps: []CallFlowStep{
			&CallFlowHangupStep{},
		},
	}
	if err := newCf.Create(mbClient); err != nil {
		t.Fatal(err)
	}

	fetchedCf, err := CallFlowByID(mbClient, newCf.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetchedCf.ID != newCf.ID {
		t.Fatalf("mismatched fetched IDs: exp %q, got %q", newCf.ID, fetchedCf.ID)
	}
}

func TestCallFlowList(t *testing.T) {
	mbClient, ok := testClient(t)
	if !ok {
		t.SkipNow()
	}

	for _, c := range []string{"foo", "bar", "baz"} {
		newCf := &CallFlow{
			Title: c,
			Steps: []CallFlowStep{
				&CallFlowHangupStep{},
			},
		}
		if err := newCf.Create(mbClient); err != nil {
			t.Fatal(err)
		}
	}

	i := 0
	for cf := range CallFlows(mbClient).Stream() {
		if err, ok := cf.(error); ok {
			t.Fatal(err)
		}
		i++
	}
	if i == 0 {
		t.Fatal("no callflows were fetched")
	}
}
