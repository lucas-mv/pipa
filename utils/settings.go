package utils

import (
	"os"
)

//PipaSettings structure of the project settings settings.json file
type PipaSettings struct {
	BingAPIKey      string `json:"bing_api_key"`
	TwitterBasicKey string `json:"twitter_basic_key"`
	Address         string `json:"address"`
}

//GetSettings returns the settings configured in settings.json
func GetSettings() PipaSettings {
	return PipaSettings{
		os.Getenv("bing_api_key"),
		os.Getenv("twitter_basic_key"),
		os.Getenv("address"),
	}
}
