package spotify

// Client provide an easy Spotify api access
type Client struct {
	Auth AuthProvider
}

// NewClient creates spotify client
func NewClient(auth AuthProvider) Client {
	return Client{Auth: auth}
}
