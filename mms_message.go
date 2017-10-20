package messagebird

import (
    "net/url"
    "time"
)

type MmsMessage struct {
    Id                string
    HRef              string
    Direction         string
    Originator        string
    Body              string
    Reference         string
    Subject           string
    MediaUrls         []string
    ScheduledDatetime *time.Time
    CreatedDatetime   *time.Time
    Recipients        Recipients
    Errors            []Error
}

type MmsMessageParams struct {
    Subject           string
    Reference         string
    ScheduledDatetime time.Time
}

func paramsForMmsMessage(params *MmsMessageParams) (*url.Values, error) {
    urlParams := &url.Values{}

    if params == nil {
        return urlParams, nil
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
