package blackjack_test

import (
	"bytes"
	"cards"
	"cards/blackjack"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	input := strings.NewReader("1\nPlanty\na\n1")

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
		blackjack.WithOutput(output),
		blackjack.WithInput(input),
	)
	if err != nil {
		t.Fatal(err)
	}

	g.PlayerSetup(output, input)
	for g.Continue() {
		g.Betting()
		g.Players = g.RemoveQuitPlayers()
		if len(g.Players) == 0 {
			blackjack.RunCLI()
		}
		g.ResetPlayers()
		g.Start()
	}

	want := blackjack.Record{
		Win:         1,
		HandsPlayed: 1,
	}
	got := g.Players[0].Record

	if !cmp.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

	wantReport := "************** Player Win-Lose-Tie Report **************\nPlayer won: 1, lost: 0 and tied: 0\n"
	gotReport := g.Players[0].Record.String()

	if wantReport != gotReport {
		t.Fatalf("want: %q, got: %q", wantReport, gotReport)
	}

	wantDealerScore := 25
	gotDealerScore := g.Dealer.Score()

	if wantDealerScore != gotDealerScore {
		t.Fatalf("want: %d, got: %d", want, got)
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

	player := blackjack.Player{
		Name:   "Planty",
		Action: blackjack.ActionStand,
	}
	g.AddPlayer(player)

	player2 := blackjack.Player{
		Name:   "Kevin",
		Action: blackjack.ActionStand,
	}
	g.AddPlayer(player2)

	player3 := blackjack.Player{
		Name:   "Donald",
		Action: blackjack.ActionHit,
	}
	g.AddPlayer(player3)

	g.Start()

	output := &bytes.Buffer{}

	g.Outcome(output)

	wantPlayer := blackjack.OutcomeBlackjack
	gotPlayer := g.Players[0].HandOutcome

	if wantPlayer != gotPlayer {
		t.Fatalf("wanted: %q, got: %q", wantPlayer.String(), gotPlayer.String())
	}

	wantPlayer2 := blackjack.OutcomeWin
	gotPlayer2 := g.Players[1].HandOutcome

	if wantPlayer2 != gotPlayer2 {
		t.Fatalf("wanted: %q, got: %q", wantPlayer2.String(), gotPlayer2.String())
	}

	wantPlayer3 := blackjack.OutcomeBust
	gotPlayer3 := g.Players[2].HandOutcome

	if wantPlayer3 != gotPlayer3 {
		t.Fatalf("wanted: %q, got: %q", wantPlayer3.String(), gotPlayer3.String())
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

func TestHumanBet(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	input := strings.NewReader("q")

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithOutput(output),
		blackjack.WithInput(input),
	)
	if err != nil {
		t.Fatal(err)
	}

	player := blackjack.Player{
		Name: "test",
		Bet:  blackjack.HumanBet,
	}

	g.AddPlayer(player)

	g.Players[0] = g.Players[0].Bet(output, input, player)

	want := blackjack.ActionQuit.String()
	got := g.Players[0].Action.String()

	if want != got {
		t.Fatalf("wanted: %q, got: %q", want, got)
	}

}
