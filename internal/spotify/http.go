package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (c Client) request(req *http.Request, result interface{}) error {
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

func (c Client) get(URL *url.URL, result interface{}) error {
	req, err := http.NewRequest("GET", URL.String(), nil)

	if err != nil {
		return err
	}

	return c.request(req, result)
}
