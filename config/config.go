package config

import (
	"log/slog"
	"os"
)

func Get_Client_id() string {
	data := os.Getenv("GOOGLE_CLIENT_ID")
	if data == "" {
		slog.Error("GOOGLE_CLIENT_ID not found")
		os.Exit(1)
	}
	return data
}

func Get_Client_secret() string {
	data := os.Getenv("GOOGLE_CLIENT_SECRET")
	if data == "" {
		slog.Error("GOOGLE_CLIENT_SECRET not found")
		os.Exit(1)
	}
	return data
}
