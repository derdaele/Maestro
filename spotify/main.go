package main

import (
	"fmt"
	"os"
)

func main() {
	token := getBearerToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	albumRange := getAlbumRange(token, "2wOqMjp9TyABvtHdOSOTUS", 0, 50)

	for _, album := range albumRange.Items {
		fmt.Println(album.Name, album.TotalTracks)
	}
}
