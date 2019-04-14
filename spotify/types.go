package main

import "github.com/derdaele/maestro/internal/spotify"

type Composition struct {
	Name   string
	Tracks []spotify.Track
	Album  spotify.Album
}
