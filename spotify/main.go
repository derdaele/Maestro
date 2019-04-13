package main

import (
	"fmt"
	"os"

	"github.com/derdaele/maestro/internal/spotify"
)

func main() {
	credentials := spotify.NewClientCredentials(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	client := spotify.NewClient(credentials, "FR")

	beethoven := "2wOqMjp9TyABvtHdOSOTUS"
	res, _ := client.GetArtistAlbumRange(beethoven, 0, 50)

	ids := make([]string, 50)
	for idx, album := range res.Items {
		ids[idx] = album.ID
	}

	albums, _ := client.GetAlbums(ids[0:20])
	for _, album := range albums.Albums {
		if album.Tracks == nil {
			continue
		}

		for _, track := range album.Tracks.Items {
			fmt.Println(track.ID, track.Name)
		}
	}
}
