package voice

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

// A Webhook is an HTTP callback to your platform. They are sent when calls are
// created and updated.
type Webhook struct {
	ID        string
	URL       string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type jsonWebhook struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Token     string `json:"token"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// MarshalJSON implements the json.Marshaler interface.
func (wh Webhook) MarshalJSON() ([]byte, error) {
	data := jsonWebhook{
		ID:        wh.ID,
		URL:       wh.URL,
		Token:     wh.Token,
		CreatedAt: wh.CreatedAt.Format(time.RFC3339),
		UpdatedAt: wh.UpdatedAt.Format(time.RFC3339),
	}
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (wh *Webhook) UnmarshalJSON(data []byte) error {
	var raw jsonWebhook
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	createdAt, err := time.Parse(time.RFC3339, raw.CreatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse Webhook CreatedAt: %v", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, raw.UpdatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse Webhook UpdatedAt: %v", err)
	}
	*wh = Webhook{
		ID:        raw.ID,
		URL:       raw.URL,
		Token:     raw.Token,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return nil
}

// Webhooks returns a paginator over all webhooks.
func Webhooks(client *messagebird.Client) *Paginator {
	return newPaginator(client, "webhooks/", reflect.TypeOf(Webhook{}))
}

// CreateWebHook creates a new webhook the specified url that will be called
// and security token.
func CreateWebHook(client *messagebird.Client, url, token string) (*Webhook, error) {
	data := struct {
		URL   string `json:"url"`
		Token string `json:"token,omitempty"`
	}{
		URL:   url,
		Token: token,
	}
	wh := &Webhook{}
	if err := client.Request(wh, "POST", "webhooks/", data); err != nil {
		return nil, err
	}
	return wh, nil
}

// Update syncs hte local state of a webhook to the MessageBird API.
func (wh *Webhook) Update(client *messagebird.Client) error {
	var data struct {
		Data []Webhook `json:"data"`
	}
	if err := client.Request(&data, "PUT", "webhooks/"+wh.ID, wh); err != nil {
		return err
	}
	*wh = data.Data[0]
	return nil
}

// Delete deletes a webhook.
func (wh *Webhook) Delete(client *messagebird.Client) error {
	return client.Request(nil, "DELETE", "webhooks/"+wh.ID, nil)
}
