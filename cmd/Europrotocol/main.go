package main

import (
	"EuroprotocolTGBot/internal/bot"
	"EuroprotocolTGBot/internal/config"
	"EuroprotocolTGBot/internal/loggin"
)

func main() {
	config := config.NewConffig()

	err := loggin.Initialize(config.Mode)
	if err != nil {
		panic("Error while zap initialize.")
	}

	botObj, err := bot.NewBot(config.Key, config.Mode != "Release")
	if err != nil {
		panic("Error while bot initialize.")
	}

	botObj.Run()
}
