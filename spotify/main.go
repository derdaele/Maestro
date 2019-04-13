package main

import (
	"fmt"
	"os"

	"github.com/derdaele/maestro/internal/spotify"
)

func getAlbumTracks(client spotify.Client, album *spotify.Album, out chan spotify.Track) error {
	current := album.Tracks
	var err error

	for current != nil {
		if err != nil {
			return err
		}

		for _, track := range current.Items {
			out <- track
		}

		if current.Next != nil {
			err = client.Get(*current.Next, current)
		} else {
			current = nil
		}
	}

	return nil
}

func getArtistTracks(client spotify.Client, artistID string) chan spotify.Track {
	res := make(chan spotify.Track)

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
	tracks := getArtistTracks(client, beethoven)
	count := 0

	for range tracks {
		count++
	}

	fmt.Println("Total tracks = ", count)
}
