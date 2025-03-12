package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"simple-chat/config"
	"simple-chat/server/functions"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	redirect = "/"

	client_id     string
	client_secret string

	redirect_uri  = "http://localhost:11000/google/callback"
	auth_url      = "https://accounts.google.com/o/oauth2/auth"
	google_token  = "https://oauth2.googleapis.com/token"
	user_info_url = "https://www.googleapis.com/oauth2/v2/userinfo"
)

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

func init() {
	client_id = config.Get_Client_id()
	client_secret = config.Get_Client_secret()
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

func GoogleAuth(w http.ResponseWriter, r *http.Request) {
  slog.Info(client_id)
	json.NewEncoder(w).Encode(map[string]string{
		"url": fmt.Sprintf(
			"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=email+profile",
			auth_url, client_id, redirect_uri,
		),
	})
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}
	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("redirect_uri", redirect_uri)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	res, err := http.PostForm(google_token, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		http.Error(w, "Failed to parse token response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req, _ := http.NewRequest("GET", user_info_url, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	userResp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()
	body, _ := io.ReadAll(userResp.Body)
  // need to modify this later!
  json.Marshal(body)
  http.Redirect(w, r, "/", http.StatusFound)
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

// need to adjust this function also!
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
