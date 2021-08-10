package card_test

import (
	"card"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func TestCard(t *testing.T) {

// 	card := card.Card{
// 		Rank: card.Ace,
// 		Suit: card.Heart,
// 	}

// 	got := card.String()

// 	want := "Ace of Hearts"

// 	if want != got {
// 		t.Fatalf("want: %s, got: %s", want, got)
// 	}

// }

func TestDeal(t *testing.T) {

	random := rand.New(rand.NewSource(1))

	deck := card.NewDeck(random)

	got, err := deck.Deal()
	if err != nil {
		t.Fatal(err)
	}

	want := card.Card{
		Suit: card.Heart,
		Rank: card.Ace,
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}

}
