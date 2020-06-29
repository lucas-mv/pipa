package utils

import (
	"os"
)

//PipaSettings structure of the project settings settings.json file
type PipaSettings struct {
	BingAPIKey      string
	TwitterBasicKey string
}

//GetSettings returns the settings configured in settings.json
func GetSettings() PipaSettings {
	return PipaSettings{
		os.Getenv("bing_api_key"),
		os.Getenv("twitter_basic_key"),
	}
}
