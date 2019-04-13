package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Request perform an HTTP request and parse the JSON body
func (c Client) Request(req *http.Request, result interface{}) error {
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", c.Auth.getToken()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	payload, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(payload, result)
}

// Get perform a GET http request on the URL and extract the json body result
func (c Client) Get(URL *url.URL, result interface{}) error {
	req, err := http.NewRequest("GET", URL.String(), nil)

	if err != nil {
		return err
	}

	return c.Request(req, result)
}
