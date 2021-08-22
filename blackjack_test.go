package card_test

import (
	"card"
	"math/rand"
	"testing"
)

func TestBlackjack(t *testing.T) {

	t.Parallel()

	random := rand.New(rand.NewSource(1))

	deck := card.NewDeck(
		card.WithNumberOfDecks(1),
	)

	shuffle := deck.Shuffle(random)

	player := card.Hand{Player: "Player1"}
	var err error

	playerCards, err := shuffle.Deal(2)
	if err != nil {
		t.Fatal(err)
	}

	player.AddCards(playerCards)

	got := player.String()
	want := "Player1 has 14: [Nine of Diamonds][Five of Spades]"

	if want != got {
		t.Fatal(err)
	}

	wantPlayerScore := 14
	gotPlayerScore := player.Score

	if want != got {
		t.Fatalf("wanted: %d, got: %d", wantPlayerScore, gotPlayerScore)
	}

	dealer := card.Hand{Player: "Dealer"}

	dealerCards, err := shuffle.Deal(2)
	if err != nil {
		t.Fatal(err)
	}

	dealer.AddCards(dealerCards)

	want = "Dealer: [Three of Spades][???]"
	got = dealer.DealerString()

	if want != got {
		t.Fatal(err)
	}

	err = player.Hit(&shuffle)
	if err != nil {
		t.Fatal(err)
	}

	want = "Player1 has 23: [Nine of Diamonds][Five of Spades][Nine of Hearts]"
	got = player.String()

	if want != got {
		t.Fatalf("wanted: %s, got: %s", want, got)
	}

}

func TestBlackjackScoringFaces(t *testing.T) {
	card1 := card.Card{
		Rank: card.Nine,
	}

	card2 := card.Card{
		Rank: card.Jack,
	}

	player := card.Hand{}

	player.ScoreCards(card1)
	player.ScoreCards(card2)

	want := 19
	got := player.Score

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestBlackjackAces(t *testing.T) {
	card1 := []card.Card{
		{
			Rank: card.Ace,
		},
	}
	player := card.Hand{Cards: card1}

	// first card ace
	got := player.EvaluateAce()

	want := 0

	if want != got {
		t.Fatalf("wanted first card ace: %d, got: %d", want, got)
	}

	card2 := []card.Card{
		{
			Rank: card.Ten,
		},
	}

}
