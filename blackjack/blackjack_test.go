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

	got := blackjack.ReportMap[g.Player.HandOutcome]

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

	g.Player.Action = blackjack.ActionStand

	g.Start()

	//g.DealerStart()

	outcome := g.Outcome()

	got := blackjack.ReportMap[outcome]

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

	g.Player.Action = blackjack.ActionHit
	g.Start()

	got := blackjack.ReportMap[g.Player.HandOutcome]

	want := "***** Bust!  Player loses *****"

	if want != got {
		t.Fatalf("wanted: %s, got: %s", want, got)
	}

}

func TestOnlyStandAI(t *testing.T) {
	// reuse the blackjack machinery
	// insert the ai logic into the flow of game

	stack := []cards.Card{
		{Rank: cards.Ace, Suit: cards.Club},
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
		blackjack.WithAiType(blackjack.AiStandOnly),
	)

	if err != nil {
		t.Fatal(err)
	}

	g.Start()

	got := g.Player.Action
	want := blackjack.ActionStand

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}
