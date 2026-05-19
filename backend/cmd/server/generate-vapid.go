package main

import (
	"backend/internal/webpush"
	"log"
)

func runGenerateVAPIDCmd() {
	private, public, err := webpush.GenerateVapidKeys()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("VAPID Private Key: %s", private)
	log.Printf("VAPID Public Key: %s", public)
}
