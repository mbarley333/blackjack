package aceornothing_test

import (
	"cards"
	"cards/aceornothing"
	"testing"
)

func TestEvaluateAceOrNothing(t *testing.T) {

	t.Parallel()

	aceHand := []cards.Card{
		{
			Suit: cards.Spade,
			Rank: cards.Ace,
		},
	}
	result, err := aceornothing.EvaluateAceOrNothing(aceHand)
	if err != nil {
		t.Fatal(err)
	}

	if result != "Ace of Spades: WIN" {
		t.Fatalf("wanted: Ace of Spades: WIN, got:%s", result)
	}

	notAceHand := []cards.Card{
		{
			Suit: cards.Spade,
			Rank: cards.Jack,
		},
	}
	result, err = aceornothing.EvaluateAceOrNothing(notAceHand)
	if err != nil {
		t.Fatal(err)
	}

	if result != "Jack of Spades: LOSE" {
		t.Fatalf("wanted: LOSE, got:%s", result)
	}

}
