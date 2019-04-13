package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/derdaele/maestro/internal/spotify"
)

// Entry is a track with it owning album
type Entry struct {
	track spotify.Track
	album spotify.Album
}

func containsArtist(track *spotify.Track, artistID string) bool {
	for _, artist := range track.Artists {
		if artist.ID == artistID {
			return true
		}
	}
	return false
}

func getAlbumTracks(client *spotify.Client, artistID string, album *spotify.Album, out chan Entry) error {
	current := album.Tracks
	var err error

	for current != nil {
		if err != nil {
			return err
		}

		for _, track := range current.Items {
			if !containsArtist(&track, artistID) {
				continue
			}

			out <- Entry{track: track, album: *album}
		}

		if current.Next != nil {
			err = client.Get(*current.Next, current)
		} else {
			current = nil
		}
	}

	return nil
}

func getArtistTracks(client *spotify.Client, artistID string) chan Entry {
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
				if high > len(ids) {
					high = len(ids)
				}
				albums, _ := client.GetAlbums(ids[low:high])

				for _, album := range albums.Albums {
					getAlbumTracks(client, artistID, &album, res)
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

type Composition struct {
	Name   string
	Tracks []spotify.Track
	Album  spotify.Album
}

func getCompositions(client *spotify.Client, artistID string) chan Composition {
	res := make(chan Composition)

	go func() {
		entries := getArtistTracks(client, artistID)
		tracks := make([]spotify.Track, 0)
		var previousComposition string
		var previousAlbum spotify.Album

		for entry := range entries {
			composition := strings.SplitN(entry.track.Name, ":", 2)[0]
			if previousComposition != "" {
				if previousComposition != composition {
					res <- Composition{
						Album:  previousAlbum,
						Name:   previousComposition,
						Tracks: tracks,
					}

					tracks = make([]spotify.Track, 0)
				}
			}

			previousComposition = composition
			previousAlbum = entry.album
			tracks = append(tracks, entry.track)
		}

		res <- Composition{
			Album:  previousAlbum,
			Name:   previousComposition,
			Tracks: tracks,
		}

		close(res)
	}()

	return res
}

type CompoStats struct {
	composition string
	count       int
}

type byComposition []*CompoStats

func (s byComposition) Len() int {
	return len(s)
}

func (s byComposition) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byComposition) Less(i, j int) bool {
	return s[i].composition < s[j].composition
}

func main() {
	credentials := spotify.NewClientCredentials(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	client := spotify.NewClient(credentials, "FR")

	// beethoven := "2wOqMjp9TyABvtHdOSOTUS"
	// tchaikovsky := "3MKCzCnpzw3TjUYs2v7vDA"
	// chopin := "7y97mc3bZRFXzT2szRM4L4"
	mozart := "4NJhFmfw43RLBLjQvxDuRS"

	compositionsFreq := make(map[string]int)
	compositions := getCompositions(&client, mozart)

	for compo := range compositions {
		compositionsFreq[compo.Name] = compositionsFreq[compo.Name] + 1
	}

	candidates := make([]*CompoStats, 0)
	for compo, count := range compositionsFreq {
		if count > 5 {
			candidates = append(candidates, &CompoStats{composition: compo, count: count})
		}
	}

	sort.Sort(byComposition(candidates))
	for _, stats := range candidates {
		fmt.Println(stats.composition, "[", stats.count, "]")
	}
}
