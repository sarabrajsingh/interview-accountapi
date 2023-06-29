# Form3 Take Home Exercise

## Purpose

Sarabraj Singh's submission to the form3-interview process. My design philosophy was to keep the design as simple and as barebones as possible. I tried staying as close to the golang standard library as possible. 

Golang has a pretty powerful and production-ready http client in the standary library (`net/http`) as well as a great mock server (`net/httptest`) that can be leveraged to great effect. This allows us NOT to rely on other popular golang restful API testing patterns/factories such as `gomock` or `httpmock`.

Security concerns in a production environment (such as TLS configuration) were ignored in this project as I believe it is out-of-scope for the purposes of this project.

## About the Client Implementation
The `net/http` client offers alot of extensibilty, and my client implementation in [client.go](src/client/client.go) primarily focuses around two areas of customization; `timeouts` and `transports`. More information about timeouts can found [here](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/). The default `http.Client` is instantiated and initialized as follows:
```go
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
```
The `Timeout` and `net.Dialier.Timeout` values and their respective effects are essentially the same, but the code has been written out this way to show a consumer that these parameters are customizeable. A user can just set or consume the default `HTTPClient.Timeout` setting, or if they need more fine-grained control over this behavoir, they can define their own `http.Transport` object and bind it to the client.
## About the Accounts API Implementation
The implementation is as simple as can be, and leverages the defaults set in the [client package](src/client/client.go) to make http calls. There is no explict control of client parameters in the [accounts package](src/accounts/api.go), but any consumer of this code can add them.
## Deploy
From the root of the repository, run:
```bash
docker-compose up
```
##### Note - you might have to clean out your docker image cache to force of a re-pull/re-build of the containers listed in [docker-compose.yml](docker-compose.yml)
## Run the Example File
<b>Make sure the backend API is reachable before running the example file</b>. From the root of the repository, run:
```bash
FORM3_ACCOUNTS_API_URL="http://localhost:8080/v1/organisation/accounts" go run src/example.go
```
### Output from API Server
```bash
accountapi_1              | [GIN] 2022/02/13 - 23:02:14 | 201 |    5.275821ms |      172.21.0.1 | POST     "/v1/organisation/accounts"
accountapi_1              | [GIN] 2022/02/13 - 23:02:14 | 200 |     543.171µs |      172.21.0.1 | GET      "/v1/organisation/accounts/bc5c052d-c486-478b-8dd2-afe82fd7725d"
accountapi_1              | [GIN] 2022/02/13 - 23:02:14 | 204 |    1.242251ms |      172.21.0.1 | DELETE   "/v1/organisation/accounts/bc5c052d-c486-478b-8dd2-afe82fd7725d?version=0"
```

## Run Tests (All Inclusive)
From the root of the repository, run:
```bash
FORM3_ACCOUNTS_API_URL="http://localhost:8080/v1/organisation/accounts" go test -v ./...
```
## Generate Testing Coverage Report
A coverage report file is generated from the docker-compose bootstrap, entitled `coverage_report_from_container.out`. If that file is not present in the root of this repositry after executing `docker-compose up`, you can run the following command:
```bash
FORM3_ACCOUNTS_API_URL="http://localhost:8080/v1/organisation/accounts" go test -v ./... -coverprofile=report.out && go tool cover -html=report.out
```
## Wait? Why Isn't Code Coverage 100% for `api.go`
The coverage report shows that I am not testing `json.Marshal()` calls. This is by design. I didn't see the value in creating mocks and extra code for the sake of code coverage. It is my opinion that if data is properly marshalled to `json.Marshal()`, that this error will not fire. A potential situation where the error would fire could be a programming error.
## Code Usage
A user can consume the following packages/imports into their own program, to create/fetch/delete Form3 accounts to a Form3 backend API.
```go
import (
  "github.com/sarabrajsingh/interview-accountapi/src/models"
  "github.com/sarabrajsingh/interview-accountapi/src/accounts"
)
```
After importing the libraries, a consumer must create a valid account object before sending this payload off to the backend API via the accounts package. 

### Example Account Struct Creation
```go
// some default and required variables
AccountClassification := "Personal"
AccountMatchingOptOut := false
Country := "GB"
JointAccount := false

account_id := uuid.New().String()
account_organisation_id := uuid.New().String()
account_version := 0

accountAttribs := models.AccountAttributes{
  Country:      &Country,
  BaseCurrency: "GBP",
  BankID:       "400300",
  BankIDCode:   "GBDSC",
  Bic:          "NWBKGB22",
  Name: []string{
    "Samantha Holder",
  },
  AlternativeNames: []string{
    "Sam Holder",
  },
  AccountClassification:   &AccountClassification,
  JointAccount:            &JointAccount,
  AccountMatchingOptOut:   &AccountMatchingOptOut,
  SecondaryIdentification: "A1B2C3D4",
}

accountData := models.AccountData{
  Attributes:     &accountAttribs,
  ID:             account_id,
  OrganisationID: account_organisation_id,
  Type:           "accounts",
}

account := models.Account{
  Data: &accountData,
}
```
Then the `create`/`fetch`/`delete` methods can be leveraged like so:
### CREATE
```go
resp, err := accounts.Create(account)
if err != nil {
  fmt.Println(err)
}
fmt.Println(resp)
```
### CREATE with a Context
```go
ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10)
resp, err := CreateWithCtx(ctx, generateAccount())
if err != nil {
  fmt.Println(err)
}
fmt.Println(resp)
```
### FETCH
```go
resp, err = accounts.Fetch(account_id)
if err != nil {
  fmt.Println(err)
}
fmt.Println(resp)
```
### DELETE
```go
resp, err = accounts.Delete(account_id, account_version)
if err != nil {
  fmt.Println(err)
}
fmt.Println(resp)
```
Please refer to [example.go](src/example.go) for the full source code to an example file. The `create/fetch/delete` functions can also be used with `Context` objects.
## Project Structure
```bash
.
├── docker-compose.yml
├── Dockerfile
//
└── src
    ├── accounts
    │   ├── api.go
    │   └── api_test.go
    ├── client
    │   ├── client.go
    │   └── client_test.go
    ├── models
    │   └── models.go
    └── example.go
```
### `docker-compose.yml`
The main `docker-compose` file that bootstraps the backend API using the docker engine. This provides a mock form3 api backend for consumption. Unit and integration tests for the client and accounts api have been baked into a container, and included into the docker-compose file. They (the testing container) will run after the backend API and database containers have initialized and bootstrapped.

### `Dockerfile`
This dockerfile builds a container environment for this projects' unit tests and integration tests. A rootless container design was choosen for this container definiton as it most closely reflects what a given user might witness in a production environment.

### `src/accounts`
`api.go` contains the wrapper methods that tie together the `models/models.go` account models with the http client module. This is an example implementation, and doesn't necessarily leverage all the tweaks to a http client that a production grade system requires.

`api_test.go` contains unit and integration tests for `api.go`

### `src/client`
`client.go` contains the main machinery and wrapper code around `net/http` to implement a http client. This module is the main crux of this demo, and was designed to be simple to consume. Any optional wrapper structs and modules around Error messages for example, was omitted by design, and instead, left to the consumers on how to handle errors recieved from a restful service.

`client_test.go` contains the unit and functional tests for the client package

### `src/models`
`models.go` contains the structs that represent accounts for the backend API

### `src/example.go`
`example.go` contains the source code to an example resource, which shows how a consumer might leverage the client and accounts pakcage(s) to write account objects to the backend API.

### Improvments for Ambassador Labs Interview
1. Decoupled integration tests from database, such that every integration test will restore the database to it's initial seed point
2. Refactored unit and integration tests by introducing testing tables
3. Unmarhalling client responses from a flat string into an Account struct object
4. Introduced new errors package to highlight and kee

## References
1. [SendGrid REST API](https://github.com/sendgrid/rest)
2. [mocking outbound http requests in go: you’re (probably) doing it wrong](https://medium.com/zus-health/mocking-outbound-http-requests-in-go-youre-probably-doing-it-wrong-60373a38d2aa)
3. [Resty](https://github.com/go-resty/resty)
4. [Official Golang Documentation](https://pkg.go.dev/net/http)
5. [Tuning the Golang net/http client](https://www.loginradius.com/blog/async/tune-the-go-http-client-for-high-performance/)
6. [net/http Transport documentation](https://go.dev/src/net/http/transport.go)
7. [More net/http client timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)