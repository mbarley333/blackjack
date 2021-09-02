package main

import (
	"cards/blackjack"
	"log"
)

func main() {
	g, err := blackjack.NewBlackjackGame()
	if err != nil {
		log.Fatal(err)
	}
	g.RunCLI()
}
