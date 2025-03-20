package bot

import (
	"EuroprotocolTGBot/internal/config"
	"EuroprotocolTGBot/internal/loggin"
	"EuroprotocolTGBot/internal/storage"
	"EuroprotocolTGBot/internal/tools"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// структура бота
type Bot struct {
	Bot   *tgbotapi.BotAPI
	cfg   config.Config
	chain tools.MsgChain
	db    *storage.DBStorage
}

// функция создания нового бота
func NewBot(cfg config.Config) (*Bot, error) {
	// подключаемся к апи телеграмма
	bot, err := tgbotapi.NewBotAPI(cfg.Key)
	if err != nil {
		loggin.Log.Error(err.Error())
		return nil, err
	}

	bot.Debug = (cfg.Mode != "Release")

	// получаем цепочку сообщений для опроса и создания европротокола
	chain := tools.NewMsgChain()
	chain.LoadAsks(cfg.ConfigFile)

	loggin.Log.Debug("NewBot:", zap.String("Authorized on account %s", bot.Self.UserName))

	// создаем БД
	db, err := storage.NewDBStorage(cfg.DBConfig)
	if err != nil {
		loggin.Log.Error("NewDBStorage:", zap.String("Database creation fail: ", err.Error()))
		return &Bot{
			Bot:   bot,
			cfg:   cfg,
			chain: *chain,
			db:    db,
		}, err
	}

	return &Bot{
		Bot:   bot,
		cfg:   cfg,
		chain: *chain,
		db:    db,
	}, nil
}

// функция запуска бота
func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			loggin.Log.Debug("", zap.String("User: ", update.Message.From.UserName+" Message: "+update.Message.Text))

			var msg tgbotapi.MessageConfig
			switch update.Message.Text {
			case "/start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я ваш Telegram-бот по заполнению европротоклов при ДТП. Напишите /help для просмотра всех команд.")
			case "/help":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID,
					`Я могу помочь с основными вопросами. 
				Напишите /new для создания нового протокола, 
				/list для ппросмотра списка своих протоколов,
				/new_on для создания протокола на основе существующего протокола,
				/edit для редактирования проткола из списка .`)
			case "/new":
				bot.chain.Reset()
				ask, _ := bot.chain.GetCurrentAsk()
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, ask.Text)
			case "/list":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Данная функция находится в разработке.")
			case "/new_on":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Данная функция находится в разработке.")
			case "/edit":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Данная функция находится в разработке.")
			default:
				if bot.chain.Start {
					bot.chain.SetCurrentAnswer(update.Message.Text)
					ask, ok := bot.chain.GetCurrentAsk()
					if ok {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, ask.Text)
					} else {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Заполнение протокола завершено. Можете его распечатать командой /print.")
					}
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваше сообщение не ясно. Воспользуйтесь /help.")
				}
			}

			bot.Bot.Send(msg)
		}
	}
}
