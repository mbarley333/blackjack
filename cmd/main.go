package main

import (
	"card"
	"log"
)

func main() {
	err := card.NewBlackjackGame()
	if err != nil {
		log.Fatal(err)
	}
}
