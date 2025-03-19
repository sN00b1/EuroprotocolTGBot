package main

import (
	"EuroprotocolTGBot/internal/bot"
	"EuroprotocolTGBot/internal/config"
	"EuroprotocolTGBot/internal/loggin"
)

func main() {
	config := config.NewConfig()

	err := loggin.Initialize(config.Mode)
	if err != nil {
		panic("Error while zap initialize.")
	}

	bot, err := bot.NewBot(config)
	if err != nil {
		panic("Error while bot initialize.")
	}

	bot.Run()
}
