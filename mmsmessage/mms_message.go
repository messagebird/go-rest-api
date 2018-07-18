package mmsmessage

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

// MMSMessage represents a MMS Message.
type MMSMessage struct {
	ID                string
	HRef              string
	Direction         string
	Originator        string
	Body              string
	Reference         string
	Subject           string
	MediaUrls         []string
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        messagebird.Recipients
}

// MMSMessageParams represents the parameters that can be supplied when creating
// a request.
type MMSMessageParams struct {
	Body              string
	MediaUrls         []string
	Subject           string
	Reference         string
	ScheduledDatetime time.Time
}

// path represents the path to the MMS resource.
const path = "mms"

// Read retrieves the information of an existing MmsMessage.
func Read(c *messagebird.Client, id string) (*MMSMessage, error) {
	mmsMessage := &MMSMessage{}
	if err := c.Request(mmsMessage, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return mmsMessage, nil
}

// Create creates a new MMS message for one or more recipients.
func Create(c *messagebird.Client, originator string, recipients []string, msgParams *MMSMessageParams) (*MMSMessage, error) {
	params, err := paramsForMMSMessage(msgParams)
	if err != nil {
		return nil, err
	}

	params.Set("originator", originator)
	params.Set("recipients", strings.Join(recipients, ","))

	mmsMessage := &MMSMessage{}
	if err := c.Request(mmsMessage, http.MethodPost, path, params); err != nil {
		return nil, err
	}

	return mmsMessage, nil
}

// paramsForMMSMessage converts the specified MMSMessageParams struct to a
// url.Values pointer and returns it.
func paramsForMMSMessage(params *MMSMessageParams) (*url.Values, error) {
	urlParams := &url.Values{}

	if params.Body == "" && params.MediaUrls == nil {
		return nil, errors.New("Body or MediaUrls is required")
	}
	if params.Body != "" {
		urlParams.Set("body", params.Body)
	}
	if params.MediaUrls != nil {
		urlParams.Set("mediaUrls[]", strings.Join(params.MediaUrls, ","))
	}
	if params.Subject != "" {
		urlParams.Set("subject", params.Subject)
	}
	if params.Reference != "" {
		urlParams.Set("reference", params.Reference)
	}
	if params.ScheduledDatetime.Unix() > 0 {
		urlParams.Set("scheduledDatetime", params.ScheduledDatetime.Format(time.RFC3339))
	}

	return urlParams, nil
}