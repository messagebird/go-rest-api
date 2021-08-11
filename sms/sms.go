package sms

import (
	"errors"
	"sync"
)

type SMSClient interface {
	Request(v interface{}, method, path string, data interface{}) error
}

var smsClient SMSClient
var mu sync.Mutex

func RegisterClient(c SMSClient) {
	mu.Lock()
	defer mu.Unlock()

	smsClient = c
}

var errNoClient = errors.New("no client is set")

func ensureClient() error {
	mu.Lock()
	defer mu.Unlock()

	if smsClient == nil {
		return errNoClient
	}

	return nil
}
