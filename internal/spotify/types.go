package spotify

// Image represent an image url with its size.
type Image struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

// AlbumArtist represent an artist/composer/...
type AlbumArtist struct {
	ExternalUrls map[string]string `json:"external_urls"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
}

// Album is a collection of tracks.
type Album struct {
	AlbumGroup           string            `json:"album_group"`
	AlbumType            string            `json:"album_type"`
	Artists              []AlbumArtist     `json:"artists"`
	ExternalUrls         map[string]string `json:"external_urls"`
	ID                   string            `json:"id"`
	Image                []Image           `json:"images"`
	Name                 string            `json:"name"`
	ReleaseDate          string            `json:"release_date"`
	ReleaseDatePrecision string            `json:"release_date_precision"`
	TotalTracks          int               `json:"total_tracks"`
}
