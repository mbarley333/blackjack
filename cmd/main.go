package main

import (
	"cards/blackjack"
	"log"
	"os"
)

func main() {
	g, err := blackjack.NewBlackjackGame()
	if err != nil {
		log.Fatal(err)
	}
	g.RunCLI(os.Stdout, os.Stdin)
}
