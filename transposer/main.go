package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

func isMusicFile(info os.FileInfo) bool {
	switch path.Ext(info.Name()) {
	case ".wma":
		return true
	default:
		return false
	}
}

func transcode(from string, to string) {
	trans := new(transcoder.Transcoder)

	err := trans.Initialize(from, to)

	if err != nil {
		fmt.Println("Count not transcode file", err)
	}

	trans.MediaFile().SetAudioBitRate(320)
	done := trans.Run(false)
	err = <-done

	if err != nil {
		fmt.Println("Could not transcode file", err)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "<input directory> <output directory>")
		os.Exit(1)
	}

	inputDirectory := os.Args[1]
	outputDirectory := os.Args[2]

	filepath.Walk(inputDirectory, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !isMusicFile(info) {
			return nil
		}

		relativeFilepath := strings.TrimPrefix(filepath, inputDirectory)
		relativeFilepath = strings.TrimSuffix(relativeFilepath, path.Ext(relativeFilepath))
		relativeFilepath = relativeFilepath + ".mp3"
		outputFilepath := path.Join(outputDirectory, relativeFilepath)
		dir := path.Dir(outputFilepath)
		os.MkdirAll(dir, os.ModePerm)

		fmt.Println("Transcoding from", filepath, "to", outputFilepath)
		transcode(filepath, outputFilepath)

		return nil
	})
}
