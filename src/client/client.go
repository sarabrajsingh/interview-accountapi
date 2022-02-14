package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// default http client that does all the http work for us. from the golang stl
var DefaultClient = &Client{
	HTTPClient: &http.Client{
		Timeout: time.Duration(10) * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			MaxIdleConns: 100,
			MaxConnsPerHost: 100,
			MaxIdleConnsPerHost: 100,
		},
	},
}

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	DELETE Method = "DELETE"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
)

type Request struct {
	Method      Method
	BaseURL     string
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       string
}

// struct around the main http engine. if a programmer needs to override any other environmental parameter for the client
// they can leverage this struct to do so
type Client struct {
	HTTPClient *http.Client
}

// helper functions to set http client timeouts
func SetTimeout(timeout int) {
	DefaultClient.HTTPClient.Timeout = time.Duration(timeout) * time.Second
}

func (c *Client) SetTimeout(timeout int) {
	c.HTTPClient.Timeout = (time.Duration(timeout) * time.Second)
}

// helper functions to set http client Transport options
func SetClientTransportOpts(t *http.Transport) {
	DefaultClient.HTTPClient.Transport = t
}

func (c *Client) SetClientTransportOpts(t *http.Transport) {
	c.HTTPClient.Transport = t
}

// helper function that generates URL encoded query params to a http request
func generateQueryParams(baseURL string, query map[string]string) string {
	baseURL += "?"
	parameters := url.Values{}
	for key, value := range query {
		parameters.Add(key, value)
	}
	return baseURL + parameters.Encode()
}

// transforms our custom Request struct to a http.Request object that can be consumed by the http client
func buildRequest(r Request) (*http.Request, error) {
	if len(r.QueryParams) != 0 {
		r.BaseURL = generateQueryParams(r.BaseURL, r.QueryParams)
	}
	
	// generate our http client compatible http.Request object. canonical pattern to send HTTP requests to a http client
	req, err := http.NewRequest(string(r.Method), r.BaseURL, bytes.NewReader(r.Body))

	if err != nil {
		return req, err
	}

	// set our headers if any
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	// a small check to enforce a Content-Type: application/json in our request headers 
	_, val := req.Header["Content-Type"]

	if len(r.Body) > 0 && !val {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, err
}

// public facing method that takes a http.Request object, marshals it to the default client in this module,
// and returns a raw http.Response
func ExecuteRequest(r *http.Request) (*http.Response, error) {
	return DefaultClient.HTTPClient.Do(r)
}

// public facing handler to the Client struct that emulates the function above
func (c *Client) ExecuteRequest(r *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(r)
}

// internal function that transforms a raw http.Response object from the http client to our consumable and custom defined Response object
func buildResponse(r *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	// transform a raw http.Response object to a marshalled Response object
	response := Response{
		StatusCode: r.StatusCode,
		Headers:    r.Header,
		Body:       string(body),
	}

	// must always close connection when using a io-op
	r.Body.Close()

	return &response, err
}

// public facing callable module function that sends Request objects to http client with a default context
func Send(r Request) (*Response, error) {
	return SendWithCtx(context.Background(), r)
}

// public facing wrapper for *Client.sendWithCtx()
func SendWithCtx(ctx context.Context, r Request) (*Response, error) {
	return DefaultClient.sendWithCtx(ctx, r)
}

// client send functionality sends with a default context timeout
func (c *Client) Send(r Request) (*Response, error) {
	return c.sendWithCtx(context.Background(), r)
}

// this function allows the caller to override the context that gets passed to the http client. called by SendWithCtx
func (c *Client) sendWithCtx(ctx context.Context, r Request) (*Response, error) {
	request, err := buildRequest(r)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	result, err := c.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}
	return buildResponse(result)
}
