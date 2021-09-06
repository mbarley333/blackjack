package blackjack_test

import (
	"cards"
	"cards/blackjack"
	"fmt"
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

	got := blackjack.ReportMap[g.Player.HandOutcome]

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
	t.Parallel()

	stack := []cards.Card{
		{Rank: cards.Ace, Suit: cards.Club},
		{Rank: cards.Seven, Suit: cards.Club},
		{Rank: cards.Jack, Suit: cards.Club},
		{Rank: cards.Queen, Suit: cards.Club},

		{Rank: cards.Ten, Suit: cards.Heart},
		{Rank: cards.King, Suit: cards.Heart},
		{Rank: cards.Two, Suit: cards.Heart},
		{Rank: cards.Queen, Suit: cards.Heart},

		{Rank: cards.Ten, Suit: cards.Heart},
		{Rank: cards.King, Suit: cards.Heart},
		{Rank: cards.Jack, Suit: cards.Heart},
		{Rank: cards.Queen, Suit: cards.Heart},
	}

	deck := cards.Deck{
		Cards: stack,
	}

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
		blackjack.WithAiType(blackjack.AiStandOnly),
		blackjack.WithAiHandsToPlay(3),
	)

	if err != nil {
		t.Fatal(err)
	}

	g.RunCLI()

	want := blackjack.ActionQuit
	got := g.Player.Action

	if want != got {
		fmt.Printf("want: %d, got:%d", want, got)
	}

	wantHandsPlayed := 3
	gotHandsPlayed := g.HandsPlayed

	if wantHandsPlayed != gotHandsPlayed {
		t.Fatalf("want: %d, got:%d", wantHandsPlayed, gotHandsPlayed)
	}

	wantPlayed := 3
	gotPlayed := g.PlayerWin + g.PlayerLose + g.PlayerTie

	fmt.Println(gotPlayed)
	if wantPlayed != gotPlayed {
		t.Fatalf("want: %d, got:%d", wantPlayed, gotPlayed)
	}

	wantReport := "Player won: 1, lost: 1 and tied: 1"
	gotReport := g.GetPlayerReport()
	fmt.Println(wantReport)

	if wantReport != gotReport {
		t.Fatalf("want: %s, got:%s", wantReport, gotReport)
	}

}
