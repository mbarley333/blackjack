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
	input := strings.NewReader("b\n1")

	stack := []cards.Card{
		{Rank: cards.Queen, Suit: cards.Club},
		{Rank: cards.Three, Suit: cards.Club},
		{Rank: cards.Five, Suit: cards.Club},
		{Rank: cards.Jack, Suit: cards.Club},
		{Rank: cards.Six, Suit: cards.Club},
		{Rank: cards.Eight, Suit: cards.Club},
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

	id := p.NextHandId()
	hand := blackjack.NewHand(id)
	p.AddHand(hand)

	g.AddPlayer(p)

	g.Betting()

	g.OpeningDeal()

	g.Players = g.RemoveQuitPlayers()

	card := g.Deal(output)

	g.Players[0].Hands[0].Hit(output, card, g.Players[0].Name)

	want := &blackjack.Player{
		Name:   "j",
		Cash:   99,
		Decide: blackjack.HumanAction,
		Bet:    blackjack.HumanBet,
		Hands: []*blackjack.Hand{
			{
				Cards: []cards.Card{
					{Rank: cards.Queen, Suit: cards.Club},
					{Rank: cards.Five, Suit: cards.Club},
					{Rank: cards.Six, Suit: cards.Club},
				},
				Id:  1,
				Bet: 1,
			},
		},
	}

	got := g.Players[0]

	if !cmp.Equal(want, got, cmpopts.IgnoreFields(blackjack.Player{}, "Decide", "Bet")) {
		t.Error(cmp.Diff(want, got))
	}

	g.DealerPlay()

	wantDealer := &blackjack.Player{
		Hands: []*blackjack.Hand{
			{
				Cards: []cards.Card{
					{Rank: cards.Three, Suit: cards.Club},
					{Rank: cards.Jack, Suit: cards.Club},
					{Rank: cards.Eight, Suit: cards.Club},
				},
				Id:     1,
				Action: blackjack.ActionStand,
			},
		},
	}

	gotDealer := g.Dealer

	if !cmp.Equal(wantDealer, gotDealer, cmpopts.IgnoreFields(blackjack.Player{}, "Decide", "Bet")) {
		t.Error(cmp.Diff(wantDealer, gotDealer))
	}

	g.Outcome(output)

	wantOutcome := &blackjack.Player{
		Name:   "j",
		Cash:   100,
		Decide: blackjack.HumanAction,
		Bet:    blackjack.HumanBet,
		Hands: []*blackjack.Hand{
			{
				Cards: []cards.Card{
					{Rank: cards.Queen, Suit: cards.Club},
					{Rank: cards.Five, Suit: cards.Club},
					{Rank: cards.Six, Suit: cards.Club},
				},
				Id:      1,
				Bet:     0,
				Outcome: blackjack.OutcomeTie,
			},
		},
		Record: blackjack.Record{
			Tie:         1,
			HandsPlayed: 1,
		},
	}

	gotOutcome := g.Players[0]

	if !cmp.Equal(wantOutcome, gotOutcome, cmpopts.IgnoreFields(blackjack.Player{}, "Decide", "Bet")) {
		t.Error(cmp.Diff(wantOutcome, gotOutcome))
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
		t.Error(cmp.Diff(want, got))
	}

	wantReport := "************** Player Win-Lose-Tie Report **************\nPlayer won: 1, lost: 0 and tied: 0\n"
	gotReport := g.Players[0].Record.RecordString()

	if wantReport != gotReport {
		t.Fatalf("want: %q, got: %q", wantReport, gotReport)
	}

	wantDealerScore := 15
	gotDealerScore := g.Dealer.Hands[0].Score()

	if wantDealerScore != gotDealerScore {
		t.Fatalf("want: %d, got: %d", wantDealerScore, gotDealerScore)
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

	g.DeckCount = 1

	player := &blackjack.Player{
		Name:   "Planty",
		Action: blackjack.None,
		Cash:   99,
		Hands: []*blackjack.Hand{
			{
				Id:     1,
				Bet:    1,
				Action: blackjack.ActionStand,
			},
		},
	}
	g.AddPlayer(player)

	player2 := &blackjack.Player{
		Name:   "Kevin",
		Action: blackjack.None,
		Cash:   99,
		Hands: []*blackjack.Hand{
			{
				Id:     1,
				Bet:    1,
				Action: blackjack.ActionStand,
			},
		},
	}
	g.AddPlayer(player2)

	player3 := &blackjack.Player{
		Name:   "Donald",
		Action: blackjack.None,
		Cash:   99,
		Hands: []*blackjack.Hand{
			{
				Id:     1,
				Bet:    1,
				Action: blackjack.ActionHit,
			},
		},
	}
	g.AddPlayer(player3)

	g.Start()

	wantPlayer := blackjack.OutcomeBlackjack
	gotPlayer := g.Players[0].Hands[0].Outcome

	if wantPlayer != gotPlayer {
		t.Fatalf("wanted: %q, got: %q", wantPlayer.String(), gotPlayer.String())
	}

	wantPlayer2 := blackjack.OutcomeWin
	gotPlayer2 := g.Players[1].Hands[0].Outcome

	if wantPlayer2 != gotPlayer2 {
		t.Fatalf("wanted: %q, got: %q", wantPlayer2.String(), gotPlayer2.String())
	}

	wantPlayer3 := blackjack.OutcomeBust
	gotPlayer3 := g.Players[2].Hands[0].Outcome

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

		want := &blackjack.Player{
			Cash: tc.cash,
			Hands: []*blackjack.Hand{
				{
					Id:      1,
					Outcome: tc.outcome,
					Payout:  tc.handPayout,
				},
			},
		}

		p := &blackjack.Player{
			Cash: 99,
			Hands: []*blackjack.Hand{
				{
					Id:      1,
					Bet:     tc.bet,
					Outcome: tc.outcome,
				},
			},
		}

		p.Payout()
		got := p

		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
		}
	}
}

func TestPlayerCash(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}

	p := &blackjack.Player{
		Name: "Planty",
		Cash: 101,
		Hands: []*blackjack.Hand{
			{
				Id:      1,
				Outcome: blackjack.OutcomeWin,
				Payout:  1,
			},
		},
	}

	p.OutcomeReport(output)

	want := "Planty won $1.  Cash available: $101\n"

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
		Name: "Planty",
		Hands: []*blackjack.Hand{
			{
				Id:      1,
				Outcome: blackjack.OutcomeLose,
				Payout:  -1,
			},
		},
		Cash: 0,
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

func TestResetFieldsAfterIncomingDeck(t *testing.T) {
	t.Parallel()

	g := blackjack.Game{
		CardsDealt: 55,
		CardCounter: blackjack.CardCounter{
			Count:     7,
			TrueCount: 3.0,
		},
	}
	g.ResetFieldsAfterIncomingDeck()

	want := blackjack.CardCounter{
		Count:     0,
		TrueCount: 0,
	}

	got := g.CardCounter

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	wantCardsDealt := 0
	gotCardsDealt := g.CardsDealt

	if wantCardsDealt != gotCardsDealt {
		t.Fatalf("want: %d, got: %d", wantCardsDealt, gotCardsDealt)
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
		dealerHand  []*blackjack.Hand
		description string
		result      bool
	}
	tcs := []testCase{
		{
			players: []*blackjack.Player{
				{Hands: []*blackjack.Hand{
					{
						Outcome: blackjack.OutcomeBust,
					},
				},
				},
			},
			dealerHand: []*blackjack.Hand{
				{
					Cards: []cards.Card{{Rank: cards.Seven, Suit: cards.Club}, {Rank: cards.Seven, Suit: cards.Club}},
				},
			},
			result:      false,
			description: "All Players Busted",
		},
		{
			players: []*blackjack.Player{
				{Hands: []*blackjack.Hand{
					{
						Outcome: blackjack.OutcomeBust,
					},
				},
				},
				{Hands: []*blackjack.Hand{
					{},
				},
				},
			},
			dealerHand: []*blackjack.Hand{
				{
					Cards: []cards.Card{{Rank: cards.Seven, Suit: cards.Club}, {Rank: cards.Seven, Suit: cards.Club}},
				},
			},
			result:      true,
			description: "All Players Not Busted",
		},
		{
			players: []*blackjack.Player{
				{Hands: []*blackjack.Hand{
					{
						Outcome: blackjack.OutcomeBust,
					},
				},
				},
				{Hands: []*blackjack.Hand{
					{
						Outcome: blackjack.OutcomeBlackjack,
					},
				},
				},
			},
			dealerHand: []*blackjack.Hand{
				{
					Cards: []cards.Card{{Rank: cards.Seven, Suit: cards.Club}, {Rank: cards.Seven, Suit: cards.Club}},
				},
			},
			result:      false,
			description: "All Players Blackjack or Bust",
		},
	}

	for _, tc := range tcs {

		g.Players = tc.players
		g.Dealer.Hands = tc.dealerHand
		want := tc.result
		got := g.IsDealerDraw()

		if want != got {
			t.Fatalf("%s: wanted: %v, got: %v", tc.description, want, got)
		}
		g.Players = nil
		g.Dealer.Hands = nil

	}

}

func TestDoubleDown(t *testing.T) {

	stack := []cards.Card{
		{Rank: cards.Six, Suit: cards.Club},
		{Rank: cards.Four, Suit: cards.Club},
		{Rank: cards.Four, Suit: cards.Club},
		{Rank: cards.Jack, Suit: cards.Club},
		{Rank: cards.Ace, Suit: cards.Club},
		{Rank: cards.Ten, Suit: cards.Club},
	}

	deck := cards.Deck{
		Cards: stack,
	}
	output := &bytes.Buffer{}
	input := strings.NewReader("d\nq")
	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
		blackjack.WithIncomingDeck(false),
		blackjack.WithOutput(output),
		blackjack.WithInput(input),
	)
	if err != nil {
		t.Fatal(err)
	}

	p := &blackjack.Player{
		Name: "planty",
		Hands: []*blackjack.Hand{
			{
				Id:  1,
				Bet: 1,
			},
		},
		Cash:   99,
		Decide: blackjack.HumanAction,
		Bet:    blackjack.HumanBet,
	}

	g.AddPlayer(p)

	g.Start()

	want := 102

	got := g.Players[0].Cash

	if want != got {
		t.Fatalf("wanted: %d, got: %d", want, got)
	}

}

func TestSplit(t *testing.T) {
	t.Parallel()
	output := &bytes.Buffer{}

	stack := []cards.Card{
		{Rank: cards.Six, Suit: cards.Heart},
		{Rank: cards.Six, Suit: cards.Club},
		{Rank: cards.Nine, Suit: cards.Spade},
		{Rank: cards.Four, Suit: cards.Diamond},
	}

	deck := cards.Deck{
		Cards: stack,
	}

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithCustomDeck(deck),
		blackjack.WithOutput(output),
		blackjack.WithIncomingDeck(false),
	)
	if err != nil {
		t.Fatal(err)
	}

	p := &blackjack.Player{
		Cash: 99,
		Hands: []*blackjack.Hand{
			{
				Id:  1,
				Bet: 1,
			},
		},
	}

	g.AddPlayer(p)

	card := g.Deal(output)
	g.Players[0].Hands[p.HandIndex].Cards = append(g.Players[0].Hands[p.HandIndex].Cards, card)
	card = g.Deal(output)
	g.Players[0].Hands[p.HandIndex].Cards = append(g.Players[0].Hands[p.HandIndex].Cards, card)

	card1 := g.Deal(output)
	card2 := g.Deal(output)

	g.Players[0].Split(output, card1, card2)

	want := &blackjack.Player{
		Cash: 98,
		Hands: []*blackjack.Hand{
			{
				Id: 1,
				Cards: []cards.Card{
					{Rank: cards.Six, Suit: cards.Heart},
					{Rank: cards.Nine, Suit: cards.Spade},
				},
				Bet:    1,
				Action: blackjack.None,
			},
			{
				Id: 2,
				Cards: []cards.Card{
					{Rank: cards.Six, Suit: cards.Club},
					{Rank: cards.Four, Suit: cards.Diamond},
				},
				Bet:    1,
				Action: blackjack.None,
			},
		},
	}
	got := g.Players[0]

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestPlayHand(t *testing.T) {
	output := &bytes.Buffer{}
	input := strings.NewReader("s\ns")

	stack := []cards.Card{
		{Rank: cards.Six, Suit: cards.Heart},
		{Rank: cards.Six, Suit: cards.Club},
		{Rank: cards.Nine, Suit: cards.Spade},
		{Rank: cards.Four, Suit: cards.Diamond},
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
		Name:   "test",
		Cash:   99,
		Decide: blackjack.HumanAction,
		Hands: []*blackjack.Hand{
			{
				Id:     1,
				Bet:    1,
				Action: blackjack.ActionSplit,
				Cards: []cards.Card{
					{Rank: cards.Three, Suit: cards.Club},
					{Rank: cards.Three, Suit: cards.Club},
				},
			},
		},
	}

	g.AddPlayer(p)

	g.Dealer = &blackjack.Player{
		Hands: []*blackjack.Hand{
			{
				Cards: []cards.Card{
					{Rank: cards.Three, Suit: cards.Club},
					{Rank: cards.Three, Suit: cards.Club},
				},
				Id:     1,
				Action: blackjack.ActionStand,
			},
		},
	}

	g.PlayHand(g.Players[0])

	want := &blackjack.Player{
		Name: "test",
		Cash: 98,
		Hands: []*blackjack.Hand{
			{
				Cards: []cards.Card{
					{Rank: cards.Three, Suit: cards.Club},
					{Rank: cards.Six, Suit: cards.Heart},
				},
				Id:     1,
				Action: blackjack.ActionStand,
				Bet:    1,
			},
			{
				Cards: []cards.Card{
					{Rank: cards.Three, Suit: cards.Club},
					{Rank: cards.Six, Suit: cards.Club},
				},
				Id:     2,
				Action: blackjack.ActionStand,
				Bet:    1,
			},
		},
	}

	got := g.Players[0]

	if !cmp.Equal(want, got, cmpopts.IgnoreFields(blackjack.Player{}, "Decide", "Bet")) {
		t.Error(cmp.Diff(want, got))
	}

}
