package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Image represent an image url with its size.
type Image struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

// AlbumArtist represent an artist/composer/...
type AlbumArtist struct {
	ExternalUrls map[string]string `json:"external_urls`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
}

// Album is a collection of tracks.
type Album struct {
	AlbumGroup           string            `json:"album_group"`
	AlbumType            string            `json:"album_type`
	Artists              []AlbumArtist     `json:"artists"`
	ExternalUrls         map[string]string `json:"external_urls`
	ID                   string            `json:"id"`
	Image                []Image           `json:"images"`
	Name                 string            `json:"name"`
	ReleaseDate          string            `json:"release_date"`
	ReleaseDatePrecision string            `json:"release_date_precision"`
	TotalTracks          int               `json:"total_tracks"`
}

type getAlbumRangeResponse struct {
	Items []Album `json:"items"`
}

func get(URL *url.URL, bearerToken string, result interface{}) error {
	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprint("Bearer ", bearerToken))
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

func getAlbumRange(authToken string, artistID string, offset int, limit int) getAlbumRangeResponse {
	url, _ := url.Parse(fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums", artistID))

	query := url.Query()
	query.Set("offset", strconv.Itoa(offset))
	query.Set("limit", strconv.Itoa(limit))
	query.Set("include_groups", "album,single,compilation,appears_on")
	query.Set("market", "FR")
	url.RawQuery = query.Encode()

	var response getAlbumRangeResponse
	err := get(url, authToken, &response)

	if err != nil {
		panic(err)
	}

	return response
}

func getAlbums(authToken string, artistID string) chan Album {
	res := make(chan Album)

	return res
}
