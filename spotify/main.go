package main

import "fmt"

const (
	clientID     = "6456ae94a67f49ba9307a3618415f44a"
	clientSecret = ""
)

func main() {
	token := getBearerToken(clientID, clientSecret)
	fmt.Println(token)
}
