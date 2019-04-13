package spotify

// Client provide an easy Spotify api access
type Client struct {
	Auth   AuthProvider
	Market string
}

// NewClient creates spotify client
func NewClient(auth AuthProvider, market string) Client {
	return Client{Auth: auth, Market: market}
}
