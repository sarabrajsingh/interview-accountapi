package main

import (
	"fmt"
	"github.com/sarabrajsingh/interview-accountapi/src/accounts"
	"github.com/sarabrajsingh/interview-accountapi/src/models"
	"github.com/google/uuid"
)

func main() {

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

	resp, err := accounts.Create(account)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

	resp, err = accounts.Fetch(account_id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

	resp, err = accounts.Delete(account_id, account_version)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
