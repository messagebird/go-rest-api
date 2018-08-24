package conversation

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/internal/mbtest"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestList(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationListObject.json", http.StatusOK)
	client := mbtest.Client(t)

	convList, err := List(client, &ListOptions{10, 20})
	if err != nil {
		t.Fatalf("unexpected error listing Conversations: %s", err)
	}

	if convList.Offset != 20 {
		t.Fatalf("got %d, expected 20", convList.Offset)
	}

	if convList.Items[0].ID != "convid" {
		t.Fatalf("got %s, expected convid", convList.Items[0].ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations")
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationObject.json", http.StatusOK)
	client := mbtest.Client(t)

	conv, err := Read(client, "convid")
	if err != nil {
		t.Fatalf("unexpected error reading Conversation: %s", err)
	}

	if conv.ID != "convid" {
		t.Fatalf("got %s, expected convid", conv.ID)
	}

	if conv.Contact.ID != "contid" {
		t.Fatalf("got %s, expected contid", conv.Contact.ID)
	}

	if conv.Contact.MSISDN != "31612345678" {
		t.Fatalf("got %s, expected 31612345678", conv.Contact.MSISDN)
	}

	if val, ok := conv.Contact.CustomDetails["userId"]; ok {
		if val != int64(12345678) {
			t.Fatalf("got %v, expected 12345678", val)
		}
	} else {
		t.Fatalf("got nil, expected 12345678")
	}

	if conv.Channels[0].Name != "chname" {
		t.Fatalf("got %s, expected chname", conv.Channels[0].Name)
	}

	if conv.Messages.TotalCount != 1 {
		t.Fatalf("got %d, expected 1", conv.Messages.TotalCount)
	}

	if conv.Status != ConversationStatusActive {
		t.Fatalf("got %s, expected active", conv.Status)
	}

	mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations/convid")
}

func TestStartHSM(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	conv, err := Start(client, &StartRequest{
		ChannelID: "chid",
		To:        "31612345678",
		Content: &MessageContent{
			HSM: &HSM{
				Namespace:    "ns",
				TemplateName: "template",
				Language: &HSMLanguage{
					Policy: HSMLanguagePolicyDeterministic,
					Code:   "en_US",
				},
				LocalizableParameters: []*HSMLocalizableParameter{
					&HSMLocalizableParameter{
						Default: "Hello!",
					},
					&HSMLocalizableParameter{
						Default: "EUR12.34",
						Currency: &HSMLocalizableParameterCurrency{
							Code:   "EUR",
							Amount: 12340,
						},
					},
					&HSMLocalizableParameter{
						Default:  "Today",
						DateTime: mustParseRFC3339(t, "2018-08-24T11:52:12+00:00"),
					},
				},
			},
		},
		Type: MessageTypeHSM,
	})
	if err != nil {
		t.Fatalf("unexpected error starting Conversation: %s", err)
	}

	if conv.ID != "convid" {
		t.Fatalf("got %s, expected convid", conv.ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/conversations/start")
	mbtest.AssertTestdata(t, "conversationStartHsmRequest.json", mbtest.Request.Body)
}

func mustParseRFC3339(t *testing.T, s string) *time.Time {
	result, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("unexpected error parsing RFC3339 time: %s", err)
	}

	return &result
}

func TestStartMedia(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	conv, err := Start(client, &StartRequest{
		ChannelID: "chid",
		To:        "31612345678",
		Content: &MessageContent{
			Video: &Video{
				URL: "https://example.com/video.mp4",
			},
		},
		Type: MessageTypeText,
	})
	if err != nil {
		t.Fatalf("unexpected error starting Conversation: %s", err)
	}

	if conv.ID != "convid" {
		t.Fatalf("got %s, expected convid", conv.ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/conversations/start")
	mbtest.AssertTestdata(t, "conversationStartVideoRequest.json", mbtest.Request.Body)
}

func TestStartText(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationObject.json", http.StatusCreated)
	client := mbtest.Client(t)

	conv, err := Start(client, &StartRequest{
		ChannelID: "chid",
		To:        "31612345678",
		Content: &MessageContent{
			Text: "Hello",
		},
		Type: MessageTypeText,
	})
	if err != nil {
		t.Fatalf("unexpected error starting Conversation: %s", err)
	}

	if conv.ID != "convid" {
		t.Fatalf("got %s, expected convid", conv.ID)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/conversations/start")
	mbtest.AssertTestdata(t, "conversationStartTextRequest.json", mbtest.Request.Body)
}

func TestUpdate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationUpdatedObject.json", http.StatusOK)
	client := mbtest.Client(t)

	conv, err := Update(client, "id", &UpdateRequest{
		Status: ConversationStatusArchived,
	})
	if err != nil {
		t.Fatalf("unexpected error updating Conversation: %s", err)
	}

	if conv.Status != ConversationStatusArchived {
		t.Fatalf("got %s, expected archived", conv.Status)
	}

	mbtest.AssertEndpointCalled(t, http.MethodPatch, "/v1/conversations/id")
	mbtest.AssertTestdata(t, "conversationUpdateRequest.json", mbtest.Request.Body)
}
