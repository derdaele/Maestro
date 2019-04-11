package main

import (
	"fmt"
	"net/http"
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

func get(URL string, bearerToken string) (*http.Response, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprint("Bearer ", bearerToken))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func getAlbumRange(authToken string, artistID string, offset int, limit int) getAlbumRangeResponse {
	res, err := get(fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums?offset=%d&limit=%d&include_groups=album,single,compilation,appears_on&market=FR",
	artistID, offset, limit))

	if err != nil {
		panic(err)
	}

	var response getAlbumRangeResponse

	payload, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(payload, &response)

	if err != nil {
		panic(err)
	}

	return response
}

func getAlbums(authToken string, artistID string) chan Album {
	res := make(chan Album)

	get(fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums?offset=1&limit=1&include_groups=album,single,compilation,appears_on&market=FR")
}
