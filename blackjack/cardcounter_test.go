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

	// need a card
	// need to assign a value based on count
	card := g.Deal(output)
	g.CountCards = blackjack.CountHiLo
	g.CardCounter.Count = 2
	g.DeckCount = 1
	g.CardsDealt = 26

	g.CardCounter.Count, g.CardCounter.TrueCount = g.CountCards(card, g.CardCounter.Count, g.CardsDealt, g.DeckCount)

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
