package number

import (
	"fmt"
	messagebird "github.com/messagebird/go-rest-api/v9"
	"net/http"
)

type BackOrderID string

type Backorder struct {
	ID          string
	ProductID   int
	Country     string
	Prefix      string
	Status      string
	ReasonCodes []string
}

type BackorderDocument struct {
	ID          int
	Name        string
	Description string
	Status      string
}

type BackorderDocuments struct {
	Limit, Count int
	Items        []*BackorderDocument
}

type EndUserDetail struct {
	ID    string
	Label string
}

type EndUserDetails struct {
	Items []*EndUserDetail
}

type PlaceBackorderRequest struct {
	ProductID int    `json:"productID"`
	Prefix    string `json:"prefix"`
	Quantity  int    `json:"quantity"`
}

type CreateBackorderDocumentRequest struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	Content  string `json:"content"`
}

type CreateBackorderEndUserDetailRequest struct {
	CompanyName  string
	Street       string
	StreetNumber string
	ZipCode      string
	City         string
	Country      string
}

func PlaceBackorder(c messagebird.MessageBirdClient, req *PlaceBackorderRequest) (BackOrderID, error) {
	resp := &struct {
		Id string `json:"id"`
	}{}

	if err := request(c, resp, http.MethodPost, pathBackorders, req); err != nil {
		return "", err
	}

	return BackOrderID(resp.Id), nil
}

func ReadBackorder(c messagebird.MessageBirdClient, backOrderID string) (*Backorder, error) {
	uri := fmt.Sprintf("%s/%s", pathBackorders, backOrderID)

	bo := &Backorder{}
	if err := request(c, bo, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return bo, nil
}

func ListBackorderDocuments(c messagebird.MessageBirdClient, backOrderID string) (*BackorderDocuments, error) {
	uri := fmt.Sprintf("%s/%s/%s", pathBackorders, backOrderID, pathDocuments)

	bd := &BackorderDocuments{}
	if err := request(c, bd, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return bd, nil
}

func CreateBackorderDocument(c messagebird.MessageBirdClient, backOrderID string, req *CreateBackorderDocumentRequest) error {
	uri := fmt.Sprintf("%s/%s/%s", pathBackorders, backOrderID, pathDocuments)

	return request(c, nil, http.MethodPost, uri, req)
}

func ListBackorderEndUserDetails(c messagebird.MessageBirdClient, backOrderID string) (*EndUserDetails, error) {
	uri := fmt.Sprintf("%s/%s/%s", pathBackorders, backOrderID, pathEndUserDetails)

	eud := &EndUserDetails{}
	if err := request(c, eud, http.MethodGet, uri, nil); err != nil {
		return nil, err
	}

	return eud, nil
}

func CreateBackorderEndUserDetail(c messagebird.MessageBirdClient, backOrderID string, req *CreateBackorderEndUserDetailRequest) error {
	uri := fmt.Sprintf("%s/%s/%s", pathBackorders, backOrderID, pathEndUserDetails)

	return request(c, nil, http.MethodPost, uri, req)
}
