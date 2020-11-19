package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type LocalTokenStorage struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

var instance *LocalTokenStorage
var once sync.Once

func (storage LocalTokenStorage) GetAccessToken() string {
	return storage.AccessToken
}

func (storage LocalTokenStorage) GetRefreshToken() string {
	return storage.RefreshToken
}

func (storage *LocalTokenStorage) SaveTokenEntity(access_token string, refresh_token string) {
	storage.AccessToken = access_token
	storage.RefreshToken = refresh_token

	storage.SaveToFile()
}

func (storage *LocalTokenStorage) SaveToFile() {
	log.Print("Trying to save tokens to file...")

	file, err := json.Marshal(storage)

	if err != nil {
		log.Print("Couldn't marshall storage")
		return
	}

	err = ioutil.WriteFile("tokens.json", file, os.ModePerm)

	if err != nil {
		log.Print("Couldn't save tokens to file")
		return
	}

	log.Print("Tokens saved to file")
}

func (storage *LocalTokenStorage) ReadFromFile() (LocalTokenStorage, error) {
	log.Print("Trying to read tokens from file...")

	file, _ := ioutil.ReadFile("tokens.json")

	data := LocalTokenStorage{}

	err := json.Unmarshal([]byte(file), &data)

	if err != nil {
		log.Print("Tokens not found in file")
	}

	return data, err
}

func NewTokenStorage() *LocalTokenStorage {
	once.Do(func() {

		instance = &LocalTokenStorage{}

		data, err := instance.ReadFromFile()

		if err == nil {
			instance = &data
		}

	})

	return instance
}
