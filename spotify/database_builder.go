package main

import (
	"fmt"
	"strings"

	"github.com/derdaele/maestro/internal/spotify"
)

type DatabaseBuilder struct {
	client *spotify.Client
}

func NewDatabaseBuilder(client *spotify.Client) *DatabaseBuilder {
	return &DatabaseBuilder{client: client}
}

func (builder *DatabaseBuilder) BuildDatabase(artistIDs []string) {
	entries := make(chan Entry)
	compositions := make(chan Composition)

	for _, artistID := range artistIDs {
		go builder.FetchAllArtistTracks(artistID, entries)
		go builder.GroupEntriesByComposition(entries, compositions)
		builder.WriteCompositionsToDB(compositions)
	}
}

func (builder *DatabaseBuilder) WriteCompositionsToDB(compositions chan Composition) {
	for compo := range compositions {
		fmt.Println(compo.Name)
	}
}

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

func (builder *DatabaseBuilder) FetchAllAlbumTracks(artistID string, album *spotify.Album, out chan Entry) error {
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
			err = builder.client.Get(*current.Next, current)
		} else {
			current = nil
		}
	}

	return nil
}

func (builder *DatabaseBuilder) FetchAllArtistTracks(artistID string, res chan Entry) {
	current, err := builder.client.GetArtistAlbumRange(artistID, nil, nil)
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
			albums, _ := builder.client.GetAlbums(ids[low:high])

			for _, album := range albums.Albums {
				builder.FetchAllAlbumTracks(artistID, &album, res)
			}
		}

		if current.Next != nil {
			err = builder.client.Get(*current.Next, current)
		} else {
			current = nil
		}
	}

	close(res)
}

func (builder *DatabaseBuilder) GroupEntriesByComposition(entries chan Entry, res chan Composition) {
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
}
