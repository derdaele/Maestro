package spotify

// Image represent an image url with its size.
type Image struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

// Artist represent an artist/composer/...
type Artist struct {
	ExternalUrls map[string]string `json:"external_urls"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Href         string            `json:"href"`
}

type GetAlbumTracksRange struct {
	Href     string  `json:"href"`
	Items    []Track `json:"items"`
	Next     string  `json:"next"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
}

// Album is a collection of tracks.
type Album struct {
	AlbumGroup           string               `json:"album_group"`
	AlbumType            string               `json:"album_type"`
	Artists              []Artist             `json:"artists"`
	ExternalUrls         map[string]string    `json:"external_urls"`
	ID                   string               `json:"id"`
	Image                []Image              `json:"images"`
	Name                 string               `json:"name"`
	ReleaseDate          string               `json:"release_date"`
	ReleaseDatePrecision string               `json:"release_date_precision"`
	TotalTracks          int                  `json:"total_tracks"`
	Tracks               *GetAlbumTracksRange `json:"tracks"`
}

// Track represent a music file
type Track struct {
	Artists      []Artist          `json:"artists"`
	DiscNumber   int               `json:"disc_number"`
	DurationMs   int               `json:"duration_ms"`
	Explicit     bool              `json:"explicit"`
	ExternalUrls map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	ID           string            `json:"id"`
	IsLocal      bool              `json:"is_local"`
	IsPlayable   bool              `json:"is_playable"`
	Name         string            `json:"name"`
	PreviewURL   string            `json:"preview_url"`
	TrackNumber  int               `json:"track_number"`
}
