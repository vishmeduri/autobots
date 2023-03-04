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

// create a webservver that will listen for a zipcode and return the chargers in that zipcode in html
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//get zipcode from url
		zipcode := r.URL.Query().Get("zipcode")
		//get apikey from url
		apikey := r.URL.Query().Get("apikey")
		//get chargers from openchargemap
		chargers, err := GetChargersByZip(zipcode, apikey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//print chargers
		fmt.Println(chargers)
		//print chargers in html
		fmt.Fprintf(w, "%v", chargers)

		//set html header in response
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//print chargers in html

		fmt.Fprintf(w, "<h1>Number of chargers: %d</h1>", len(chargers))
		for _, charger := range chargers {
			fmt.Fprintf(w, "<h2>%s</h2>", charger.AddressInfo.Title)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.AddressLine1)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.AddressLine2)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.Town)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.StateOrProvince)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.Postcode)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.Country.Title)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.ContactTelephone1)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.ContactTelephone2)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.ContactEmail)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.AccessComments)
			fmt.Fprintf(w, "<p>%s</p>", charger.AddressInfo.RelatedURL)
			fmt.Fprintf(w, "<p>%d</p>", charger.AddressInfo.Distance)
			fmt.Fprintf(w, "<p>%d</p>", charger.AddressInfo.DistanceUnit)
		}

	})
	http.ListenAndServe(":8080", nil)
}
