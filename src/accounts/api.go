// example implementation of code that can consume the go package to interface with the form3 backend accounts api

package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/sarabrajsingh/interview-accountapi/src/client"
	. "github.com/sarabrajsingh/interview-accountapi/src/models"
	"strconv"
	"os"
)

type URL struct {
	BaseURL string
}

var DefaultUrl = &URL{
	BaseURL: os.Getenv("FORM3_ACCOUNTS_API_URL"),
}

// helper function for the URL struct used in this module
func (u *URL) SetBaseURL(url string) {
	u.BaseURL = url
}

// CREATE account without custom user context
func Create(acc Account) (*Response, error) {
	accEncoded, err := json.Marshal(acc)

	if err != nil {
		return nil, err
	}

	return Send(Request{
		Method:  POST,
		BaseURL: DefaultUrl.BaseURL,
		Body:    accEncoded,
	})
}

// there is no method overloading in golang, nor are there default params like in python, so we need to create
// concrete methods that encompass all functionalities
// Create with custom context
func CreateWithCtx(ctx context.Context, acc Account) (*Response, error) {
	accEncoded, err := json.Marshal(acc)
	if err != nil {
		return nil, err
	}
	return SendWithCtx(ctx, Request{
		Method:  POST,
		BaseURL: DefaultUrl.BaseURL,
		Body:    accEncoded,
	})
}

// fetch implementation
func Fetch(id string) (*Response, error) {
	return Send(Request{
		Method:  GET,
		BaseURL: fmt.Sprintf("%s/%s", DefaultUrl.BaseURL, id),
	})
}
// fetch with context implementation
func FetchWithCtx(ctx context.Context, id string) (*Response, error) {
	return SendWithCtx(ctx, Request{
		Method:  GET,
		BaseURL: fmt.Sprintf("%s/%s", DefaultUrl.BaseURL, id),
	})
}

// delete implementation
func Delete(id string, version int) (*Response, error) {
	return Send(Request{
		Method:  DELETE,
		BaseURL: fmt.Sprintf("%s/%s", DefaultUrl.BaseURL, id),
		QueryParams: map[string]string{
			"version": strconv.Itoa(version),
		},
	})
}

// delete with context implementation
func DeleteWithCtx(ctx context.Context, id string, version int) (*Response, error) {
	return SendWithCtx(ctx, Request{
		Method:  DELETE,
		BaseURL: fmt.Sprintf("%s/%s", DefaultUrl.BaseURL, id),
		QueryParams: map[string]string{
			"version": strconv.Itoa(version),
		},
	})
}
