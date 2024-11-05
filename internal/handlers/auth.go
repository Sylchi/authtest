package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"oauth-demo/internal/config"
)

var (
	store  *sessions.CookieStore
	Config *config.Config
)

func InitHandlers(sessionStore *sessions.CookieStore, cfg *config.Config) {
	store = sessionStore
	Config = cfg
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the home page! <a href='/login'>Login</a>"))
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	oauthConfig := oauth2.Config{
		ClientID:     Config.OAuth.ClientID,
		ClientSecret: Config.OAuth.ClientSecret,
		RedirectURL:  Config.OAuth.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  Config.OAuth.AuthURL,
			TokenURL: Config.OAuth.TokenURL,
		},
		Scopes: []string{"openid", "profile", "email"},
	}

	url := oauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	oauthConfig := oauth2.Config{
		ClientID:     Config.OAuth.ClientID,
		ClientSecret: Config.OAuth.ClientSecret,
		RedirectURL:  Config.OAuth.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  Config.OAuth.AuthURL,
			TokenURL: Config.OAuth.TokenURL,
		},
	}

	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", Config.OAuth.UserInfoURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["email"] = userInfo["email"]
	session.Save(r, w)

	http.Redirect(w, r, "/protected", http.StatusTemporaryRedirect)
}

func HandleProtected(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userEmail, ok := session.Values["email"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Protected page. Welcome " + userEmail + "! <a href='/logout'>Logout</a>"))
}
