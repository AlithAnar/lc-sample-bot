package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"lc-sample-bot/config"
	"log"
	"net/http"
	"time"
)

type AuthorizePayload struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	LicenceId    int    `json:"licence_id"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
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

	response, err := http.Post("https://accounts.labs.livechat.com/token", "application/json", bytes.NewBuffer(requestBody))

	defer response.Body.Close()

	if err != nil {
		return err
	}

	rawbody, err := ioutil.ReadAll(response.Body)

	var payload AuthorizePayload
	err = json.Unmarshal(rawbody, &payload)

	if err != nil {
		return err
	}

	tokenStorage := NewTokenStorage()
	tokenStorage.SaveTokenEntity(payload.AccessToken, payload.RefreshToken)

	return nil
}

func ValidateToken() error {
	tokenStorage := NewTokenStorage()
	request, err := http.NewRequest(http.MethodGet, "https://accounts.labs.livechat.com/v2/info", nil)

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenStorage.GetAccessToken()))

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode == 422 || response.StatusCode == 401 {
		return errors.New("Invalid token")
	}

	rawbody, err := ioutil.ReadAll(response.Body)

	log.Print(string(rawbody))

	var payload AuthorizePayload
	err = json.Unmarshal(rawbody, &payload)

	if err != nil {
		return err
	}

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

	res, err := http.Post("https://accounts.labs.livechat.com/v2/token", "application/json", bytes.NewBuffer(requestBody))

	defer res.Body.Close()

	if err != nil {
		return err
	}

	rawbody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

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
