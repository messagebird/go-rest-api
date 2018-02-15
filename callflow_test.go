package messagebird

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

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
	SetServerResponse(http.StatusOK, []byte(`{
		"data": [
			{
				"id": "the-id",
				"title": "the-title",
				"createdAt": "2018-01-29T13:46:06Z",
				"updatedAt": "2018-01-30T16:00:34Z"
			}
		]
	}`))
	newCf, err := mbClient.CreateCallFlow(&CallFlow{})
	if err != nil {
		t.Fatal(err)
	}
	if newCf.ID != "the-id" {
		t.Fatalf("Unexpected ID: %q", newCf.ID)
	}
	if newCf.Title != "the-title" {
		t.Fatalf("Unexpected Title: %q", newCf.Title)
	}
}
