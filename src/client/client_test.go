package client

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewDefaultClient(t *testing.T) {
	client := NewDefaultClient()
	assert.True(t, true, reflect.DeepEqual(client, DefaultClient))
}

// unit-test-1a - test our setTimeout method, which sets the default req/resp timeout in the http client
func TestSetTimeoutFromFunc(t *testing.T) {
	t.Parallel()
	SetTimeout(123)
	assert.Equal(t, DefaultClient.HTTPClient.Timeout, (time.Duration(123) * time.Second), "failed to get timeout from http client")
}

// unit-test-1b - test setTimeout via the client struct
func TestSetTimeoutFromStruct(t *testing.T) {
	t.Parallel()
	fakeClient := &Client{
		HTTPClient: &http.Client{
			Timeout: time.Duration(10) * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).DialContext,
				MaxIdleConns:        100,
				MaxConnsPerHost:     100,
				MaxIdleConnsPerHost: 100,
			},
		},
	}
	fakeClient.SetTimeout(123)
	assert.Equal(t, fakeClient.HTTPClient.Timeout, (time.Duration(123) * time.Second), "failed to get timeout from http client")
}

// unit-test-2a - test context handling in http client. we are sleeping our mock httpServer (serving as our psuedo-API)
// and making sure the context times-out before the server is operational. we want to capture this case, when consuming custom contexts,
// espesically in a production environment
func TestSendWithCtxGoodRequest(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Millisecond * 10)
		fmt.Fprintln(writer, "{\"message\": \"superfakeapi\"}")
	}))
	defer mockServer.Close()

	req := Request{
		Method:  http.MethodGet,
		BaseURL: mockServer.URL + "/some_endpoint",
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*5)
	_, err := SendWithCtx(ctx, req)
	if err == nil {
		t.Error("a timeout exception based on the context passed in, should have trigered a timeout error")
	}
	assert.NotNil(t, err.Error(), "error from context timeout should not be nil")
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Error("we were expected a context timeout error but didn't recieve one")
	}
}

// unit-test-2b - testing send with context, but with a malformed request. expecting to capture an error
func TestSendWithCtxBadRequest(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(writer, "{\"message\": \"superfakeapi\"}")
	}))
	defer mockServer.Close()

	req := Request{
		Method:  "@",
		BaseURL: mockServer.URL + "/some_endpoint",
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10)
	_, err := SendWithCtx(ctx, req)
	if err == nil {
		t.Error("we expected this test to fail because we're passing a malformed request object")
	}
}

// unit-test-2c - tests the Send() function in the client package, end to end
func TestSend(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Millisecond * 10)
		fmt.Fprintln(writer, "{\"message\": \"superfakeapi\"}")
	}))
	defer mockServer.Close()

	req := Request{
		Method:  http.MethodGet,
		BaseURL: mockServer.URL + "/some_endpoint",
	}

	resp, err := Send(req)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 200, "invalid status code returned")
	assert.NotNil(t, resp.Body, "empty body should not have been returned")
	assert.NotNil(t, resp.Headers, "empty headers should not  have been returned")
}

// unit-test-2b - test outer Send() function in client package

// unit-test-3 - test our built-in URL query params helper function
func TestGenerateQueryParams(t *testing.T) {
	t.Parallel()
	query := make(map[string]string)
	hostname := "http://superfake.com"
	query["foo"] = "bar"
	query["foofoo"] = "barbar"
	generated := generateQueryParams(hostname, query)
	expected := "http://superfake.com?foo=bar&foofoo=barbar"
	assert.Equal(t, generated, expected, "generated and expected should have equaled each other")
}

// unit-test-4 - test our buildRequest method which transforms our Request strcut to something consumable by the http client
// in the form of a http.Request object
func TestBuildGoodRequest(t *testing.T) {
	t.Parallel()
	baseURL := "http://superfake.com"
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	query := make(map[string]string)
	query["foo"] = "bar"
	req := Request{
		Method:      http.MethodGet,
		BaseURL:     baseURL,
		QueryParams: query,
		Headers:     headers,
	}

	request, err := buildRequest(req)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, request)
	assert.Equal(t, request.Method, "GET", "buildRequest didn't set request method properly")
	assert.Equal(t, request.Header, http.Header(http.Header{"Accept": []string{"application/json"}}), "failed to set proper headers")
}

// unit-test-5 - test buildRequeset with a fudged-up Request object
func TestBuildBadRequest(t *testing.T) {
	t.Parallel()
	req := Request{
		Method: "@",
	}
	resp, err := buildRequest(req)
	if err == nil {
		t.Error("we expected this to error out")
	}
	assert.Nil(t, resp, "we expected resp to be nil")
}

// unit-test-6 - testing ExecuteRequest and BuildResponse. More of a functional test than a unit test.
func TestBuildGoodResponse(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(writer, "{\"message\": \"superfakeapi\"}")
	}))
	defer mockServer.Close()

	req := Request{
		Method:  http.MethodGet,
		BaseURL: mockServer.URL,
	}
	// step 1- build a http client compatible request object
	request, err := buildRequest(req)
	if err != nil {
		t.Error(err)
	}
	// step 2- execute the request and capture the raw http.Response object
	rawResponse, err := ExecuteRequest(request)
	if err != nil {
		t.Error(err)
	}
	// step 3- take the http.Response object from the previous step and run it through our Response object func
	resp, err := buildResponse(rawResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 200, "invalid status code returned")
	assert.NotNil(t, resp.Body, "empty body should not have been returned")
	assert.NotNil(t, resp.Headers, "empty headers should not  have been returned")
}

// some helper structs for unit-test-7
type panic struct{}

func (*panic) Read([]byte) (n int, e error) {
	return 0, errors.New("error")
}
func (*panic) Close() error {
	return nil
}

// unit-test-7 - force buildResponse to panic and capture event
func TestBuildBadResponse(t *testing.T) {
	t.Parallel()
	badResponse := &http.Response{
		Body: new(panic),
	}
	resp, err := buildResponse(badResponse)
	if err == nil {
		t.Error("bad response to buildResponse should have thrown an error")
	}
	assert.Nil(t, resp, "response should have been nil")
}

// unit-test-8 - using strict contextTimeouts to test a bad client situation.
// this should mimic a RL scenario where we are sending Requests to a dead / malformed client
func TestBadHTTPClient(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Millisecond * 10)
		fmt.Fprintln(writer, "{\"message\": \"superfakeapi\"}")
	}))
	defer mockServer.Close()

	req := Request{
		Method:  http.MethodGet,
		BaseURL: mockServer.URL,
	}

	mockClient := &Client{
		HTTPClient: &http.Client{
			Timeout: time.Duration(1) * time.Millisecond,
		},
	}

	resp, err := mockClient.Send(req)
	if err == nil {
		t.Error("this should have timed-out")
	}
	assert.Nil(t, resp, "response should have been nil")
}

// not running test 10 and 11 in parallel because there might be a race condition to DefaultClient.HTTPClient.Transport
// unit-test-10 - test SetDefaultClientTransportOpts
func TestDefaultClientTransportOpts(t *testing.T) {
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxIdleConns, 100, "default http transport settings impropely set")
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxConnsPerHost, 100, "default http transport settings impropely set")
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxIdleConnsPerHost, 100, "default http transport settings impropely set")
}

// unit-test-11 - test custom
func TestSetClientTransportOpts(t *testing.T) {
	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	SetClientTransportOpts(defaultTransport)
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxIdleConns, 100, "default http transport settings impropely set")
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxConnsPerHost, 0, "default http transport settings impropely set")
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxIdleConnsPerHost, 0, "default http transport settings impropely set")
}

// unit-test-12 - test setting custom transport (http.Transport) from client struct
func TestSetClientTransportOptsFromStruct(t *testing.T) {
	t.Parallel()
	fakeClient := &Client{
		HTTPClient: &http.Client{
			Timeout: time.Duration(10) * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).DialContext,
				MaxIdleConns:        100,
				MaxConnsPerHost:     100,
				MaxIdleConnsPerHost: 100,
			},
		},
	}
	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	fakeClient.SetClientTransportOpts(defaultTransport)
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxIdleConns, 100, "default http transport settings impropely set")
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxConnsPerHost, 0, "default http transport settings impropely set")
	assert.Equal(t, DefaultClient.HTTPClient.Transport.(*http.Transport).MaxIdleConnsPerHost, 0, "default http transport settings impropely set")
}

// unit-test-13 - set setting default Http headers when request body is present
func TestBuildRequestDefaultHeaders(t *testing.T) {
	t.Parallel()
	baseURL := "http://superfake.com"
	req := Request{
		Method:  http.MethodGet,
		BaseURL: baseURL,
		Body:    []byte("test"),
	}

	request, err := buildRequest(req)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, request)
	assert.Equal(t, request.Method, "GET", "buildRequest didn't set request method properly")
	assert.Equal(t, request.Header, http.Header(http.Header{"Content-Type": []string{"application/json"}}), "failed to set proper headers")
}
