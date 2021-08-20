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
		card.WithNumberOfDecks(1),
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

	wantCardsRemaining := 51

	gotCardsRemaining := len(shuffle.Cards)

	if wantCardsRemaining != gotCardsRemaining {
		t.Fatalf("want: %d, got: %d", wantCardsRemaining, gotCardsRemaining)
	}

}
