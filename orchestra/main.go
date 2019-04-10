package main

import (
	"fmt"
	"os"

	"github.com/dhowden/tag"

	"github.com/derdaele/maestro/internal/filesystem"
)

func getTags(filepath string) (tag.Metadata, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	m, err := tag.ReadFrom(file)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "<classical music directory>")
		os.Exit(1)
	}

	for musicFilePath := range filesystem.ListMusicFiles(os.Args[1]) {
		m, err := getTags(musicFilePath)

		if err != nil {
			fmt.Println("Error", err)
			continue
		}

		fmt.Println(m.Title())
	}
}
