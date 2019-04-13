package main

import (
	"fmt"
	"os"

	"github.com/derdaele/maestro/internal/spotify"
)

type Entry struct {
	track *spotify.Track
	album *spotify.Album
}

func getAlbumTracks(client spotify.Client, album *spotify.Album, out chan Entry) error {
	current := album.Tracks
	var err error

	for current != nil {
		if err != nil {
			return err
		}

		for _, track := range current.Items {
			out <- Entry{track: &track, album: album}
		}

		if current.Next != nil {
			err = client.Get(*current.Next, current)
		} else {
			current = nil
		}
	}

	return nil
}

func getArtistTracks(client spotify.Client, artistID string) chan Entry {
	res := make(chan Entry)

	go func() {
		current, err := client.GetArtistAlbumRange(artistID, nil, nil)
		for current != nil {
			if err != nil {
				panic(err)
			}

			ids := make([]string, len(current.Items))
			for idx, album := range current.Items {
				ids[idx] = album.ID
			}

			// We make small batch of albums to get their tracks
			for batch := 0; batch < len(ids)/spotify.MaxRequestAlbumCount; batch++ {
				low, high := batch*spotify.MaxRequestAlbumCount, (batch+1)*spotify.MaxRequestAlbumCount
				albums, _ := client.GetAlbums(ids[low:high])

				for _, album := range albums.Albums {
					getAlbumTracks(client, &album, res)
				}
			}

			if current.Next != nil {
				err = client.Get(*current.Next, current)
			} else {
				current = nil
			}
		}

		close(res)
	}()

	return res
}

func main() {
	credentials := spotify.NewClientCredentials(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	client := spotify.NewClient(credentials, "FR")

	beethoven := "2wOqMjp9TyABvtHdOSOTUS"
	entries := getArtistTracks(client, beethoven)

	for entry := range entries {
		fmt.Println(entry.track.ID, entry.track.Name)
	}
}
