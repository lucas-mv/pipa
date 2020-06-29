package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type pipaSettings struct {
	BingAPIKey      string `json:"bing_api_key"`
	TwitterBasicKey string `json:"twitter_basic_key"`
	Address         string `json:"address"`
}

func GetSettings() pipaSettings {
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
