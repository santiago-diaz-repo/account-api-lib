package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	country := "GB"
	accountClassification := "Personal"
	jointAccount := false
	accountMatchingOptOut := false
	data := Data{
		Data: AccountData{
			ID: "ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6",
			OrganisationID: "ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6",
			Type: "accounts",
			Attributes: &AccountAttributes{
				Country: &country,
				BaseCurrency: "GBP",
				BankID: "400302",
				BankIDCode: "GBDSC",
				Bic: "NWBKGB42",
				Name: []string{"Santiago"},
				AlternativeNames: []string{"test"},
				AccountClassification: &accountClassification,
				JointAccount: &jointAccount,
				AccountMatchingOptOut: &accountMatchingOptOut,
				SecondaryIdentification: "A1B2C3D4",
			},
		},
	}

	re, err := json.Marshal(data)
	if err != nil{
		log.Fatalln(err)
	}

	reader := strings.NewReader(string(re))

	request, err := http.NewRequest(http.MethodPost,"http://localhost:8090/v1/organisation/accounts",reader)
	if err != nil{
		log.Fatalln(err)
	}

	t := http.DefaultTransport.(*http.Transport).Clone()

	client := &http.Client{
		Transport: t,
	}

	response, err := client.Do(request)
	if err != nil{
		log.Fatalln(err)
	}

	if response.StatusCode == 201{
		fmt.Printf("%s", response.Body)
	}else {
		fmt.Printf("%d",response.StatusCode)
	}
}
