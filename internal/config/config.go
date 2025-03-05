package config

import (
	"flag"
	"log"
	"os"
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

	log.Println("os args: ", args)

	return Config{
		Key:  *key,
		Mode: *mode,
	}
}
