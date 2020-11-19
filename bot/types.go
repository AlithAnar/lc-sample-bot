package bot

type BotPayload struct {
	Name     string `json:"name"`
	Webhooks WebhooksPayload
}

type RoutingStatusPayload struct {
	Status  string `json:"status"`
	AgentId string `json:"agent_id"`
}

type WebhooksPayload struct {
	URL       string `json:"url"`
	SecretKey string `json:"secret_key"`
	Actions   []ActionPayload
}

type ActionPayload struct {
	Name string `json:"name"`
}

type BotConfig struct {
	BotAgentId string `json:"bot_agent_id"`
}
