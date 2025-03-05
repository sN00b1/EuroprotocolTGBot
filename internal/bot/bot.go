package bot

import (
	"EuroprotocolTGBot/internal/loggin"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	Bot *tgbotapi.BotAPI
}

func NewBot(apiTocken string, mode bool) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(apiTocken)
	if err != nil {
		loggin.Log.Error(err.Error())
		return nil, err
	}

	bot.Debug = mode

	loggin.Log.Debug("NewBot:", zap.String("Authorized on account %s", bot.Self.UserName))
	return &Bot{
		Bot: bot,
	}, nil
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			bot.Bot.Send(msg)
		}
	}
}
