package config

import (
	"EuroprotocolTGBot/internal/storage"
	"flag"
	"log"
	"os"
)

type Config struct {
	Key        string
	Mode       string
	ConfigFile string
	DBConfig   storage.DBConfig
}

func NewConfig() Config {
	key := flag.String("k", "", "Telegram API token")
	mode := flag.String("m", "Debug", "Mode for loggin. Should be Zap like debug mode")
	cf := flag.String("f", "Asks.json", "Name of config file with answers in json format")
	flag.Parse()

	var args string

	for _, v := range os.Args {
		args = args + string(" ") + v
	}

	log.Println("os args: ", args)

	return Config{
		Key:        *key,
		Mode:       *mode,
		ConfigFile: *cf,
	}
}
