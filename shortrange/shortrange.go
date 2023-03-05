package shortrange

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/umahmood/haversine"
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

// create a function to filter through the list of charger structs and return a list of chargers that are within a certain distance of the zipcode
func FilterChargersByDistance(chargers []Charger, zipcode string, homelat string, homelong string) []Charger {
	var filteredChargers []Charger
	for _, charger := range chargers {
		if charger.AddressInfo.Postcode == zipcode {
			filteredChargers = append(filteredChargers, charger)
		}
	}

	return NearestFromReference(filteredChargers, homelat, homelong)
}

// sort the nearest chargers by charger.AddressInfo.Distance
func SortByDistance(chargers []Charger) []Charger {
	sort.Slice(chargers, func(i, j int) bool {
		return chargers[i].AddressInfo.Distance < chargers[j].AddressInfo.Distance
	})
	return chargers
}

func NearestFromReference(chargers []Charger, homelat string, homelong string) []Charger {

	var nearestChargers []Charger
	for _, charger := range chargers {

		//get the latitude and longitude of the charger
		chargerLat := charger.AddressInfo.Latitude
		chargerLong := charger.AddressInfo.Longitude

		//get the latitude and longitude of the reference point
		//convert the string to a float

		referenceLat, err1 := strconv.ParseFloat(homelat, 64)
		if err1 != nil {
			fmt.Println(err1)
		}

		referenceLong, err2 := strconv.ParseFloat(homelong, 64)
		if err2 != nil {
			fmt.Println(err2)
		}

		// Define the two points you want to measure the distance between.
		point1 := haversine.Coord{Lat: chargerLat, Lon: chargerLong}
		point2 := haversine.Coord{Lat: referenceLat, Lon: referenceLong}

		// Calculate the distance between the two points.
		distance, _ := haversine.Distance(point1, point2)

		charger.AddressInfo.Distance = int(distance)
		// Print the distance between the two points.
		fmt.Println(distance)
		nearestChargers = append(nearestChargers, charger)

	}
	return nearestChargers

}
