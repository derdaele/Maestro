package spotify

import (
	"fmt"
	"net/url"
	"strconv"
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
	query.Set("market", "FR")
	url.RawQuery = query.Encode()

	var response GetArtistAlbumRangeResponse
	err := c.get(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
