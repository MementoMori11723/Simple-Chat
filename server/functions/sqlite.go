package functions

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "/data/chat.db")
	return db, err
}

func TestConnection() {
	db, err := connect()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	} else {
    slog.Info("Connected to database!")
  }
}

func Add(string, string, time.Duration) {}

func Get_data(key string) string {
	return key
}

func Get_keys(key string) string {
	return key
}

func Delete() {}

func Encrypt() {}

func Decrypt() {}
