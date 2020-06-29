package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

//PipaSettings structure of the project settings settings.json file
type PipaSettings struct {
	BingAPIKey      string `json:"bing_api_key"`
	TwitterBasicKey string `json:"twitter_basic_key"`
	Address         string `json:"address"`
}

//GetSettings returns the settings configured in settings.json
func GetSettings() PipaSettings {
	data, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		fmt.Print(err)
	}

	var settings PipaSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		log.Fatal(err)
	}

	return settings
}
