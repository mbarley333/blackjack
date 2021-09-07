package main

import (
	"cards/blackjack"
	"log"
)

func main() {
	g, err := blackjack.NewBlackjackGame(
		blackjack.WithAiType(blackjack.AiStandOnly),
		blackjack.WithAiHandsToPlay(100),
	)
	if err != nil {
		log.Fatal(err)
	}
	g.RunCLI()
}
