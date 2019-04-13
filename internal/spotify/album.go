package spotify

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// GetArtistAlbumRangeResponse represent a range of album for a given artist
type GetArtistAlbumRangeResponse struct {
	Href     string  `json:"href"`
	Items    []Album `json:"items"`
	Next     string  `json:"next"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
}

// GetArtistAlbumRange retrieve a range of albums for a given artist
func (c Client) GetArtistAlbumRange(artistID string, offset int, limit int) (*GetArtistAlbumRangeResponse, error) {
	url, _ := url.Parse(fmt.Sprintf("%s/artists/%s/albums", baseEndpoint, artistID))
	query := url.Query()
	query.Set("offset", strconv.Itoa(offset))
	query.Set("limit", strconv.Itoa(limit))
	query.Set("include_groups", "album,single,compilation,appears_on")
	query.Set("market", c.Market)
	url.RawQuery = query.Encode()

	var response GetArtistAlbumRangeResponse
	err := c.Get(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetAlbums retrieve album data, tracks included, for a batch of album ids
func (c Client) GetAlbums(ids []string) (*GetAlbums, error) {
	if len(ids) > maxRequestAlbumCount {
		return nil, errors.New(fmt.Sprintf("Too many album provided: %d (maximum is %d)", len(ids), maxRequestAlbumCount))
	}

	url, _ := url.Parse(fmt.Sprintf("%s/albums", baseEndpoint))
	query := url.Query()
	query.Set("market", c.Market)
	query.Set("ids", strings.Join(ids, ","))
	url.RawQuery = query.Encode()

	var response GetAlbums
	err := c.Get(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type GetAlbums struct {
	Albums []Album `json:"albums"`
}
