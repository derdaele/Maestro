package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func getBearerToken(clientID string, clientSecret string) string {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	formData := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(formData.Encode()))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprint("Basic ", basicAuth))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	payload, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var authResponse authResponse
	err = json.Unmarshal(payload, &authResponse)

	if err != nil {
		panic(err)
	}

	return authResponse.AccessToken
}
