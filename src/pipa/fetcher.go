package pipa

import (
	"net/http"
)

func GetTrends() {
	settings := GetSettings()

	lat, long := GetLatLong(settings.Address, settings.BingAPIKey)

	client := &http.Client{}

	twitterAuthentication := GetAccessToken(client, settings.TwitterBasicKey)

	WOEID := GetWOEID(client, lat, long, twitterAuthentication.AccessToken)

	trendLocations := GetTrendingTopics(client, WOEID, twitterAuthentication.AccessToken)
}
