package functions

import (
	"log/slog"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("X-json-token", "1234567890")
  username := r.FormValue("username")
  password := r.FormValue("password")
  slog.Info("login", "Username: ", username)
  slog.Info("login", "Password: ", password)
  http.Redirect(w, r, "/", http.StatusFound)
}

func Signup(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("X-json-token", "1234567890")
  username := r.FormValue("username")
  password := r.FormValue("password")
  email := r.FormValue("email")
  confirm_password := r.FormValue("confirm-password")
  slog.Info("signup", "Username: ", username)
  slog.Info("signup", "Password: ", password)
  slog.Info("signup", "Email: ", email)
  slog.Info("signup", "Confirm Password: ", confirm_password)
  if password != confirm_password {
    slog.Error("Passwords do not match")
  }
  http.Redirect(w, r, "/", http.StatusFound)
}
