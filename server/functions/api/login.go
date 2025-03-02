package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"simple-chat/server/functions"
	"strings"
	"time"

	"github.com/google/uuid"
)

var redirect = "/"

type request struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

type response struct {
	Token string `json:"token"`
	Route string `json:"route"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	data := request{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}
	if err := data.validate(false); err != nil {
		responseError(w, err)
		return
	}
	token, err := generate_token(data, true)
	if err != nil {
		responseError(w, err)
		return
	}
	response{
		Token: token,
		Route: redirect,
	}.write(w)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm_password := r.FormValue("confirm-password")
	data := request{
		Email:           strings.TrimSpace(email),
		Username:        strings.TrimSpace(username),
		Password:        strings.TrimSpace(password),
		ConfirmPassword: strings.TrimSpace(confirm_password),
	}
	if err := data.validate(true); err != nil {
		responseError(w, err)
		return
	}
	token, err := generate_token(data, false)
	if err != nil {
		responseError(w, err)
		return
	}
	response{
		Token: token,
		Route: redirect,
	}.write(w)
}

func (r response) write(w http.ResponseWriter) {
	data, err := json.Marshal(r)
	if err != nil {
		slog.Error("Error marshalling response")
	}
	w.Header().Set("X-Auth-Token", string(data))
}

func (r request) validate(is_signup bool) error {
	if r.Username == "" || r.Password == "" {
		return errors.New("Username or password is required")
	}
	if is_signup {
		if r.Email == "" {
			return errors.New("Email is required")
		}
		if r.ConfirmPassword == "" {
			return errors.New("Confirm password is required")
		}
		if r.Password != r.ConfirmPassword {
			return errors.New("Passwords do not match")
		}
	}
	return nil
}

func responseError(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func generate_token(r request, is_login bool) (string, error) {
	if is_login {
		key := functions.Get_keys(
			fmt.Sprintf(
				"user:%s:passwd:%s:email:*",
				r.Username, r.Password,
			),
		)
		data := functions.Get_data(key)
		if data == "" {
			return "", errors.New("User not found")
		}
		var json_data map[string]string
		err := json.Unmarshal([]byte(data), &json_data)
		if err != nil {
			return "", err
		}
		header := base64.StdEncoding.EncodeToString([]byte(json_data["id"]))
		base_data := base64.StdEncoding.EncodeToString([]byte(data))
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s.%s", header, base_data)))
		return fmt.Sprintf("%s.%s.%s", header, base_data, hex.EncodeToString(hash[:])), nil
	}
	data, err := json.Marshal(map[string]string{
		"id": uuid.New().String(),
	})
	if err != nil {
		return "", err
	}
	functions.Add(
		fmt.Sprintf(
			"user:%s:passwd:%s:email:%s",
			r.Username, r.Password, r.Email,
		),
		string(data), time.Duration(9999*time.Hour),
	)
	return r.Username, nil
}
