package card_test

import (
	"card"
	"math/rand"
	"strconv"
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

	player.Cards, err = shuffle.Deal(2)
	if err != nil {
		t.Fatal(err)
	}

	got := player.String()
	want := "Player1 has 14: [Nine of Diamonds][Five of Spades]"

	if want != got {
		t.Fatal(err)
	}

	wantPlayerScore := 14
	gotPlayerScore := player.Score()

	if want != got {
		t.Fatalf("wanted: %d, got: %d", wantPlayerScore, gotPlayerScore)
	}

	dealer := card.Hand{Player: "Dealer"}

	dealer.Cards, err = shuffle.Deal(2)
	if err != nil {
		t.Fatal(err)
	}

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
	t.Parallel()

	cards := []card.Card{
		{
			Rank: card.Nine,
		},
		{
			Rank: card.Jack,
		},
	}

	player := card.Hand{Cards: cards}

	want := 19
	got := player.Score()

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestBlackjackScoring21(t *testing.T) {
	t.Parallel()

	cards := []card.Card{
		{
			Rank: card.Ace,
		},
		{
			Rank: card.Jack,
		},
	}

	player := card.Hand{Cards: cards}

	want := 21
	got := player.Score()

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestWinner(t *testing.T) {
	t.Parallel()

	playerCards := []card.Card{
		{
			Rank: card.Ace,
		},
		{
			Rank: card.Jack,
		},
	}

	player := card.Hand{
		Cards: playerCards,
	}

	dealerCards := []card.Card{
		{
			Rank: card.Eight,
		},
		{
			Rank: card.Jack,
		},
	}

	dealer := card.Hand{
		Cards: dealerCards,
	}

	got := card.IsPlayerWinner(player.Score(), dealer.Score())
	want := true

	if want != got {
		t.Fatalf("want: %s, got:%s", strconv.FormatBool(want), strconv.FormatBool(got))
	}

}
