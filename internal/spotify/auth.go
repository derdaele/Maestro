package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AuthProvider provide a valid bearer token.
type AuthProvider interface {
	getToken() string
}

type authError struct {
	ErrorName        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *authError) Error() string {
	return fmt.Sprintf("Authentication error [%s]: %s", e.ErrorName, e.ErrorDescription)
}

type authSuccess struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// ClientCredentials represent a client credential auth provider.
type ClientCredentials struct {
	ClientID     string
	ClientSecret string
	cachedToken  string
	expiration   time.Time
}

func getClientCredentialToken(clientID string, clientSecret string) (*authSuccess, error) {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	formData := url.Values{"grant_type": {"client_credentials"}}
	req, err := http.NewRequest("POST", authTokenEndpoint, strings.NewReader(formData.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprint("Basic ", basicAuth))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	payload, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		err := new(authError)
		jsonErr := json.Unmarshal(payload, err)
		if jsonErr != nil {
			return nil, jsonErr
		}
		return nil, err
	}

	var success authSuccess
	jsonErr := json.Unmarshal(payload, &success)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &success, nil
}

// GetToken retrieve the current cached token if not expired, retrieve one if needed
func (cc *ClientCredentials) getToken() string {
	if cc.expiration.Before(time.Now().Add(time.Minute)) {
		authSuccess, err := getClientCredentialToken(cc.ClientID, cc.ClientSecret)

		if err != nil {
			panic(err)
		}

		cc.expiration = time.Now().Add(time.Second * time.Duration(authSuccess.ExpiresIn))
		cc.cachedToken = authSuccess.AccessToken
	}

	return cc.cachedToken
}

// NewClientCredentials initialize a new client credential auth provider from a client id and secret.
func NewClientCredentials(clientID string, clientSecret string) *ClientCredentials {
	return &ClientCredentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		cachedToken:  "",
		expiration:   time.Now().Add(-time.Hour),
	}
}
