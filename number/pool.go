package number

import (
	"fmt"
	messagebird "github.com/messagebird/go-rest-api/v9"
	"net/http"
	"strings"
	"time"
)

type Pool struct {
	ID            string
	Name          string
	Service       string
	Configuration *PoolConfiguration
	NumbersCount  int
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type PoolConfiguration struct {
	ByCountry bool `json:"byCountry"`
}

type Pools struct {
	Limit, Offset, Count, TotalCount int
	Items                            []*Pool
}

type PoolNumbers struct {
	Limit, Offset, Count, TotalCount int
	Numbers                          []string
}

type AddNumberToPollResult struct {
	Success []string
	Fail    []FailResult
}

type FailResult struct {
	Number, Error string
}

type CreatePoolRequest struct {
	PoolName      string             `json:"poolName"`
	Service       string             `json:"service"`
	Configuration *PoolConfiguration `json:"configuration"`
}

type UpdatePoolRequest struct {
	PoolName      string             `json:"poolName"` // new pool name
	Configuration *PoolConfiguration `json:"configuration"`
}

type ListPoolRequest struct {
	PoolName string `json:"poolName,omitempty"` // new pool name
	Service  string `json:"service,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

type ListPoolNumbersRequest struct {
	Number string `json:"number,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

func CreatePool(c messagebird.MessageBirdClient, req *CreatePoolRequest) (*Pool, error) {
	p := &Pool{}
	if err := request(c, p, http.MethodPost, pathPools, req); err != nil {
		return nil, err
	}

	return p, nil
}

func ReadPool(c messagebird.MessageBirdClient, poolName string) (*Pool, error) {
	uri := fmt.Sprintf("%s/%s", pathPools, poolName)

	p := &Pool{}
	if err := request(c, p, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return p, nil
}

func UpdatePool(c messagebird.MessageBirdClient, poolName string, req *UpdatePoolRequest) (*Pool, error) {
	uri := fmt.Sprintf("%s/%s", pathPools, poolName)

	p := &Pool{}
	if err := request(c, p, http.MethodPut, uri, req); err != nil {
		return nil, err
	}

	return p, nil
}

func DeletePool(c messagebird.MessageBirdClient, poolName string) error {
	uri := fmt.Sprintf("%s/%s", pathPools, poolName)

	return request(c, nil, http.MethodDelete, uri, nil)
}

func ListPool(c messagebird.MessageBirdClient, req *ListPoolRequest) (*Pools, error) {
	p := &Pools{}
	if err := request(c, p, http.MethodGet, pathPools, req); err != nil {
		return nil, err
	}

	return p, nil
}

func ListPoolNumbers(c messagebird.MessageBirdClient, poolName string, req *ListPoolNumbersRequest) (*PoolNumbers, error) {
	uri := fmt.Sprintf("%s/%s/%s", pathPools, poolName, pathNumbers)

	p := &PoolNumbers{}
	if err := request(c, p, http.MethodGet, uri, req); err != nil {
		return nil, err
	}

	return p, nil
}

func AddNumberToPool(c messagebird.MessageBirdClient, poolName string, numbers []string) (*AddNumberToPollResult, error) {
	uri := fmt.Sprintf("%s/%s/%s", pathPools, poolName, pathNumbers)

	req := &struct {
		numbers []string
	}{numbers}

	p := &AddNumberToPollResult{}
	if err := request(c, p, http.MethodPost, uri, req); err != nil {
		return nil, err
	}

	return p, nil
}

func DeleteNumberFromPool(c messagebird.MessageBirdClient, poolName string, numbers []string) error {
	uri := fmt.Sprintf("%s/%s/%s", pathPools, poolName, pathNumbers)

	req := &struct {
		numbers string
	}{strings.Join(numbers, ",")}

	return request(c, nil, http.MethodDelete, uri, req)
}
