package main

import (
	"fmt"
	"os"

	"github.com/derdaele/maestro/internal/spotify"
)

func main() {
	credentials := spotify.NewClientCredentials(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	client := spotify.NewClient(credentials)

	beethoven := "2wOqMjp9TyABvtHdOSOTUS"
	res, _ := client.GetArtistAlbumRange(beethoven, 50, 50)

	for _, album := range res.Items {
		fmt.Println(album.Name, album.TotalTracks)
	}
}
