package main

import (
	"lc-sample-bot/api"
	"lc-sample-bot/bot"
	"lc-sample-bot/config"
	"lc-sample-bot/utils"
	"log"
	"net/http"
)

func main() {
	utils.NewTokenStorage()
	bot.NewBotConfig()
	config.NewAppConfig()
	http.HandleFunc("/", api.IndexHandler)
	http.HandleFunc("/token", api.TokenHandler)
	http.HandleFunc("/register_bot", api.RegisterBotHandler)
	http.HandleFunc("/webhook", api.WebhookHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
