package external

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetWOEID(client *http.Client, lat, long float64, accessToken string) int {
	trendsURL := "https://api.twitter.com/1.1/trends/closest.json"
	req, err := http.NewRequest("GET", trendsURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	query := url.Values{
		"lat":  {FloatToString(lat)},
		"long": {FloatToString(long)},
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

func GetTrendingTopics(client *http.Client, location int, accessToken string) []trendingTopic {
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

func GetAccessToken(client *http.Client, twiterBasicKey string) twitterAuthentication {
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

type twitterPlace struct {
	Name  string `json:"name"`
	WOEID int    `json:"woeid"`
}
