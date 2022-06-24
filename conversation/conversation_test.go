package conversation

import (
	"net/http"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api/v8/internal/mbtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mbtest.EnableServer(m)
}

func TestList(t *testing.T) {
	t.Run("limit_offset", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "conversationListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		convList, err := List(client, &ListOptions{10, 20})
		assert.NoError(t, err)

		assert.Equal(t, 20, convList.Offset)

		assert.Equal(t, "convid", convList.Items[0].ID)

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "limit=10&offset=20", query)
	})

	t.Run("all", func(t *testing.T) {
		mbtest.WillReturnTestdata(t, "allConversationListObject.json", http.StatusOK)
		client := mbtest.Client(t)

		convList, err := List(client, nil)
		assert.NoError(t, err)

		assert.Equal(t, 0, convList.Offset)

		assert.Equal(t, 10, convList.Limit)

		assert.Equal(t, "convid", convList.Items[0].ID)

		mbtest.AssertEndpointCalled(t, http.MethodGet, "/v1/conversations")

		query := mbtest.Request.URL.RawQuery
		assert.Equal(t, "", query)
	})
}

func TestRead(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationObject.json", http.StatusOK)
	client := mbtest.Client(t)

	conv, err := Read(client, "convid")
	assert.NoError(t, err)
	assert.Equal(t, "convid", conv.ID)
	assert.Equal(t, "contid", conv.Contact.ID)
	assert.Equal(t, "31612345678", conv.Contact.MSISDN)

	val, ok := conv.Contact.CustomDetails["userId"]
	assert.True(t, ok)
	assert.Equal(t, int64(12345678), val)
	assert.Equal(t, "chname", conv.Channels[0].Name)
	assert.Equal(t, 1, conv.Messages.TotalCount)
	assert.Equal(t, ConversationStatusActive, conv.Status)

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
					{
						Default: "Hello!",
					},
					{
						Default: "EUR12.34",
						Currency: &HSMLocalizableParameterCurrency{
							Code:   "EUR",
							Amount: 12340,
						},
					},
					{
						Default:  "Today",
						DateTime: mustParseRFC3339(t, "2018-08-24T11:52:12+00:00"),
					},
				},
			},
		},
		Type: MessageTypeHSM,
	})
	assert.NoError(t, err)
	assert.Equal(t, "convid", conv.ID)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/conversations/start")
	mbtest.AssertTestdata(t, "conversationStartHsmRequest.json", mbtest.Request.Body)
}

func mustParseRFC3339(t *testing.T, s string) *time.Time {
	result, err := time.Parse(time.RFC3339, s)
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	assert.Equal(t, "convid", conv.ID)

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
	assert.NoError(t, err)
	assert.Equal(t, "convid", conv.ID)

	mbtest.AssertEndpointCalled(t, http.MethodPost, "/v1/conversations/start")
	mbtest.AssertTestdata(t, "conversationStartTextRequest.json", mbtest.Request.Body)
}

func TestUpdate(t *testing.T) {
	mbtest.WillReturnTestdata(t, "conversationUpdatedObject.json", http.StatusOK)
	client := mbtest.Client(t)

	conv, err := Update(client, "id", &UpdateRequest{
		Status: ConversationStatusArchived,
	})
	assert.NoError(t, err)
	assert.Equal(t, ConversationStatusArchived, conv.Status)

	mbtest.AssertEndpointCalled(t, http.MethodPatch, "/v1/conversations/id")
	mbtest.AssertTestdata(t, "conversationUpdateRequest.json", mbtest.Request.Body)
}
