package api

import (
	"encoding/json"
	"io/ioutil"
	"lc-sample-bot/bot"
	"lc-sample-bot/common"
	"lc-sample-bot/config"
	"lc-sample-bot/utils"
	"log"
	"net/http"
)

func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("200 - cool"))
}

func TokenHandler(writer http.ResponseWriter, request *http.Request) {
	rawbody, err := ioutil.ReadAll(request.Body)

	enableCors(&writer)

	if err != nil {
		log.Fatalln(err)
	}

	var payload ExchangeTokenRequestPayload

	err = json.Unmarshal(rawbody, &payload)

	if err != nil {
		log.Fatalln(err)
	}

	err = utils.ExchangeToken(payload.Code)

	if err != nil {
		log.Fatalln(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - not cool"))
		return
	}

	config := config.NewAppConfig()
	bot.RegisterBot(config.WebhookUrl, config.BotSecret)

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("200 - cool"))
}

func RegisterBotHandler(writer http.ResponseWriter, request *http.Request) {
	rawbody, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var payload struct {
		WebhookUrl string `json:"webhook_url"`
		BotSecret  string `json:"bot_secret"`
	}

	err = json.Unmarshal(rawbody, &payload)

	if payload.WebhookUrl == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("400 - not cool"))
		return
	}

	err = bot.RegisterBot(payload.WebhookUrl, payload.BotSecret)

	if err != nil {
		log.Fatalln(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - not cool"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("200 - cool"))
}

func WebhookHandler(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)

	rawbody, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var webhookPayload json.RawMessage

	webhookBody := common.WebhookBody{
		Payload: &webhookPayload,
	}
	err = json.Unmarshal(rawbody, &webhookBody)

	log.Printf("[API] Incoming webhook: %s", webhookBody.Action)

	switch webhookBody.Action {
	case common.IncomingEvent:
		var payload common.WebhookChatEventPayload
		err := json.Unmarshal(webhookPayload, &payload)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		webhookBody.Payload = payload
	case common.IncomingChat:
		var payload common.WebhookChatPayload
		err := json.Unmarshal(webhookPayload, &payload)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		webhookBody.Payload = payload

	default:
		log.Printf("[API] Unhandled webhook action %s", webhookBody.Action)
		return
	}

	err = bot.HandleWebhook(webhookBody)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("200 - cool"))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
