package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Frozelo/music-rate-service/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOAuthConfig *oauth2.Config
	stateString       string
)

func InitOauth(cfg *config.Config) {
	githubOAuthConfig = &oauth2.Config{
		ClientID:     cfg.Oauth.ClientID,
		ClientSecret: cfg.Oauth.ClientSecret,
		RedirectURL:  cfg.Oauth.RedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	stateString = cfg.Oauth.OauthStateString
}

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOAuthConfig.AuthCodeURL(stateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != stateString {
		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
		return
	}

	code := r.FormValue("code")
	token, err := githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	client := githubOAuthConfig.Client(context.Background(), token)
	userInfo, err := getUserInfo(client)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User Info: %s\n", userInfo)
}

func getUserInfo(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", user), nil
}
