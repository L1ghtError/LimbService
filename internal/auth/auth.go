package auth

import (
	"encoding/json"
	"fmt"
	"io"
	model "light-backend/internal/model"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func ConfigGoogle() *oauth2.Config {
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	conf := &oauth2.Config{ClientID: os.Getenv("Client"),
		ClientSecret: os.Getenv("Secret"),
		RedirectURL:  os.Getenv("redirect_url"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid"},
		Endpoint: google.Endpoint,
	}
	return conf
}

func GetGoogleResponse(token string) (*model.GoogleResponse, error) {
	reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return nil, err
	}

	ptoken := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{Method: "GET", URL: reqURL, Header: map[string][]string{"Authorization": {ptoken}}}

	req, err := http.DefaultClient.Do(res)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}
	var data model.GoogleResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
