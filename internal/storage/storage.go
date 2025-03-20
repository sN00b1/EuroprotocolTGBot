package storage

import (
	"EuroprotocolTGBot/internal/loggin"
	"database/sql"
	"fmt"
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
	DB        *sql.DB
	IsActive  bool
	Semaphore chan struct{}
}

func NewDBStorage(cfg DBConfig) (*DBStorage, error) {
	db := createConnectionPool(cfg)

	err := db.Ping()
	if err != nil {
		return &DBStorage{
			DB:        nil,
			IsActive:  false,
			Semaphore: make(chan struct{}, cfg.MaxOpenCons),
		}, err
	}

	dbStorage := &DBStorage{
		DB:        db,
		IsActive:  true,
		Semaphore: make(chan struct{}, cfg.MaxOpenCons),
	}

	err = dbStorage.createDBStruct()
	if err != nil {
		return nil, err
	}

	return dbStorage, err
}

func (db *DBStorage) Close() {
	db.DB.Close()
	db.IsActive = false
}

func (db *DBStorage) createDBStruct() error {
	createQuery := `
		CREATE TABLE IF NOT EXISTS europrotocol (
		id VARCHAR(255) PRIMARY KEY,
		str1 TEXT,
		str2 TEXT,
		str3 TEXT,
		str4 TEXT,
		str5 TEXT,
		str6 TEXT,
		str7 TEXT,
		str8 TEXT,
		str9 TEXT,
		str10 TEXT,
		str11 TEXT,
		str12 TEXT,
		str13 TEXT,
		str14 TEXT,
		str15 TEXT,
		str16 TEXT,
		str17 TEXT,
		str18 TEXT,
		str19 TEXT,
		str20 TEXT,
		str21 TEXT,
		str22 TEXT,
		str23 TEXT,
		str24 TEXT,
		str25 TEXT,
		str26 TEXT,
		str27 TEXT,
		str28 TEXT,
		str29 TEXT,
		str30 TEXT,
		str31 TEXT,
		str32 TEXT,
		str33 TEXT,
		str34 TEXT,
		str35 TEXT,
		str36 TEXT,
		str37 TEXT,
		str38 TEXT,
		str39 TEXT,
		str40 TEXT,
		str41 TEXT,
		str42 TEXT,
		str43 TEXT,
		str44 TEXT,
		str45 TEXT,
		str46 TEXT,
		str47 TEXT,
		str48 TEXT,
		str49 TEXT,
		str50 TEXT,
		str51 TEXT,
		str52 TEXT,
		str53 TEXT,
		str54 TEXT,
		str55 TEXT,
		str56 TEXT,
		);`

	_, err := db.DB.Exec(createQuery)
	if err != nil {
		loggin.Log.Debug("err:", zap.String("err:", err.Error()))
		return err
	}

	return nil
}
