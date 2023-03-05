package main

//import shortrange package

import (
	"autobots/shortrange"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetChargersByZip(zipcode string, apikey string, homelat string, homelong string) ([]shortrange.Charger, error) {
	//set url with zipcode and apikey
	url := fmt.Sprintf("https://api.openchargemap.io/v3/poi/?output=json&countrycode=US&maxresults=1000000&postalcode=%s&compact=true&verbose=false&key=%s", zipcode, apikey)
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
	var chargers []shortrange.Charger
	if err := json.Unmarshal(body, &chargers); err != nil {
		return nil, err
	}

	//return filtered chargers
	return shortrange.FilterChargersByDistance(chargers, zipcode, homelat, homelong), nil

}

// create a webserver that will listen for a zipcode and return the chargers in that zipcode in html
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//get zipcode from url
		zipcode := r.URL.Query().Get("zipcode")
		//printf zipcode
		fmt.Println(zipcode)

		//get apikey from url
		apikey := r.URL.Query().Get("apikey")

		//get latitude from url
		homelat := r.URL.Query().Get("latitude")

		//get longitude from url
		homelong := r.URL.Query().Get("longitude")

		//get chargers from openchargemap
		chargers, err := GetChargersByZip(zipcode, apikey, homelat, homelong)
		//sort chargers by distance
		chargers = shortrange.SortByDistance(chargers)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
			//fmt.Fprintf(w, "<p>%d</p>", charger.AddressInfo.DistanceUnit)
			fmt.Fprintf(w, "<p>%f</p>", charger.AddressInfo.Latitude)
			fmt.Fprintf(w, "<p>%f</p>", charger.AddressInfo.Longitude)
			//fmt.Fprintf(w, "<p>%d</p>", charger.AddressInfo.ID)
			//fmt.Fprintf(w, "<p>%d</p>", charger.ID)
			fmt.Fprintf(w, "<hr>")

		}

	})
	http.ListenAndServe(":8080", nil)
}
