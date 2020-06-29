package main

import (
	"fmt"
	"net/http"

	"github.com/lucas-mv/pipa/external"
	"github.com/lucas-mv/pipa/utils"
)

func main() {
	trends := getTrends()
	fmt.Println(utils.PrettyPrint(trends))
}

func getTrends() []external.TrendingTopic {
	settings := utils.GetSettings()

	lat, long := external.GetLatLong(settings.Address, settings.BingAPIKey)

	client := &http.Client{}

	twitterAuthentication := external.GetAccessToken(client, settings.TwitterBasicKey)

	WOEID := external.GetWOEID(client, lat, long, twitterAuthentication.AccessToken)

	return external.GetTrendingTopics(client, WOEID, twitterAuthentication.AccessToken)
}
