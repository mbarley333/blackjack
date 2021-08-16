package card_test

import (
	"card"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeal(t *testing.T) {
	t.Parallel()
	random := rand.New(rand.NewSource(1))

	deck := card.NewDeck(
		card.WithNumberOfDecks(3),
	)

	shuffle := deck.Shuffle(random)

	got, err := shuffle.Deal(1)
	if err != nil {
		t.Fatal(err)
	}

	want := []card.Card{
		{
			Suit: card.Diamond,
			Rank: card.Nine,
		},
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}

}

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
