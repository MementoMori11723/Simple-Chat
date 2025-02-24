package database

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	database_url = "data/base.db"
	create_query = `CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY UNIQUE NOT NULL, 
    email TEXT NOT NULL UNIQUE, 
    password TEXT NOT NULL 
  ); pragma journal_mode = WAL;`
)

func init() {
  slog.Info("Initializing database")
  err := os.Mkdir("data", 0755)
  if err != nil {
    slog.Info("Directory already exists")
  }
	db := connect()
	if db == nil {
		slog.Error("Failed to connect to database")
		return
	}
	defer db.Close()
	_, err = db.Exec(create_query)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func connect() *sql.DB {
	db, err := sql.Open("sqlite3", database_url)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	return db
}

func Test() {
  slog.Info("Test")
  db := connect()
  if db == nil {
    slog.Error("Failed to connect to database")
    return
  } else {
    slog.Info("Connected to database")
  }
  defer db.Close()
  _, err := db.Exec("INSERT OR IGNORE INTO users (id, email, password) VALUES ('3', 'test-3', 'test-3');")
  if err != nil {
    slog.Error(err.Error())
    return
  }
}
