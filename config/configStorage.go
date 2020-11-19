package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

const configFile = "config.json"

var instance *AppConfig
var once sync.Once

type AppConfig struct {
	WebhookUrl  string `json:"webhook_url"`
	BotSecret   string `json:"bot_secret"`
	AppClientId string `json:"app_client_id"`
	AppSecretId string `json:"app_secret_id"`
	RedirectURI string `json:"redirect_uri"`
}

func (storage *AppConfig) ReadFromFile() (AppConfig, error) {
	log.Print("Trying to read app config from file...")

	file, _ := ioutil.ReadFile(configFile)

	data := AppConfig{}

	err := json.Unmarshal([]byte(file), &data)

	if err != nil {
		log.Fatal("App config not found in file")
	}

	return data, err
}

func NewAppConfig() *AppConfig {
	once.Do(func() {

		instance = &AppConfig{}

		data, err := instance.ReadFromFile()

		if err == nil {
			instance = &data
		}

	})

	return instance
}
