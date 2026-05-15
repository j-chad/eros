package main

import (
	"backend/internal/crypto"
	"log"
)

func runGenerateVAPIDCmd() {
	private, public, err := crypto.GenerateVapidKeys()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("VAPID Private Key: %s", private)
	log.Printf("VAPID Public Key: %s", public)
}
