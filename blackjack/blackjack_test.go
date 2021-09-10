package blackjack_test

import (
	"bytes"
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

	output := &bytes.Buffer{}

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
		blackjack.WithOutput(output),
	)
	if err != nil {
		t.Fatal(err)
	}
	player := blackjack.Player{
		Name: "Planty",
	}

	g.Players = append(g.Players, player)

	g.OpeningDeal()

	g.ShowPlayerCards(output)

	wantCards := "Planty has 21: [Ace of Clubs][Jack of Clubs]\n"
	gotCarts := output.String()

	if wantCards != gotCarts {
		t.Fatalf("wanted: %q, got: %q", wantCards, gotCarts)
	}

	g.Players[0].HandOutcome = blackjack.OutcomeBlackjack

	wantContinue := false
	gotContinue := g.Players[0].PlayerContinue()

	if wantContinue != gotContinue {
		t.Fatalf("wanted: %v, got: %v", wantContinue, gotContinue)
	}

	g.DealerPlay()

	outputOutcome := &bytes.Buffer{}

	g.Outcome(outputOutcome)

	wantOutcome := "Planty: ***** Blackjack!  Player wins *****\n"
	gotOutcome := outputOutcome.String()

	if wantOutcome != gotOutcome {
		t.Fatalf("wanted: %q, got: %q", wantOutcome, gotOutcome)
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

	player := blackjack.Player{
		Name: "Planty",
	}

	g.Players = append(g.Players, player)

	g.Players[0].Action = blackjack.ActionStand

	g.Start()

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

		{Rank: cards.Jack, Suit: cards.Club},
		{Rank: cards.Ten, Suit: cards.Club},
		{Rank: cards.Six, Suit: cards.Club},

		{Rank: cards.Seven, Suit: cards.Club},
		{Rank: cards.Four, Suit: cards.Club},
		{Rank: cards.Three, Suit: cards.Club},
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
	}

	g.Players = players

	g.Players[0].Action = blackjack.ActionStand
	g.Players[1].Action = blackjack.ActionStand

	g.Start()

	output := &bytes.Buffer{}

	g.Outcome(output)

	want := "Planty: ***** Blackjack!  Player wins *****\nKevin: ***** Player wins! *****\n"
	got := output.String()

	if want != got {
		t.Fatalf("wanted:%q, got %q", want, got)
	}

}

func TestResetPlayers(t *testing.T) {
	t.Parallel()

	g := blackjack.Game{
		Players: []blackjack.Player{
			{Action: blackjack.ActionStand},
			{Action: blackjack.ActionQuit},
			{Action: blackjack.ActionQuit},
		},
	}

	g.ResetPlayers()

	want := 1

	got := len(g.Players)

	if want != got {
		t.Fatalf("wanted: %d, got: %d", want, got)
	}

}

func TestDeal(t *testing.T) {
	t.Parallel()

	stack := []cards.Card{
		{Rank: cards.Ace, Suit: cards.Club},
		{Rank: cards.Seven, Suit: cards.Club},
		{Rank: cards.Jack, Suit: cards.Club},
		{Rank: cards.Queen, Suit: cards.Club},
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

	player := blackjack.Player{
		Name: "Test Player",
	}

	g.Players = append(g.Players, player)

	g.OpeningDeal()

	want := 4
	got := len(g.Players[0].Hand) + len(g.Dealer.Hand)

	if want != got {
		t.Fatalf("wanted:%d, got:%d", want, got)
	}

}

// func TestPlayerBust(t *testing.T) {
// 	t.Parallel()

// 	stack := []cards.Card{
// 		{Rank: cards.Queen, Suit: cards.Club},
// 		{Rank: cards.Eight, Suit: cards.Club},
// 		{Rank: cards.Six, Suit: cards.Club},
// 		{Rank: cards.Seven, Suit: cards.Club},
// 		{Rank: cards.Ten, Suit: cards.Club},
// 		{Rank: cards.King, Suit: cards.Club},
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

// 	g.Player.Action = blackjack.ActionHit
// 	g.Start()

// 	got := blackjack.ReportMap[g.Player.HandOutcome]

// 	want := "***** Bust!  Player loses *****"

// 	if want != got {
// 		t.Fatalf("wanted: %s, got: %s", want, got)
// 	}

// }

// func TestOnlyStandAI(t *testing.T) {
// 	t.Parallel()

// 	stack := []cards.Card{
// 		{Rank: cards.Ace, Suit: cards.Club},
// 		{Rank: cards.Seven, Suit: cards.Club},
// 		{Rank: cards.Jack, Suit: cards.Club},
// 		{Rank: cards.Queen, Suit: cards.Club},

// 		{Rank: cards.Ten, Suit: cards.Heart},
// 		{Rank: cards.King, Suit: cards.Heart},
// 		{Rank: cards.Two, Suit: cards.Heart},
// 		{Rank: cards.Queen, Suit: cards.Heart},

// 		{Rank: cards.Ten, Suit: cards.Heart},
// 		{Rank: cards.King, Suit: cards.Heart},
// 		{Rank: cards.Jack, Suit: cards.Heart},
// 		{Rank: cards.Queen, Suit: cards.Heart},
// 	}

// 	deck := cards.Deck{
// 		Cards: stack,
// 	}

// 	g, err := blackjack.NewBlackjackGame(
// 		blackjack.WithCustomDeck(deck),
// 		blackjack.WithAiType(blackjack.AiStandOnly),
// 		blackjack.WithAiHandsToPlay(2),
// 	)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	g.RunCLI()

// 	want := blackjack.ActionQuit
// 	got := g.Player.Action

// 	if want != got {
// 		fmt.Printf("want: %d, got:%d", want, got)
// 	}

// 	wantHandsPlayed := 3
// 	gotHandsPlayed := g.HandsPlayed

// 	if wantHandsPlayed != gotHandsPlayed {
// 		t.Fatalf("want: %d, got:%d", wantHandsPlayed, gotHandsPlayed)
// 	}

// 	wantPlayed := 3
// 	gotPlayed := g.PlayerWin + g.PlayerLose + g.PlayerTie

// 	fmt.Println(gotPlayed)
// 	if wantPlayed != gotPlayed {
// 		t.Fatalf("want: %d, got:%d", wantPlayed, gotPlayed)
// 	}

// 	wantReport := "Player won: 1, lost: 1 and tied: 1"
// 	gotReport := g.GetPlayerReport()
// 	fmt.Println(wantReport)

// 	if wantReport != gotReport {
// 		t.Fatalf("want: %s, got:%s", wantReport, gotReport)
// 	}

// }

// func TestGetPlayerAction(t *testing.T) {
// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("h")
// 	wantPrompt := "Please choose (H)it, (S)tand or (Q)uit\n"

// 	got := blackjack.GetPlayerAction(output, input)
// 	gotPrompt := output.String()

// 	want := blackjack.ActionHit

// 	if want != got {
// 		t.Fatalf("wanted: %v, got: %v", want, got)
// 	}

// 	if wantPrompt != gotPrompt {
// 		t.Fatalf("wanted: %q, got: %q", wantPrompt, gotPrompt)
// 	}

// }

// func TestPlayerSetup(t *testing.T) {

// 	t.Parallel()
// 	g, err := blackjack.NewBlackjackGame()

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("2")

// 	want := "Please enter number of Blackjack players:\n"
// 	g.PlayerSetup(output, input)
// 	got := output.String()

// 	if want != got {
// 		t.Fatalf("wanted: %q, got:%q", want, got)
// 	}

// 	wantPlayers := 2
// 	gotPlayers := len(g.Players)
// 	if wantPlayers != gotPlayers {
// 		t.Fatalf("wanted: %d, got:%d", wantPlayers, gotPlayers)
// 	}

// }
