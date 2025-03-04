package config

import (
	"EuroprotocolTGBot/internal/loggin"
	"flag"
	"os"

	"go.uber.org/zap"
)

type Config struct {
	Key  string
	Mode string
}

func NewConffig() Config {
	key := flag.String("k", "", "Telegram API token")
	mode := flag.String("m", "Debug", "Mode for loggin. Should be Zap like debug mode")
	flag.Parse()

	var args string

	for _, v := range os.Args {
		args = args + string(" ") + v
	}

	loggin.Log.Debug("os: ", zap.String("args: ", args))

	return Config{
		Key:  *key,
		Mode: *mode,
	}
}
