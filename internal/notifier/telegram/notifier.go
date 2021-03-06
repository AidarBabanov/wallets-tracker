package telegram

import (
	"github.com/beego/beego/v2/core/logs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const parseModeMarkdown = "markdown"

type Notifier struct {
	*tgbotapi.BotAPI
	chatID int64
}

func New(tlgToken string, chatID int64) (*Notifier, error) {
	n := new(Notifier)
	bot, err := tgbotapi.NewBotAPI(tlgToken)
	if err != nil {
		return nil, err
	}
	n.BotAPI = bot
	n.chatID = chatID

	return n, nil
}

func (n *Notifier) Notify(msg string) {
	msg = "```\n" + msg + "\n```"
	msgConfig := tgbotapi.NewMessage(n.chatID, msg)
	msgConfig.ParseMode = parseModeMarkdown
	_, err := n.Send(msgConfig)
	if err != nil {
		logs.Error(err)
	}
}
