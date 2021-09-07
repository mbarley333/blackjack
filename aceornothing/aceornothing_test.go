package aceornothing_test

import (
	"cards"
	"cards/aceornothing"
	"testing"
)

func TestEvaluateAceOrNothing(t *testing.T) {

	t.Parallel()

	deck := cards.Deck{
		Cards: []cards.Card{
			{Suit: cards.Spade, Rank: cards.Ace},
			{Suit: cards.Spade, Rank: cards.Five},
		},
	}

	g := aceornothing.NewGame(
		aceornothing.WithCustomDeck(deck),
	)

	g.Hand.Deal(&g.Shoe)

	result, err := g.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	if result != "Ace of Spades: WIN" {
		t.Fatalf("wanted: Ace of Spades: WIN, got:%s", result)
	}

	g.Hand = aceornothing.Hand{}

	g.Hand.Deal(&g.Shoe)

	result, err = g.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	if result != "Five of Spades: LOSE" {
		t.Fatalf("wanted: LOSE, got:%s", result)
	}

}
