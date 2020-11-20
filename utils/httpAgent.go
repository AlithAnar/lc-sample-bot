package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func CreateRequest(method string, url string, body []byte) ([]byte, error) {
	return CreateRequestWithCustomHeaders(method, url, body, nil)
}

func CreateRequestWithCustomHeaders(method string, url string, body []byte, extraHeaders map[string]string) ([]byte, error) {
	err := ValidateToken()

	if err != nil {
		err = RefreshToken()
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	tokenStorage := NewTokenStorage()
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenStorage.GetAccessToken()))

	if extraHeaders != nil {
		for key := range extraHeaders {
			request.Header.Add(key, extraHeaders[key])
		}
	}

	response, err := client.Do(request)

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	rawbody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	log.Printf(string(rawbody))

	return rawbody, nil
}
