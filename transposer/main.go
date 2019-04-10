package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/derdaele/maestro/internal/filesystem"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "<input directory> <output directory>")
		os.Exit(1)
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]
	transcode, errChan := startTranscoder(4)

	for inputFile := range filesystem.ListMusicFiles(inputDir) {
		var outputFile string

		// Relative path
		outputFile = strings.TrimPrefix(inputFile, inputDir)

		// We change the file extension to mp3
		outputFile = strings.TrimSuffix(outputFile, path.Ext(inputFile)) + ".mp3"

		// We set the absolute path relative to the output directory
		outputFile = path.Join(outputDir, outputFile)

		os.MkdirAll(path.Dir(outputFile), os.ModePerm)

		transcode <- transcodeRequest{inputFile: inputFile, outputFile: outputFile}
	}
	close(transcode)

	err := <-errChan

	if err != nil {
		fmt.Println("Error", err)
	}
}
