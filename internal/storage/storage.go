package storage

import (
	"EuroprotocolTGBot/internal/loggin"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Структура конфигурации подключения
type DBConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	MaxOpenCons int
	MaxIdleCons int
}

// Создаем пул подключений
func createConnectionPool(config DBConfig) *sql.DB {
	// Формируем строку подключения
	connString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User,
		config.Password, config.DBName)

	// Создаем пул подключений
	db, err := sql.Open("postgres", connString)
	if err != nil {
		loggin.Log.Fatal("Ошибка при открытии соединения: ", zap.String("DBError", err.Error()))
	}

	// Настраиваем параметры пула
	db.SetMaxOpenConns(config.MaxOpenCons) // Максимальное количество открытых подключений
	db.SetMaxIdleConns(config.MaxIdleCons) // Максимальное количество idle-подключений
	db.SetConnMaxLifetime(time.Minute * 5) // Время жизни подключения

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		loggin.Log.Fatal("Ошибка при проверке подключения: ", zap.String("DBError", err.Error()))
	}

	return db
}

type DBStorage struct {
	DB       *sql.DB
	IsActive bool
	wg       sync.WaitGroup
}

func NewDBStorage(cfg DBConfig) (*DBStorage, error) {
	db := createConnectionPool(cfg)

	err := db.Ping()
	isActive := true
	if err != nil {
		isActive = false
		return nil, err
	}

	return &DBStorage{
		DB:       db,
		IsActive: isActive,
	}, nil
}
