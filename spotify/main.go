package main

import (
	"os"

	"github.com/derdaele/maestro/internal/spotify"
)

func main() {
	credentials := spotify.NewClientCredentials(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	client := spotify.NewClient(credentials, "FR")
	builder := NewDatabaseBuilder(&client)

	// beethoven := "2wOqMjp9TyABvtHdOSOTUS"
	// tchaikovsky := "3MKCzCnpzw3TjUYs2v7vDA"
	// chopin := "7y97mc3bZRFXzT2szRM4L4"
	mozart := "4NJhFmfw43RLBLjQvxDuRS"

	builder.BuildDatabase([]string{mozart})
}
