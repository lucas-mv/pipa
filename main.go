package main

import (
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/lucas-mv/pipa/external"
	"github.com/lucas-mv/pipa/utils"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", run)
	http.ListenAndServe(":"+port, nil)
}

func run(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	address := query.Get("address")

	trends := getTrends(address)

	io.WriteString(w, "Welcome to pipa! 🐶\nI'll be your guide🦮 to twitter🐦! Here are your most relevant trends...")
	io.WriteString(w, "\n---------------------------------------------------------------------------")

	for i := 0; i < len(trends); i++ {
		printNamedTrends(w, trends[i])
	}

	io.WriteString(w, "\nThat's all for now! Come back later for more relevant trends! 🐕")
}

func printNamedTrends(w http.ResponseWriter, namedTrends namedTrends) {
	io.WriteString(w, "\n"+namedTrends.name+"\n")
	for i := 0; i < 5; i++ {
		io.WriteString(w, "\t#"+strconv.Itoa(i+1)+"\n")
		printTrendingTopic(w, namedTrends.trends[i])
		io.WriteString(w, "\n")
	}
	io.WriteString(w, "\n---------------------------------------------------------------------------")
}

func printTrendingTopic(w http.ResponseWriter, topic external.TrendingTopic) {
	io.WriteString(w, "\tName: "+topic.Name+"\n")
	io.WriteString(w, "\tURL: "+topic.URL+"\n")
	io.WriteString(w, "\tTweet Volume: "+strconv.FormatInt(topic.TweetVolume, 10)+"\n")
	io.WriteString(w, "\tPromoted content: "+strconv.FormatBool(topic.PromotedContent != "")+"\n")
}

func getTrends(address string) []namedTrends {
	client := &http.Client{}
	settings := utils.GetSettings()
	twitterAuthentication := external.GetAccessToken(client, settings.TwitterBasicKey)

	globalTrends := make(chan namedTrends)
	go func() {
		globalTrends <- namedTrends{3, "🛰️  Global", external.GetTrendingTopics(client, 1, twitterAuthentication.AccessToken)}
	}()

	localRegionalTrends := getLocalAndRegionalTrends(client, settings, twitterAuthentication.AccessToken, address)

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

func getLocalAndRegionalTrends(client *http.Client, settings utils.PipaSettings, accessToken, address string) []namedTrends {
	lat, long := external.GetLatLong(address, settings.BingAPIKey)

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
