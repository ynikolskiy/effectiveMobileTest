package db

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Подключение к базе данных
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения: ", err)
	}

	// Проверка подключения
	err = DB.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения: ", err)
	}

	log.Println("Успешное подключение")

	runMigrations()
}

// runMigrations выполняет миграции при старте
func runMigrations() {

	m, err := migrate.New(
		"file://migrations",
		"postgres://postgres:123@localhost:5432/musicdb?sslmode=disable",
	)
	if err != nil {
		log.Fatal("Ошибка создания миграции: ", err)
	}

	// Выполнение миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Ошибка выполнения миграции: ", err)
	}

	log.Println("Миграции выполнены успешно")
}
