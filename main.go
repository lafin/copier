package main

import (
	"log"
	"os"

	api "github.com/lafin/vk"
)

func main() {
	clientID := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")

	log.Println("start")
	_, err := api.GetAccessToken(clientID, email, password)
	if err != nil {
		log.Fatalf("[main:api.GetAccessToken] error: %s", err)
		return
	}

	log.Println("done")
}
