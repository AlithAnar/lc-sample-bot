package bot

import "sync"

type BotChatActivity struct {
	sync.RWMutex
	items map[string]int
}

func (bca *BotChatActivity) GetChatActivity(chatId string) int {
	bca.Lock()
	defer bca.Unlock()
	return botChatActivity.items[chatId]
}

func (bca *BotChatActivity) IncrementChatActivity(chatId string) {
	bca.Lock()
	defer bca.Unlock()
	botChatActivity.items[chatId] += 1
}
