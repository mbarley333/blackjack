package blackjack_test

import (
	"cards"
	"cards/blackjack"
	"testing"
)

func TestNewBlackjackGame(t *testing.T) {
	t.Parallel()

	// what do we need to start a blackjack game game
	// we need a game server layer.  NewGame func return *Game
	// we need to deal cards from a deck.  Deal func
	// we need to control the cards dealt in order to test.  Pass in stacked deck
	// we need a player and a dealer.  add player and dealer to Game
	// add deck to Game
	// deal cards to player and dealer
	// dealer logic for hit/stand
	// player prompt for hit/stand

	stack := []cards.Card{
		{Rank: cards.Ace, Suit: cards.Club},
		{Rank: cards.Eight, Suit: cards.Club},
		{Rank: cards.Jack, Suit: cards.Club},
		{Rank: cards.Seven, Suit: cards.Club},
		{Rank: cards.Ten, Suit: cards.Club},
		{Rank: cards.King, Suit: cards.Club},
	}

	deck := cards.Deck{
		Cards: stack,
	}

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
	)
	if err != nil {
		t.Fatal(err)
	}

	g.Start()

	g.DealerStart()

	outcome := g.Outcome()

	got := blackjack.ReportWinner(outcome)

	want := "***** Blackjack!  Player wins *****"

	if want != got {
		t.Fatalf("wanted: %s, got: %s", want, got)
	}

}

func TestDealerBust(t *testing.T) {
	t.Parallel()

	stack := []cards.Card{
		{Rank: cards.Ace, Suit: cards.Club},
		{Rank: cards.Eight, Suit: cards.Club},
		{Rank: cards.Nine, Suit: cards.Club},
		{Rank: cards.Seven, Suit: cards.Club},
		{Rank: cards.Ten, Suit: cards.Club},
		{Rank: cards.King, Suit: cards.Club},
	}

	deck := cards.Deck{
		Cards: stack,
	}

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
	)
	if err != nil {
		t.Fatal(err)
	}

	g.Player.Action = blackjack.Stand

	g.Start()

	g.DealerStart()

	outcome := g.Outcome()

	got := blackjack.ReportWinner(outcome)

	want := "***** Player wins! *****"

	if want != got {
		t.Fatalf("wanted: %s, got: %s", want, got)
	}

}

func TestPlayerBust(t *testing.T) {
	t.Parallel()

	stack := []cards.Card{
		{Rank: cards.Queen, Suit: cards.Club},
		{Rank: cards.Eight, Suit: cards.Club},
		{Rank: cards.Six, Suit: cards.Club},
		{Rank: cards.Seven, Suit: cards.Club},
		{Rank: cards.Ten, Suit: cards.Club},
		{Rank: cards.King, Suit: cards.Club},
	}

	deck := cards.Deck{
		Cards: stack,
	}

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
	)
	if err != nil {
		t.Fatal(err)
	}

	g.Player.Action = blackjack.Hit
	g.Start()

	outcome := g.Outcome()

	got := blackjack.ReportWinner(outcome)

	want := "***** Bust!  Player loses *****"

	if want != got {
		t.Fatalf("wanted: %s, got: %s", want, got)
	}

}
