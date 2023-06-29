// example implementation of code that can consume the go package to interface with the form3 backend accounts api

package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/sarabrajsingh/interview-accountapi/src/client"
	"github.com/sarabrajsingh/interview-accountapi/src/models"
)

type URL struct {
	BaseURL string
}

var DefaultUrl URL

func (u *URL) defaultBaseURL() {
	url := os.Getenv("FORM3_ACCOUNTS_API_URL")
	if url == "" {
		url = "http://localhost:8080/v1/organisation/accounts"
	}
	u.BaseURL = url
}

func (u *URL) GetDefaultBaseURL() string {
	u.defaultBaseURL()
	return u.BaseURL
}

// helper function for the URL struct used in this module
func (u *URL) SetBaseURL(url string) {
	u.BaseURL = url
}

// CREATE account without custom user context
func Create(acc models.Account) (*client.Response, error) {
	accEncoded, err := json.Marshal(acc)

	if err != nil {
		return nil, err
	}

	return client.Send(client.Request{
		Method:  http.MethodPost,
		BaseURL: DefaultUrl.GetDefaultBaseURL(),
		Body:    accEncoded,
	})
}

// there is no method overloading in golang, nor are there default params like in python, so we need to create
// concrete methods that encompass all functionalities
// Create with custom context
func CreateWithCtx(ctx context.Context, acc models.Account) (*client.Response, error) {
	accEncoded, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}
	return client.SendWithCtx(ctx, client.Request{
		Method:  http.MethodPost,
		BaseURL: DefaultUrl.GetDefaultBaseURL(),
		Body:    accEncoded,
	})
}

// fetch implementation
func Fetch(id string) (*client.Response, error) {
	return client.Send(client.Request{
		Method:  http.MethodGet,
		BaseURL: fmt.Sprintf("%s/%s", DefaultUrl.GetDefaultBaseURL(), id),
	})
}

// fetch with context implementation
func FetchWithCtx(ctx context.Context, id string) (*client.Response, error) {
	return client.SendWithCtx(ctx, client.Request{
		Method:  http.MethodGet,
		BaseURL: fmt.Sprintf("%s/%s", DefaultUrl.GetDefaultBaseURL(), id),
	})
}

// delete implementation
func Delete(id string, version int) (*client.Response, error) {
	return client.Send(client.Request{
		Method:  http.MethodDelete,
		BaseURL: fmt.Sprintf("%s/%s", DefaultUrl.GetDefaultBaseURL(), id),
		QueryParams: map[string]string{
			"version": strconv.Itoa(version),
		},
	})
}

// delete with context implementation
func DeleteWithCtx(ctx context.Context, id string, version int) (*client.Response, error) {
	return client.SendWithCtx(ctx, client.Request{
		Method:  http.MethodDelete,
		BaseURL: fmt.Sprintf("%s/%s", DefaultUrl.GetDefaultBaseURL(), id),
		QueryParams: map[string]string{
			"version": strconv.Itoa(version),
		},
	})
}
