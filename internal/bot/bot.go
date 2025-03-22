package bot

import (
	"EuroprotocolTGBot/internal/config"
	"EuroprotocolTGBot/internal/loggin"
	"EuroprotocolTGBot/internal/storage"
	"EuroprotocolTGBot/internal/tools"
	"os"
	"os/exec"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var FatherChan tools.MsgChain

// структура бота
type Bot struct {
	Bot   *tgbotapi.BotAPI
	cfg   config.Config
	chain map[string]*tools.MsgChain
	db    *storage.DBStorage
	mux   sync.RWMutex
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
	FatherChan = tools.NewMsgChain()
	FatherChan.LoadAsks(cfg.ConfigFile)

	loggin.Log.Debug("NewBot:", zap.String("Authorized on account %s", bot.Self.UserName))

	// создаем БД
	db, err := storage.NewDBStorage(cfg.DBConfig)
	if err != nil {
		loggin.Log.Error("NewDBStorage:", zap.String("Database creation fail: ", err.Error()))
		return &Bot{
			Bot:   bot,
			cfg:   cfg,
			chain: make(map[string]*tools.MsgChain),
			db:    db,
		}, err
	}

	return &Bot{
		Bot:   bot,
		cfg:   cfg,
		chain: make(map[string]*tools.MsgChain),
		db:    db,
	}, nil
}

// функция асинхронного добавления пользовательской сессии
func (bot *Bot) AddChain(id string, chain tools.MsgChain) {
	bot.mux.Lock()
	bot.chain[id] = &chain
	bot.mux.Unlock()
}

// функция получения данных из пользовательской сессии
func (bot *Bot) PrintAnswer(id string) string {
	var result []rune

	bot.mux.RLock()
	for k, v := range bot.chain[id].AskList {
		ask := []rune(v.Text)
		answer := []rune(bot.chain[id].AnswerList[k].Text)
		result = append(result, ask...)
		result = append(result, '\n')
		result = append(result, answer...)
		result = append(result, '\n')
	}
	bot.mux.RUnlock()
	return string(result)
}

// функция запуска бота
func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil || update.CallbackQuery != nil {
			go bot.Handle(update)
		}
	}
}

func (bot *Bot) Handle(update tgbotapi.Update) {
	var id string
	var msg tgbotapi.MessageConfig

	id = update.FromChat().UserName

	_, ok := bot.chain[id]
	if !ok {
		bot.AddChain(id, FatherChan)
	}

	if update.Message != nil {
		switch update.Message.Text {
		case "/start":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я ваш Telegram-бот по заполнению европротоклов при ДТП. Напишите /help для просмотра всех команд.")
		case "/help":
			msg = bot.HelpHandler(id, update)
		case "/new":
			msg = bot.NewHandler(id, update)
		case "/list":
			msg = DevelopmentHandler(update)
		case "/new_on":
			msg = DevelopmentHandler(update)
		case "/edit":
			msg = DevelopmentHandler(update)
		case "/print":
			msg = bot.PrintHandler(id, update)
		default:
			msg = bot.DefualtHandler(id, update)
		}
	}

	// Обработка нажатий кнопок
	if update.CallbackQuery != nil {
		callback := update.CallbackQuery
		switch callback.Data {
		case "next":
			bot.chain[id].Next()
			ask, _ := bot.chain[id].GetCurrentAsk()
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ask.Text)
			msg.ReplyMarkup = GetNaviKeyBoard() // используем ту же клавиатуру
		case "back":
			bot.chain[id].Back()
			ask, _ := bot.chain[id].GetCurrentAsk()
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ask.Text)
			msg.ReplyMarkup = GetNaviKeyBoard()
		case "cancel":
			bot.chain[id].Reset()
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Заполнение протокола отменено. Воспользутесь /help.")
		}
	}

	bot.Bot.Send(msg)
}

func DevelopmentHandler(update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Данная функция находится в разработке.")
}

func GetNaviKeyBoard() tgbotapi.InlineKeyboardMarkup {
	// Создаем клавиатуру с кнопками
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅ Назад", "back"),
			tgbotapi.NewInlineKeyboardButtonData("Далее ➡", "next"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить", "cancel"),
		),
	)
	return keyboard
}

func (bot *Bot) DefualtHandler(id string, update tgbotapi.Update) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	if bot.chain[id].Start {
		bot.chain[id].SetCurrentAnswer(update.Message.Text)
		ask, ok := bot.chain[id].GetCurrentAsk()
		if ok {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, ask.Text)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Заполнение протокола завершено. Можете его распечатать командой /print.")
		}
		msg.ReplyMarkup = GetNaviKeyBoard()
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваше сообщение не ясно. Воспользуйтесь /help.")
	}
	return msg
}

func (bot *Bot) HelpHandler(id string, update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID,
		`Я могу помочь с основными вопросами. 
	Напишите /new для создания нового протокола, 
	/list для ппросмотра списка своих протоколов,
	/new_on для создания протокола на основе существующего протокола,
	/edit для редактирования проткола из списка
	/print для печати европротокола.`)
}

func (bot *Bot) NewHandler(id string, update tgbotapi.Update) tgbotapi.MessageConfig {
	bot.chain[id].Reset()
	ask, _ := bot.chain[id].GetCurrentAsk()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, ask.Text)
	msg.ReplyMarkup = GetNaviKeyBoard()

	return msg
}

func (bot *Bot) PrintHandler(id string, update tgbotapi.Update) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	outJsonFP := bot.cfg.BinPath + "/" + strconv.FormatInt(update.Message.Chat.ID, 10) + ".json"
	docxFile := bot.cfg.BinPath + "/" + strconv.FormatInt(update.Message.Chat.ID, 10) + ".docx"

	err := bot.chain[id].PrintAnswer(outJsonFP)
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка создания протокола. Данные инцедента отправлены администратору.")
		return msg
	}

	// Запуск Python скрипта с параметрами
	cmd := exec.Command("/usr/bin/python3", bot.cfg.ScriptPath, outJsonFP, docxFile)

	// Получение вывода
	out, err := cmd.CombinedOutput()
	if err != nil {
		loggin.Log.Error("PrintHandler:", zap.String("Ошибка запуска скрипта", err.Error()))
	}

	loggin.Log.Info("PrintHandler: ", zap.String("mkdoc.py output:", string(out)))

	filePath := docxFile

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		loggin.Log.Error("PrintHandler:", zap.String("Ошибка чтения файла, полученного от скрипта.", err.Error()))
	}
	docxFileBytes := tgbotapi.FileBytes{
		Name:  "Европротокол.docx",
		Bytes: fileBytes,
	}

	bot.Bot.Send(tgbotapi.NewDocument(update.Message.Chat.ID, docxFileBytes))

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш европротокол успешно сгенерирован.")

	return msg
}
