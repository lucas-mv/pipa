package main

import (
	"fmt"
	"net/http"
)

func main() {
	trends := GetTrends()
	fmt.Println(PrettyPrint(trends))
}

func GetTrends() []TrendingTopic {
	settings := GetSettings()

	lat, long := GetLatLong(settings.Address, settings.BingAPIKey)

	client := &http.Client{}

	twitterAuthentication := GetAccessToken(client, settings.TwitterBasicKey)

	WOEID := GetWOEID(client, lat, long, twitterAuthentication.AccessToken)

	return GetTrendingTopics(client, WOEID, twitterAuthentication.AccessToken)
}
