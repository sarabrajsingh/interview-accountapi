package accounts

import (
	. "github.com/sarabrajsingh/interview-accountapi/src/models"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
	"golang.org/x/net/context"
	"os"
)

const (
	account_id              = "f773707e-769e-4ed6-9194-ab69ff639d39"
	account_organisation_id = "4fd712d9-e281-4add-8d66-800f6960b57c"
	account_version = 0
)

// random number generator; useful for generating near-unique account numbers in testing
func randomAccountNumberGenerator(n int) string {
	var accountNumberRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = accountNumberRunes[rand.Intn(len(accountNumberRunes))]
	}

	return string(b)
}

// helper function to generate working accounts
func generateAccount() Account {

	rand.Seed(time.Now().UnixNano())

	AccountClassification := "Personal"
	AccountMatchingOptOut := false
	Country := "GB"
	JointAccount := false

	return Account{
		Data: &AccountData{
			Attributes: &AccountAttributes{
				Country:      &Country,
				BaseCurrency: "GBP",
				BankID:       randomAccountNumberGenerator(6),
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
			},
			ID:             account_id,
			OrganisationID: account_organisation_id,
			Type:           "accounts",
		},
	}
}

// Unittest-1 - Make sure the SetBaseURL() function on our custom struct URL type, works as expected.
func TestSetBaseURL(t *testing.T) {
	DefaultUrl.SetBaseURL("super.fake.com")
	assert.Equal(t, DefaultUrl.BaseURL, "super.fake.com", "setting custom BaseURL failed")
}

/* INTEGRATION TESTS START HERE
  NOTES:
- Take caution when modifying these tests. TestSetBaseURL change the Accounts.api.URL.BaseURL around
  and this variable is leveraged in the Create(), Fetch() and Delete() methods that are also in that package. If they point to the wrong
  baseURL, the following integration tests will fail.
- Preserve the ordering of the Integration tests as well, as they will break if the proper account(s) exist or don't exist in the API backend
- The first integration test sets the BaseURL to the backend api for the entire api.go package so this needs to be set correctly
*/

// IntegrationTest-1 - Create a canonical Accounts account against the API. Expect success 201
func TestCreateValidAccount(t *testing.T) {
	DefaultUrl.SetBaseURL(os.Getenv("FORM3_ACCOUNTS_API_URL"))
	resp, err := Create(generateAccount())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 201, "failed to create new account against API")
	assert.NotNil(t, resp.Body)
}

// IntegreationTest-2 - Try and recreate the same account from IntegrationTest-1. Except failure 409
func TestCreateValidButDuplicatedAccount(t *testing.T) {
	resp, err := Create(generateAccount())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 409, "this test should have failed, as we are trying to duplicate an account_id against the API")
	assert.NotNil(t, resp.Body)
}

// IntegrationTest-3 - create an invalid account against the backend api
func TestCreateInvalidAccount(t *testing.T) {
	acc := generateAccount()
	acc.Data.ID = "abc123"
	resp, err := Create(acc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 400, "failed to create bad account against API")
	assert.NotNil(t, resp.Body)
}

// IntegrationTest-4 - fetch a valid account from the backend API
func TestFetchValidAccount(t *testing.T) {
	resp, err := Fetch(account_id)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 200, "failed to GET resource from API backend")
}

// IntegrationTest-5 - try to fetch an invalid account from the backend API
func TestFetchInvalidAccount(t *testing.T) {
	resp, err := Fetch("superfake.com")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 400, "failed to grab an invalid account")
}

// IntegrationTest-6 - delete a valid account from the backend API
func TestDeleteValidAccount(t *testing.T) {
	resp, err := Delete(account_id, account_version)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 204, "failed to delete resource from API backend")
}

// IntegrationTest7 - create a valid account with a custom context
func TestCreateValidAccountWithCtx(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10)
	resp, err := CreateWithCtx(ctx, generateAccount())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 201, "failed to create new account against API")
	assert.NotNil(t, resp.Body)
}

// IntegrationTest8 - fetch a valid account with a custom context
func TestFetchValidAccountWithCtx(t *testing.T){
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10)
	resp, err := FetchWithCtx(ctx, account_id)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 200, "failed to GET resource from API backend")
}

// IntegrationTest9 - dekete a valid account with a custom context
func TestDeleteValidAccountWithCtx(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10)
	resp, err := DeleteWithCtx(ctx, account_id, account_version)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 204, "failed to delete resource from API backend")	
}

// IntegrationTest10 - delete an ivalid account
func TestDeleteInvalidAccount(t *testing.T) {
	resp, err := Delete("superfake.com", 0)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.StatusCode, 400, "failed to delete an invalid account")
}

// IntegrationTest11 - Test a bad context when trying to create an account
func TestCreateValidAccountWithBadCtx(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond*1)
	resp, err := CreateWithCtx(ctx, generateAccount())
	if err == nil{
		t.Error("this error should have fired")
	}
	assert.Nil(t, resp)
}

// IntegrationTest12 - Test a bad context when trying to fetch an account
func TestFetchValidAccountWithBadCtx(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond*1)
	resp, err := FetchWithCtx(ctx, account_id)
	if err == nil{
		t.Error("this error should have fired")
	}
	assert.Nil(t, resp)	
}

// IntegrationTest13 - Test a bad context when trying to delete an account
func TestDeleteValidAccountWithBadCtx(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond*1)
	resp, err := DeleteWithCtx(ctx, account_id, account_version)
	if err == nil{
		t.Error("this error should have fired")
	}
	assert.Nil(t, resp)	
}