package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type bingAddress struct {
	StatusCode   int            `json:"statusCode"`
	ResourceSets []resourceSets `json:"resourceSets"`
}

type resourceSets struct {
	EstimatedTotal int        `json:"estimatedTotal"`
	Resources      []resource `json:"resources"`
}

type resource struct {
	Name  string   `json:"name"`
	Point geoPoint `json:"point"`
}

type geoPoint struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func GetLatLong(address, bingAPIKey string) (lat, long float64) {
	query := url.Values{
		"q":   {address},
		"o":   {"json"},
		"key": {bingAPIKey},
	}

	res, err := http.Get("http://dev.virtualearth.net/REST/v1/Locations?" + query.Encode())
	if err != nil {
		fmt.Print(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var bingAddress bingAddress
	err = json.Unmarshal(body, &bingAddress)
	if err != nil {
		log.Fatal(err)
	}

	if bingAddress.StatusCode != 200 {
		log.Fatal("Invalid status code")
	}

	coordinates := bingAddress.ResourceSets[0].Resources[0].Point.Coordinates

	return coordinates[0], coordinates[1]
}
