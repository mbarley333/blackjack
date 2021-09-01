package cards_test

import (
	"cards"
	"math/rand"
	"testing"
)

func TestShuffle(t *testing.T) {
	t.Parallel()
	random := rand.New(rand.NewSource(1))

	deck := cards.NewDeck(
		cards.WithNumberOfDecks(1),
	)

	shuffle := deck.Shuffle(random)

	got := shuffle.Cards[0].String()

	want := "Nine of Diamonds"

	if want != got {
		t.Fatalf("want: %s, got: %s", want, got)
	}

}
