package blackjack_test

import (
	"bytes"
	"cards"
	"cards/blackjack"
	"strings"
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

	setupOutput := &bytes.Buffer{}
	setupInput := strings.NewReader("1\nPlanty\na\n3")

	g.RunCLI(setupOutput, setupInput)

	want := 3
	got := g.Players[0].Win + g.Players[0].Lose + g.Players[0].Tie

	if want != got {
		t.Fatalf("want: %d, got:%d", want, got)
	}

	// outputReport := &bytes.Buffer{}

	// g.Players[0].GetPlayerReport(outputReport)

	// wantReport := "************** Player Win-Lose-Tie Report **************\nPlayer won: 1, lost: 0 and tied: 0\n"
	// gotReport := outputReport.String()

	// if wantReport != gotReport {
	// 	t.Fatalf("want: %q, got: %q", wantReport, gotReport)
	// }
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

	setupOutput := &bytes.Buffer{}
	setupInput := strings.NewReader("1\nPlanty\na\n1")

	g.RunCLI(setupOutput, setupInput)

	got := blackjack.ReportMap[g.Players[0].HandOutcome]

	want := "***** Player wins! *****"

	if want != got {
		t.Fatalf("wanted: %s, got: %s", want, got)
	}

}

func TestMultiPlayers(t *testing.T) {
	stack := []cards.Card{
		{Rank: cards.Ace, Suit: cards.Club},
		{Rank: cards.Eight, Suit: cards.Club},
		{Rank: cards.Nine, Suit: cards.Club},
		{Rank: cards.Ten, Suit: cards.Spade},

		{Rank: cards.Jack, Suit: cards.Club},
		{Rank: cards.Ten, Suit: cards.Club},
		{Rank: cards.Six, Suit: cards.Club},
		{Rank: cards.Seven, Suit: cards.Spade},

		{Rank: cards.Seven, Suit: cards.Club},
		{Rank: cards.Four, Suit: cards.Club},
		{Rank: cards.Three, Suit: cards.Club},
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

	players := []blackjack.Player{
		{Name: "Planty"},
		{Name: "Kevin"},
		{Name: "Donald"},
	}

	g.Players = players

	g.Players[0].Action = blackjack.ActionStand
	g.Players[1].Action = blackjack.ActionStand
	g.Players[2].Action = blackjack.ActionHit

	g.Start()

	output := &bytes.Buffer{}

	g.Outcome(output)

	want := "BlackjackWinBust"                                                                                       //"Planty: ***** Blackjack!  Player wins *****\nKevin: ***** Player wins! *****\nDonald: ***** Bust!  Player loses *****\n"
	got := g.Players[0].HandOutcome.String() + g.Players[1].HandOutcome.String() + g.Players[2].HandOutcome.String() //output.String()

	if want != got {
		t.Fatalf("wanted:%q, got %q", want, got)
	}

}

func TestRemoveQuitPlayers(t *testing.T) {
	t.Parallel()

	g := blackjack.Game{
		Players: []blackjack.Player{
			{Action: blackjack.ActionStand},
			{Action: blackjack.ActionQuit},
			{Action: blackjack.ActionQuit},
		},
	}

	g.Players = g.RemoveQuitPlayers()

	want := 1

	got := len(g.Players)

	if want != got {
		t.Fatalf("wanted: %d, got: %d", want, got)
	}

}

// func TestDeal(t *testing.T) {
// 	t.Parallel()

// 	stack := []cards.Card{
// 		{Rank: cards.Ace, Suit: cards.Club},
// 		{Rank: cards.Seven, Suit: cards.Club},
// 		{Rank: cards.Jack, Suit: cards.Club},
// 		{Rank: cards.Queen, Suit: cards.Club},
// 	}

// 	deck := cards.Deck{
// 		Cards: stack,
// 	}

// 	g, err := blackjack.NewBlackjackGame(
// 		blackjack.WithCustomDeck(deck),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	player := blackjack.Player{
// 		Name: "Test Player",
// 	}

// 	g.Players = append(g.Players, player)

// 	g.OpeningDeal()

// 	want := 4
// 	got := len(g.Players[0].Hand) + len(g.Dealer.Hand)

// 	if want != got {
// 		t.Fatalf("wanted:%d, got:%d", want, got)
// 	}

// }

// func TestPlayerSetup(t *testing.T) {

// 	t.Parallel()
// 	g, err := blackjack.NewBlackjackGame()

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("1")

// 	want := "Please enter number of Blackjack players:\n"
// 	g.PlayerSetup(output, input)
// 	got := output.String()

// 	if want != got {
// 		t.Fatalf("wanted: %q, got:%q", want, got)
// 	}

// 	wantPlayers := 1
// 	gotPlayers := len(g.Players)
// 	if wantPlayers != gotPlayers {
// 		t.Fatalf("wanted: %d, got:%d", wantPlayers, gotPlayers)
// 	}

// }

// func TestNewPlayer(t *testing.T) {
// 	t.Parallel()

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("1\ntest\nA\n10")

// 	g, err := blackjack.NewBlackjackGame(
// 		blackjack.WithOutput(output),
// 		blackjack.WithInput(input),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	g.RunCLI(g.Output, g.Input)

// 	want := 10
// 	got := g.Players[0].Lose + g.Players[0].Tie + g.Players[0].Win

// 	if want != got {
// 		t.Fatalf("wanted: %d, got: %d", want, got)
// 	}
// }

// func TestGameStage(t *testing.T) {
// 	t.Parallel()

// 	g, err := blackjack.NewBlackjackGame()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	player := blackjack.Player{
// 		Name: "Test",
// 	}

// 	g.Players = append(g.Players, player)

// 	want := "Invalid Stage"

// 	got := g.Stage.String()

// 	if want != got {
// 		t.Fatalf("wanted: %s, got:%s", want, got)
// 	}

// 	g.Betting()

// 	wantBet := "Betting"

// 	gotBet := g.Stage.String()

// 	if wantBet != gotBet {
// 		t.Fatalf("wanted: %s, got:%s", want, got)
// 	}

// }
