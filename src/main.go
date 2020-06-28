package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func main() {
	settings := getSettings()

	lat, long := getLatLong(settings.Address, settings.BingAPIKey)

	client := &http.Client{}

	twitterAuthentication := getAccessToken(client, settings.TwitterBasicKey)

	WOEID := getWOEID(client, lat, long, twitterAuthentication.AccessToken)

	trendLocations := getTrendingTopics(client, WOEID, twitterAuthentication.AccessToken)

	fmt.Println(prettyPrint(trendLocations))
}

func getSettings() pipaSettings {
	data, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		fmt.Print(err)
	}

	var settings pipaSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		log.Fatal(err)
	}

	return settings
}

func getLatLong(address, bingAPIKey string) (lat, long float64) {
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

func getWOEID(client *http.Client, lat, long float64, accessToken string) int {
	trendsURL := "https://api.twitter.com/1.1/trends/closest.json"
	req, err := http.NewRequest("GET", trendsURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	query := url.Values{
		"lat":  {floatToString(lat)},
		"long": {floatToString(long)},
	}
	req.URL.RawQuery = query.Encode()

	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var twitterPlaces []twitterPlace
	err = json.Unmarshal(body, &twitterPlaces)
	if err != nil {
		log.Fatal(err)
	}

	return twitterPlaces[0].WOEID
}

func getTrendingTopics(client *http.Client, location int, accessToken string) []trendingTopic {
	trendsURL := "https://api.twitter.com/1.1/trends/place.json"
	req, err := http.NewRequest("GET", trendsURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	query := url.Values{
		"id": {strconv.Itoa(location)},
	}
	req.URL.RawQuery = query.Encode()

	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var trendLocations []trendLocation
	err = json.Unmarshal(body, &trendLocations)
	if err != nil {
		log.Fatal(err)
	}

	return trendLocations[0].Trends
}

func getAccessToken(client *http.Client, twiterBasicKey string) twitterAuthentication {
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")

	url := "https://api.twitter.com/oauth2/token"

	tokenRequest, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	tokenRequest.Header.Add("Authorization", "Basic "+twiterBasicKey)
	tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	tokenRequest.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

	tokenResponse, err := client.Do(tokenRequest)
	if err != nil {
		log.Fatal(err)
	}
	defer tokenResponse.Body.Close()

	body, err := ioutil.ReadAll(tokenResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	var twitterAuthentication twitterAuthentication
	json.Unmarshal(body, &twitterAuthentication)

	return twitterAuthentication
}

func floatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 6, 64)
}

func prettyPrint(i interface{}) string {
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	return string(s)
}

type pipaSettings struct {
	BingAPIKey      string `json:"bing_api_key"`
	TwitterBasicKey string `json:"twitter_basic_key"`
	Address         string `json:"address"`
}

type twitterAuthentication struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type trendLocation struct {
	Trends    []trendingTopic `json:"trends"`
	AsOf      time.Time       `json:"as_of"`
	CreatedOn time.Time       `json:"created_on"`
	Locations []twitterPlace  `json:"locations"`
}

type trendingTopic struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	PromotedContent string `json:"promoted_content"`
	TweetVolume     int64  `json:"tweet_volume"`
}

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

type twitterPlace struct {
	Name  string `json:"name"`
	WOEID int    `json:"woeid"`
}
