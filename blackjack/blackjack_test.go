package blackjack_test

import (
	"bytes"
	"cards"
	"cards/blackjack"
	"math/rand"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// what do we need to start a blackjack game game
// we need a game server layer.  NewGame func return *Game
// we need to deal cards from a deck.  Deal func
// we need to control the cards dealt in order to test.  Pass in stacked deck
// we need a player and a dealer.  add player and dealer to Game
// add deck to Game
// deal cards to player and dealer
// dealer logic for hit/stand
// player prompt for hit/stand

func TestBlackjack(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	input := strings.NewReader("b\n1\ns\nq")

	stack := []cards.Card{
		{Rank: cards.Queen, Suit: cards.Club},
		{Rank: cards.Three, Suit: cards.Club},
		{Rank: cards.Ten, Suit: cards.Club},
		{Rank: cards.Jack, Suit: cards.Club},
		{Rank: cards.King, Suit: cards.Club},
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

	p := &blackjack.Player{
		Name:   "j",
		Cash:   100,
		Decide: blackjack.HumanAction,
		Bet:    blackjack.HumanBet,
	}

	g.AddPlayer(p)

	g.ResetPlayers()

	g.Betting()

	g.Players = g.RemoveQuitPlayers()

	g.Start()

	want := 101

	got := g.Players[0].Cash

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

}

func TestNewBlackjackGame(t *testing.T) {
	t.Parallel()

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
	input := strings.NewReader("1\nPlanty\na\ns\n1")

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
		blackjack.WithOutput(output),
		blackjack.WithInput(input),
		blackjack.WithIncomingDeck(false),
	)
	if err != nil {
		t.Fatal(err)
	}

	g.PlayerSetup(output, input)

	g.ResetPlayers()
	g.Start()

	want := blackjack.Record{
		Win:         1,
		HandsPlayed: 1,
	}
	got := g.Players[0].Record

	if !cmp.Equal(want, got) {
		cmp.Diff(want, got)
	}

	wantReport := "************** Player Win-Lose-Tie Report **************\nPlayer won: 1, lost: 0 and tied: 0\n"
	gotReport := g.Players[0].Record.RecordString()

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
		blackjack.WithIncomingDeck(false),
	)
	if err != nil {
		t.Fatal(err)
	}

	player := blackjack.Player{
		Name:    "Planty",
		Action:  blackjack.ActionStand,
		Cash:    99,
		HandBet: 1,
	}
	g.AddPlayer(&player)

	player2 := blackjack.Player{
		Name:    "Kevin",
		Action:  blackjack.ActionStand,
		Cash:    99,
		HandBet: 1,
	}
	g.AddPlayer(&player2)

	player3 := blackjack.Player{
		Name:    "Donald",
		Action:  blackjack.ActionHit,
		Cash:    99,
		HandBet: 1,
	}
	g.AddPlayer(&player3)

	g.Start()

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
		Players: []*blackjack.Player{
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

func TestBetting(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	input := strings.NewReader("b\n1")

	player := blackjack.Player{
		Bet:         blackjack.HumanBet,
		Cash:        100,
		HandOutcome: blackjack.OutcomeWin,
	}

	err := player.Bet(output, input, &player)
	if err != nil {
		t.Fatal(err)
	}

	player.Payout()

	want := blackjack.Player{
		HandBet: 0,
		Cash:    101,
	}

	got := player

	if !cmp.Equal(want, got) {
		cmp.Diff(want, got)
	}

}

func TestPayout(t *testing.T) {
	t.Parallel()

	type testCase struct {
		bet        int
		cash       int
		outcome    blackjack.Outcome
		handPayout int
	}
	tcs := []testCase{
		{bet: 1, cash: 101, outcome: blackjack.OutcomeWin, handPayout: 1},
		{bet: 1, cash: 99, outcome: blackjack.OutcomeLose, handPayout: -1},
		{bet: 1, cash: 100, outcome: blackjack.OutcomeTie, handPayout: 0},
		{bet: 1, cash: 102, outcome: blackjack.OutcomeBlackjack, handPayout: 2},
	}

	for _, tc := range tcs {
		want := blackjack.Player{
			HandBet:     tc.bet,
			Cash:        tc.cash,
			HandOutcome: tc.outcome,
			HandPayout:  tc.handPayout,
		}

		p := blackjack.Player{
			Cash:        100,
			HandBet:     tc.bet,
			HandOutcome: tc.outcome,
		}

		p.Payout()
		got := p

		if !cmp.Equal(want, got) {
			cmp.Diff(want, got)
		}
	}
}

func TestPlayerCash(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}

	p := blackjack.Player{
		Name:        "Planty",
		Cash:        100,
		HandOutcome: blackjack.OutcomeWin,
		HandPayout:  1,
	}

	p.OutcomeReport(output)

	want := "Planty won $1.  Cash available: $100\n"

	got := output.String()

	if want != got {
		t.Fatalf("wanted: %q, got: %q", want, got)
	}

}

func TestPlayerBroke(t *testing.T) {
	t.Parallel()

	g, err := blackjack.NewBlackjackGame()

	if err != nil {
		t.Fatal(err)
	}

	p := &blackjack.Player{
		Name:        "Planty",
		HandBet:     1,
		Cash:        0,
		HandOutcome: blackjack.OutcomeLose,
	}

	g.AddPlayer(p)

	p.Broke()

	want := blackjack.ActionQuit

	got := p.Action

	if want != got {
		t.Fatalf("want: %q, got: %q", want.String(), got.String())
	}

}

func TestIncomingDeck(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	random := rand.New(rand.NewSource(1))

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithOutput(output),
		blackjack.WithDeckCount(3),
		blackjack.WithRandom(random),
	)

	if err != nil {
		t.Fatal(err)
	}

	want := g.Shoe

	got := g.IncomingDeck()

	if cmp.Equal(want, got, cmpopts.IgnoreUnexported(cards.Deck{})) {
		t.Fatal("wanted a new deck, got old deck")
	}

}

func TestAiBasicAction(t *testing.T) {

	g, err := blackjack.NewBlackjackGame()
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		playerCards []cards.Card
		dealerCard  cards.Card
		action      blackjack.Action
		description string
	}
	tcs := []testCase{
		{
			playerCards: []cards.Card{{Rank: cards.Ace, Suit: cards.Club}, {Rank: cards.Jack, Suit: cards.Club}},
			dealerCard:  cards.Card{Rank: cards.Four, Suit: cards.Club},
			action:      blackjack.ActionStand,
			description: "Blackjack",
		},
		{
			playerCards: []cards.Card{{Rank: cards.Five, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}},
			dealerCard:  cards.Card{Rank: cards.Four, Suit: cards.Club},
			action:      blackjack.ActionHit,
			description: "Eleven or less",
		},
		{
			playerCards: []cards.Card{{Rank: cards.Five, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}},
			dealerCard:  cards.Card{Rank: cards.Four, Suit: cards.Club},
			action:      blackjack.ActionHit,
			description: "Soft 15 or less",
		},
		{
			playerCards: []cards.Card{{Rank: cards.Queen, Suit: cards.Club}, {Rank: cards.Nine, Suit: cards.Club}},
			dealerCard:  cards.Card{Rank: cards.Four, Suit: cards.Club},
			action:      blackjack.ActionStand,
			description: "Soft 19 or higher",
		},
		{
			playerCards: []cards.Card{{Rank: cards.Ten, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}, {Rank: cards.Five, Suit: cards.Club}},
			dealerCard:  cards.Card{Rank: cards.Four, Suit: cards.Club},
			action:      blackjack.ActionStand,
			description: "Hard 17 to 21",
		},
		{
			playerCards: []cards.Card{{Rank: cards.Ten, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}, {Rank: cards.Two, Suit: cards.Club}},
			dealerCard:  cards.Card{Rank: cards.Four, Suit: cards.Club},
			action:      blackjack.ActionStand,
			description: "Hard 12 to 16 w/ Dealer <= 6",
		},
		{
			playerCards: []cards.Card{{Rank: cards.Six, Suit: cards.Club}, {Rank: cards.Two, Suit: cards.Club}, {Rank: cards.Four, Suit: cards.Club}},
			dealerCard:  cards.Card{Rank: cards.Three, Suit: cards.Club},
			action:      blackjack.ActionHit,
			description: "Hard 12 w/ Dealer <= 3",
		},
		{
			playerCards: []cards.Card{{Rank: cards.Ten, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}, {Rank: cards.Two, Suit: cards.Club}},
			dealerCard:  cards.Card{Rank: cards.Seven, Suit: cards.Club},
			action:      blackjack.ActionHit,
			description: "Hard 12 to 16 w/ Dealer >= 7",
		},
	}

	output := &bytes.Buffer{}
	input := strings.NewReader("")

	for _, tc := range tcs {
		p := &blackjack.Player{
			Hand:   tc.playerCards,
			Decide: blackjack.AiActionBasic,
		}
		g.AddPlayer(p)
		g.Dealer.Hand = append(g.Dealer.Hand, tc.dealerCard)

		want := tc.action

		got := g.Players[0].Decide(output, input, g.Players[0], g.Dealer.Hand[0])

		if want != got {
			t.Fatalf("%q: wanted: %q, got: %q", tc.description, want.String(), got.String())
		}
		g.Players = []*blackjack.Player{}
		g.Dealer.Hand = nil

	}

}

func TestScoreDealerHoleCard(t *testing.T) {

	type testCase struct {
		card        cards.Card
		score       int
		description string
	}
	tcs := []testCase{
		{card: cards.Card{Rank: cards.Ace, Suit: cards.Club}, score: 11, description: "Ace"},
		{card: cards.Card{Rank: cards.King, Suit: cards.Club}, score: 10, description: "King"},
		{card: cards.Card{Rank: cards.Three, Suit: cards.Club}, score: 3, description: "Three"},
	}

	for _, tc := range tcs {
		want := tc.score
		got := blackjack.ScoreDealerHoleCard(tc.card)

		if want != got {
			t.Fatalf("wanted: %d, got: %d", want, got)
		}

	}

}

func TestDealerAi(t *testing.T) {

	g, err := blackjack.NewBlackjackGame()
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		players     []*blackjack.Player
		dealerHand  []cards.Card
		description string
		result      bool
	}
	tcs := []testCase{
		{
			players:     []*blackjack.Player{{HandOutcome: blackjack.OutcomeBust}},
			dealerHand:  []cards.Card{{Rank: cards.Seven, Suit: cards.Club}, {Rank: cards.Seven, Suit: cards.Club}},
			result:      false,
			description: "All Players Busted",
		},
		{
			players:     []*blackjack.Player{{HandOutcome: blackjack.OutcomeNone}},
			dealerHand:  []cards.Card{{Rank: cards.Seven, Suit: cards.Club}, {Rank: cards.Seven, Suit: cards.Club}},
			result:      true,
			description: "All Players Not Busted",
		},
	}

	for _, tc := range tcs {

		g.Players = tc.players
		g.Dealer.Hand = tc.dealerHand
		want := tc.result
		got := g.IsDealerDraw()

		if want != got {
			t.Fatalf("%s: wanted: %v, got: %v", tc.description, want, got)
		}
		g.Players = nil
		g.Dealer.Hand = nil

	}

}
