package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const botConfigFile = "bot_config.json"

var instance *BotConfig
var once sync.Once

func (botConfig *BotConfig) SaveToFile() {
	log.Print("Trying to save bot config to file...")

	file, err := json.Marshal(botConfig)

	if err != nil {
		log.Print("Couldn't marshall bot config")
		return
	}

	err = ioutil.WriteFile(botConfigFile, file, os.ModePerm)

	if err != nil {
		log.Print("Couldn't save bot config to file")
		return
	}

	log.Print("Bot config saved to file")
}

func (storage *BotConfig) ReadFromFile() (BotConfig, error) {
	log.Print("Trying to read bot confi from file...")

	file, _ := ioutil.ReadFile(botConfigFile)

	data := BotConfig{}

	err := json.Unmarshal([]byte(file), &data)

	if err != nil {
		log.Print("Bot config not found in file")
	}

	return data, err
}

func NewBotConfig() *BotConfig {
	once.Do(func() {

		instance = &BotConfig{}

		data, err := instance.ReadFromFile()

		if err == nil {
			instance = &data
		}

	})

	return instance
}
