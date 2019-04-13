package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Request perform an HTTP request and parse the JSON body
func (c Client) Request(req *http.Request, result interface{}) error {
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", c.Auth.getToken()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	for i := 0; i < maxRetryCount; i++ {
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return err
		}

		switch code := res.StatusCode; {
		case code == 429: // Too many request
			retryAfterSec, _ := strconv.Atoi(res.Header.Get("Retry-After"))
			time.Sleep(time.Second * time.Duration(retryAfterSec))

		case code >= 200 && code < 300: // Success
			payload, err := ioutil.ReadAll(res.Body)

			if err != nil {
				return err
			}

			return json.Unmarshal(payload, result)
		default:
			return fmt.Errorf("Unexpected status code %d", code)
		}
	}

	return nil
}

// Get perform a GET http request on the URL and extract the json body result
func (c Client) Get(URL string, result interface{}) error {
	req, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		return err
	}

	return c.Request(req, result)
}
