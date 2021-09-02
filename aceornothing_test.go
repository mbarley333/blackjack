package card_test

import (
	"card"
	"testing"
)

func TestEvaluateAceOrNothing(t *testing.T) {

	t.Parallel()

	aceHand := []card.Card{
		{
			Suit: card.Spade,
			Rank: card.Ace,
		},
	}
	result, err := card.EvaluateAceOrNothing(aceHand)
	if err != nil {
		t.Fatal(err)
	}

	if result != "Ace of Spades: WIN" {
		t.Fatalf("wanted: Ace of Spades: WIN, got:%s", result)
	}

	notAceHand := []card.Card{
		{
			Suit: card.Spade,
			Rank: card.Jack,
		},
	}
	result, err = card.EvaluateAceOrNothing(notAceHand)
	if err != nil {
		t.Fatal(err)
	}

	if result != "Jack of Spades: LOSE" {
		t.Fatalf("wanted: LOSE, got:%s", result)
	}

}
