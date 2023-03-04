package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Charger struct {
	ID          int         `json:"ID"`
	AddressInfo AddressInfo `json:"AddressInfo"`
}

type AddressInfo struct {
	ID                int     `json:"ID"`
	AddressLine1      string  `json:"AddressLine1"`
	AddressLine2      string  `json:"AddressLine2"`
	Town              string  `json:"Town"`
	StateOrProvince   string  `json:"StateOrProvince"`
	Postcode          string  `json:"Postcode"`
	CountryID         int     `json:"CountryID"`
	Country           Country `json:"Country"`
	Latitude          float64 `json:"Latitude"`
	Longitude         float64 `json:"Longitude"`
	ContactTelephone1 string  `json:"ContactTelephone1"`
	ContactTelephone2 string  `json:"ContactTelephone2"`
	ContactEmail      string  `json:"ContactEmail"`
	AccessComments    string  `json:"AccessComments"`
	RelatedURL        string  `json:"RelatedURL"`
	Distance          int     `json:"Distance"`
	DistanceUnit      int     `json:"DistanceUnit"`
	Title             string  `json:"Title"`
}

type Country struct {
	ID            int    `json:"ID"`
	ISOCode       string `json:"ISOCode"`
	ContinentCode string `json:"ContinentCode"`
	Title         string `json:"Title"`
}

func GetChargersByZip(zipcode string, apikey string) ([]Charger, error) {
	//set url with zipcode and apikey
	url := fmt.Sprintf("https://api.openchargemap.io/v3/poi/?output=json&countrycode=US&maxresults=10&postalcode=%s&compact=true&verbose=false&key=%s", zipcode, apikey)
	//printf url
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var chargers []Charger
	if err := json.Unmarshal(body, &chargers); err != nil {
		return nil, err
	}
	return chargers, nil
}

// creart main function with zipcode as argument and prompt user for input

func main() {

	//prompt user for zipcode
	fmt.Println("Enter a zipcode: ")
	var zipcode string
	fmt.Scanln(&zipcode)

	fmt.Println("Enter apikey: ")
	var apikey string
	fmt.Scanln(&apikey)

	chargers, err := GetChargersByZip(zipcode, apikey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(chargers)
	fmt.Println("Number of chargers: ", len(chargers))

}
