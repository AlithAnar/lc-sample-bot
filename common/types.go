package common

type Dictionary map[string]interface{}

type ActionType string

const (
	IncomingEvent ActionType = "incoming_event"
	IncomingChat             = "incoming_chat"
)

type RecipientsType string

const (
	All    RecipientsType = "all"
	Agents                = "agents"
)

type WebhookBody struct {
	WebhookId string      `json:"webhook_id"`
	SecretKey string      `json:"secret_key"`
	Action    ActionType  `json:"action"`
	Payload   interface{} `json:"payload"`
}

type WebhookChatEventPayload struct {
	ChatId   string    `json:"chat_id"`
	ThreadId string    `json:"thread_id"`
	Event    ChatEvent `json:"event"`
}

type WebhookChatPayload struct {
	Chat ChatPayload `json:"chat"`
}

type ChatPayload struct {
	Id     string        `json:"id"`
	Thread ThreadPayload `json:"thread"`
}

type ThreadPayload struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Active    bool   `json:"active"`
}

type ChatEventType string

const (
	Message     ChatEventType = "message"
	File                      = "file"
	FilledForm                = "filled_form"
	RichMessage               = "rich_message"
	Custom                    = "custom"
)

type ChatEvent struct {
	Type       ChatEventType  `json:"type"`
	Text       string         `json:"text"`
	Recipients RecipientsType `json:"recipients"`
	AuthorId   string         `json:"author_id"`
}

type ChatEventPayload struct {
	ChatId string    `json:"chat_id"`
	Event  ChatEvent `json:"event"`
}
