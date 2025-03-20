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
	host := flag.String("dbh", "localhost", "IP adress of postgresql Database")
	port := flag.String("dbp", "5432", "Port for postgresql databese connection")
	user := flag.String("dbu", "postgres", "User for postgresql database")
	pass := flag.String("dbpass", "sqlpwdusr", "Password for postgresql database")
	name := flag.String("dbn", "europrotocol", "Database name for postgresql")
	ocon := flag.Int("dbo", 300, "Maximum opened connections for postgresql database")
	icon := flag.Int("dbi", 150, "Maximum idle connections for postgresql database")
	flag.Parse()

	var args string

	for _, v := range os.Args {
		args = args + string(" ") + v
	}

	log.Println("os args: ", args)

	dbConfig := storage.DBConfig{
		Host:        *host,
		Port:        *port,
		User:        *user,
		Password:    *pass,
		DBName:      *name,
		MaxOpenCons: *ocon,
		MaxIdleCons: *icon,
	}

	return Config{
		Key:        *key,
		Mode:       *mode,
		ConfigFile: *cf,
		DBConfig:   dbConfig,
	}
}
