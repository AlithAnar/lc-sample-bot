package bot

import (
	"encoding/json"
	"lc-sample-bot/common"
	"lc-sample-bot/utils"
	"log"
	"net/http"
)

var botChatActivity = make(map[string]int)

func RegisterBot(webhookUrl string, botSecret string) error {
	body, err := json.Marshal(BotPayload{
		Name: "LC SIMPLE BOT",
		Webhooks: WebhooksPayload{
			URL:       webhookUrl,
			SecretKey: botSecret,
			Actions: []ActionPayload{
				{
					Name: "incoming_chat",
				},
				{
					Name: "incoming_event",
				},
			},
		},
	})

	if err != nil {
		return err
	}

	response, err := utils.CreateRequest(http.MethodPost, "https://api.labs.livechatinc.com/v3.2/configuration/action/create_bot", body)

	botConfig := NewBotConfig()

	json.Unmarshal(response, &botConfig)

	botConfig.SaveToFile()

	return SetBotRouting(botConfig.BotAgentId, "accepting_chats")
}

func SetBotRouting(agentId string, status string) error {
	body, err := json.Marshal(RoutingStatusPayload{
		Status:  status,
		AgentId: agentId,
	})

	if err != nil {
		return err
	}

	_, err = utils.CreateRequest(http.MethodPost, "https://api.labs.livechatinc.com/v3.2/agent/action/set_routing_status", body)

	return err
}

func HandleWebhook(webhookBody common.WebhookBody) error {
	switch webhookBody.Action {
	case common.IncomingChat:
		return handleIncomingChat(webhookBody.Payload.(common.WebhookChatPayload))
	case common.IncomingEvent:
		return handleIncomingChatEvent(webhookBody.Payload.(common.WebhookChatEventPayload))
	default:
		log.Printf("[BOT] Unhandled webhook action %s", webhookBody.Action)
		return nil
	}
}

func handleIncomingChat(payload common.WebhookChatPayload) error {
	log.Printf("[BOT] Handling incoming chat. ChatID: %s", payload.Chat.Id)
	return sendMessageToCustomer(payload.Chat.Id, "Hi dude!")
}

func handleIncomingChatEvent(payload common.WebhookChatEventPayload) error {
	switch payload.Event.Type {
	case common.Message:
		return handleIncomingMessageEvent(payload)
	default:
		log.Printf("[BOT] skipping %s event", payload.Event.Type)
		return nil
	}
}

func handleIncomingMessageEvent(payload common.WebhookChatEventPayload) error {
	botConfig := NewBotConfig()
	if payload.Event.AuthorId == botConfig.BotAgentId {
		return nil
	}
	log.Printf("[BOT] Handling incoming chat event. ChatID: %s Type: %s", payload.ChatId, payload.Event.Type)

	currentChatActivity := botChatActivity[payload.ChatId]

	var err error

	switch currentChatActivity {
	case 0:
		err = sendMessageToCustomer(payload.ChatId, "Cool, tell me more :)")
	case 1:
		err = sendMessageToCustomer(payload.ChatId, "Haha, that's funny ^^")
	case 2:
		err = sendMessageToCustomer(payload.ChatId, "I don't want to talk with you anymore ;/")
	default:
	}
	botChatActivity[payload.ChatId] += 1

	return err
}

func sendMessageToCustomer(chatId string, messageText string) error {
	botConfig := NewBotConfig()

	log.Printf("[BOT] Sending message to ChatID: %s as Bot:", chatId, botConfig.BotAgentId)

	body, err := json.Marshal(common.ChatEventPayload{
		ChatId: chatId,
		Event: common.ChatEvent{
			Recipients: common.All,
			Type:       common.Message,
			Text:       messageText,
		},
	})

	if err != nil {
		return err
	}

	customHeaders := map[string]string{
		"X-Author-Id": botConfig.BotAgentId,
	}

	_, err = utils.CreateRequestWithCustomHeaders(
		http.MethodPost,
		"https://api.labs.livechatinc.com/v3.2/agent/action/send_event",
		body,
		customHeaders,
	)

	return err
}
