package blackjack_test

import (
	"bytes"
	"cards"
	"cards/blackjack"
	"strings"
	"testing"
)

func TestCardCounting(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	input := strings.NewReader("b\n1")

	stack := []cards.Card{
		{Rank: cards.Queen, Suit: cards.Club},
	}

	deck := cards.Deck{
		Cards: stack,
	}

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
		blackjack.WithOutput(output),
		blackjack.WithInput(input),
		blackjack.WithIncomingDeck(false),
	)
	if err != nil {
		t.Fatal(err)
	}

	g.CountCards = blackjack.CountHiLo
	g.DeckCount = 1
	g.CardsDealt = 25

	c := blackjack.CardCounter{
		Count:     2,
		TrueCount: 0,
	}

	g.CardCounter = c

	_ = g.Deal(output)

	want := 1
	got := g.CardCounter.Count

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

	wantTrueCount := 2.0
	gotTrueCount := g.CardCounter.TrueCount

	if wantTrueCount != gotTrueCount {
		t.Fatalf("want: %v, got: %v", wantTrueCount, gotTrueCount)
	}

}
