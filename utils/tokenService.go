package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"lc-sample-bot/config"
	"log"
	"net/http"
)

type AuthorizePayload struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	LicenceId    int    `json:"licence_id,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

func ExchangeToken(code string) error {
	config := config.NewAppConfig()
	requestBody, err := json.Marshal(map[string]string{
		"code":          code,
		"grant_type":    "authorization_code",
		"client_id":     config.AppClientId,
		"client_secret": config.AppSecretId,
		"redirect_uri":  config.RedirectURI,
	})

	if err != nil {
		return err
	}

	response, err := CreateRequest(http.MethodPost, "https://accounts.labs.livechat.com/token", requestBody)

	if err != nil {
		return err
	}

	var payload AuthorizePayload
	err = json.Unmarshal(response, &payload)

	if err != nil {
		return err
	}

	tokenStorage := NewTokenStorage()
	tokenStorage.SaveTokenEntity(payload.AccessToken, payload.RefreshToken)

	return nil
}

func RefreshToken() error {
	config := config.NewAppConfig()
	tokenStorage := NewTokenStorage()

	requestBody, err := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     config.AppClientId,
		"client_secret": config.AppSecretId,
		"refresh_token": tokenStorage.RefreshToken,
	})

	if err != nil {
		return err
	}

	res, err := http.NewRequest(http.MethodPost, "https://accounts.labs.livechat.com/v2/token", bytes.NewBuffer(requestBody))

	defer res.Body.Close()

	if err != nil {
		return err
	}

	rawbody, err := ioutil.ReadAll(res.Body)

	log.Print(string(rawbody))

	var payload AuthorizePayload
	err = json.Unmarshal(rawbody, &payload)

	if err != nil {
		return err
	}

	if payload.AccessToken != "" && payload.RefreshToken != "" {
		tokenStorage.SaveTokenEntity(payload.AccessToken, payload.RefreshToken)
	} else {
		return errors.New("Malformed token response")
	}

	return nil
}
