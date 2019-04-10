package main

import (
	"fmt"
	"path"

	"github.com/xfrr/goffmpeg/transcoder"
)

// TranscodeRequest
type transcodeRequest struct {
	inputFile  string
	outputFile string
}

func transcode(from string, to string) {
	fmt.Println("Transcoding", path.Base(from))

	trans := new(transcoder.Transcoder)

	err := trans.Initialize(from, to)

	if err != nil {
		fmt.Println("Count not transcode file", err)
	}

	done := trans.Run(false)
	err = <-done

	if err != nil {
		fmt.Println("Could not transcode file", err)
	}
}

// Transcoder start a transcoder
func startTranscoder(parallelism int) (chan transcodeRequest, chan error) {
	res := make(chan error)
	requests := make(chan transcodeRequest)

	go func() {
		sem := make(chan bool, parallelism)
		for request := range requests {
			sem <- true

			go func(request transcodeRequest) {
				defer func() { <-sem }()

				transcode(request.inputFile, request.outputFile)
			}(request)
		}

		for i := 0; i < cap(sem); i++ {
			sem <- true
		}

		res <- nil
	}()
	return requests, res
}
