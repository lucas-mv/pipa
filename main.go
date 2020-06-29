package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/lucas-mv/pipa/external"
	"github.com/lucas-mv/pipa/utils"
)

func main() {
	fmt.Println("Welcome to pipa! 🐶")
	fmt.Println("I'll be your guide🦮 to twitter🐦! Let me fetch your most relevant trends...")

	trends := getTrends()

	fmt.Println("All done, here are your top trends!")
	fmt.Println("---------------------------------------------------------------------------")

	for i := 0; i < len(trends); i++ {
		printNamedTrends(trends[i])
	}

	fmt.Println("That's all for now! Come back later for more relevant trends! 🐕")
}

func printNamedTrends(namedTrends namedTrends) {
	fmt.Println(namedTrends.name)
	fmt.Println()
	for i := 0; i < 5; i++ {
		fmt.Println("\t#" + strconv.Itoa(i+1))
		printTrendingTopic(namedTrends.trends[i])
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("---------------------------------------------------------------------------")
}

func printTrendingTopic(topic external.TrendingTopic) {
	fmt.Println("\tName: " + topic.Name)
	fmt.Println("\tURL: " + topic.URL)
	fmt.Println("\tTweet Volume: " + strconv.FormatInt(topic.TweetVolume, 10))
	fmt.Println("\tPromoted content: " + strconv.FormatBool(topic.PromotedContent != ""))
}

func getTrends() []namedTrends {
	client := &http.Client{}
	settings := utils.GetSettings()
	twitterAuthentication := external.GetAccessToken(client, settings.TwitterBasicKey)

	globalTrends := make(chan namedTrends)
	go func() {
		globalTrends <- namedTrends{3, "🛰️  Global", external.GetTrendingTopics(client, 1, twitterAuthentication.AccessToken)}
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
}

func getLocationTrends(client *http.Client, settings utils.PipaSettings, accessToken string) []namedTrends {
	lat, long := external.GetLatLong(settings.Address, settings.BingAPIKey)

	WOEID, ParentWOEID := external.GetWOEID(client, lat, long, accessToken)

	localTrends := make(chan namedTrends)
	go func() {
		localTrends <- namedTrends{1, "🛵  Local", external.GetTrendingTopics(client, WOEID, accessToken)}
	}()

	regionalTrends := make(chan namedTrends)
	go func() {
		regionalTrends <- namedTrends{2, "🚌  Regional", external.GetTrendingTopics(client, ParentWOEID, accessToken)}
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
