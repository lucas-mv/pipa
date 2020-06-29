package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/lucas-mv/pipa/external"
	"github.com/lucas-mv/pipa/utils"
)

func main() {
	trends := getTrends()
	for i := 0; i < len(trends); i++ {
		printNamedTrends(trends[i])
	}
}

func printNamedTrends(namedTrends namedTrends) {
	fmt.Println(namedTrends.name)
	fmt.Println(utils.PrettyPrint(namedTrends.trends[0:5]))
	fmt.Println("------------------------------------------------------------------------------------------")
}

func getTrends() []namedTrends {
	client := &http.Client{}
	settings := utils.GetSettings()
	twitterAuthentication := external.GetAccessToken(client, settings.TwitterBasicKey)

	globalTrends := make(chan namedTrends)
	go func() {
		globalTrends <- namedTrends{3, "Global", external.GetTrendingTopics(client, 1, twitterAuthentication.AccessToken)}
	}()

	localRegionalTrends := getLocationTrends(client, settings, twitterAuthentication.AccessToken)

	var trends []namedTrends
	trends = append(trends, localRegionalTrends...)
	trends = append(trends, <-globalTrends)

	orderNamedTrends(trends)

	return trends
}

func orderNamedTrends(trends []namedTrends) {
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].relevance < trends[j].relevance
	})

	for k := 0; k < len(trends); k++ {
		sort.Slice(trends[k].trends, func(i, j int) bool {
			return trends[k].trends[i].TweetVolume > trends[k].trends[j].TweetVolume
		})
	}
}

func getLocationTrends(client *http.Client, settings utils.PipaSettings, accessToken string) []namedTrends {
	lat, long := external.GetLatLong(settings.Address, settings.BingAPIKey)

	WOEID, ParentWOEID := external.GetWOEID(client, lat, long, accessToken)

	localTrends := make(chan namedTrends)
	go func() {
		localTrends <- namedTrends{1, "Local", external.GetTrendingTopics(client, WOEID, accessToken)}
	}()

	regionalTrends := make(chan namedTrends)
	go func() {
		regionalTrends <- namedTrends{2, "Regional", external.GetTrendingTopics(client, ParentWOEID, accessToken)}
	}()

	var trends []namedTrends
	trends = append(trends, <-localTrends)
	trends = append(trends, <-regionalTrends)

	return trends
}

type namedTrends struct {
	relevance int
	name      string
	trends    []external.TrendingTopic
}
