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

// объект владельца автомобиля
type Owner struct {
}

// объект водителя автомобиля
type Driver struct {
}

// объект ДТП
type Incedent struct {
}

// объект автомобиля
type Car struct {
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
	CREATE TABLE IF NOT EXISTS "Car" (
	"id" varchar(2000) NOT NULL UNIQUE,
	"owner_id" varchar(2000) NOT NULL UNIQUE,
	"driver_id" varchar(2000) NOT NULL UNIQUE,
	"model" varchar(2000) NOT NULL,
	"vin" varchar(2000) NOT NULL,
	"srp" varchar(15) NOT NULL,
	"cert" varchar(2000) NOT NULL,
	"insurance" varchar(2000) NOT NULL,
	"insurance_number" varchar(2000) NOT NULL,
	"casco" varchar(2000) NOT NULL,
	"damage" varchar(2000) NOT NULL,
	"info" varchar(2000) NOT NULL,
	PRIMARY KEY ("id")
	);

	CREATE TABLE IF NOT EXISTS "Incident" (
	"id" serial NOT NULL UNIQUE,
	"car_a_id" varchar(2000) NOT NULL,
	"car_b_id" varchar(2000) NOT NULL,
	"car_a_obs_id" varchar(2000) NOT NULL,
	"car_b_obs_id" varchar(2000) NOT NULL,
	"place" varchar(2000) NOT NULL,
	"date" timestamp without time zone NOT NULL,
	"cars_number" bigint NOT NULL DEFAULT '2',
	"hurt_number" bigint NOT NULL,
	"dead_number" bigint NOT NULL,
	"alco_check" boolean NOT NULL,
	"other_cars" boolean NOT NULL,
	"other_property" boolean NOT NULL,
	"witness" varchar(2000) NOT NULL,
	"gai_work" boolean NOT NULL,
	"gai_number" varchar(2000) NOT NULL,
	"other_cars_info" varchar(2000) NOT NULL,
	"damaged_property" varchar(2000) NOT NULL,
	"dproperty_owner" varchar(2000) NOT NULL,
	"new_field" bigint NOT NULL,
	PRIMARY KEY ("id")
	);

	CREATE TABLE IF NOT EXISTS "Owner" (
	"id" serial NOT NULL UNIQUE,
	"name" varchar(2000) NOT NULL,
	"adress" varchar(2000) NOT NULL,
	PRIMARY KEY ("id")
	);

	CREATE TABLE IF NOT EXISTS "Driver" (
	"id" serial NOT NULL UNIQUE,
	"fio" varchar(2000) NOT NULL,
	"birthday" timestamp without time zone NOT NULL,
	"adress" varchar(2000) NOT NULL,
	"phone" varchar(20) NOT NULL,
	"certificate" varchar(2000) NOT NULL,
	"way_doc" varchar(2000) NOT NULL,
	PRIMARY KEY ("id")
	);

	CREATE TABLE IF NOT EXISTS "Case" (
	"id" serial NOT NULL UNIQUE,
	"fact" varchar(2000) NOT NULL,
	"owner_drive" boolean NOT NULL,
	"own_move" boolean NOT NULL,
	"note" bigint NOT NULL,
	PRIMARY KEY ("id")
	);

	ALTER TABLE "Incident" ADD CONSTRAINT "Incident_fk0" FOREIGN KEY ("id") REFERENCES "Car"("incident_id");
	ALTER TABLE "Incident" ADD CONSTRAINT "Incident_fk1" FOREIGN KEY ("car_a_id") REFERENCES "Car"("id");
	ALTER TABLE "Incident" ADD CONSTRAINT "Incident_fk2" FOREIGN KEY ("car_b_id") REFERENCES "Car"("id");
	ALTER TABLE "Incident" ADD CONSTRAINT "Incident_fk3" FOREIGN KEY ("car_a_obs_id") REFERENCES "Case"("id");
	ALTER TABLE "Incident" ADD CONSTRAINT "Incident_fk4" FOREIGN KEY ("car_b_obs_id") REFERENCES "Case"("id");
	ALTER TABLE "Owner" ADD CONSTRAINT "Owner_fk0" FOREIGN KEY ("id") REFERENCES "Car"("owner_id");
	ALTER TABLE "Driver" ADD CONSTRAINT "Driver_fk0" FOREIGN KEY ("id") REFERENCES "Car"("driver_id");`

	_, err := db.DB.Exec(createQuery)
	if err != nil {
		loggin.Log.Debug("err:", zap.String("err:", err.Error()))
		return err
	}

	return nil
}

// Вспомогательная функция для соединения строк
func join(elements []string, separator string) string {
	if len(elements) == 0 {
		return ""
	}
	result := elements[0]
	for _, s := range elements[1:] {
		result += separator + s
	}
	return result
}
