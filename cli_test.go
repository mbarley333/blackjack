package blackjack_test

import (
	"blackjack"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func TestNewBlackjackGame(t *testing.T) {
// 	t.Parallel()

// 	stack := []cards.Card{
// 		{Rank: cards.Ace, Suit: cards.Club},
// 		{Rank: cards.Eight, Suit: cards.Club},
// 		{Rank: cards.Jack, Suit: cards.Club},
// 		{Rank: cards.Seven, Suit: cards.Club},
// 		{Rank: cards.Ten, Suit: cards.Club},
// 		{Rank: cards.King, Suit: cards.Club},
// 	}

// 	deck := cards.Deck{
// 		Cards: stack,
// 	}

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("1\nPlanty\na\ns\n1")

// 	g, err := blackjack.NewBlackjackGame(
// 		blackjack.WithCustomDeck(deck),
// 		blackjack.WithOutput(output),
// 		blackjack.WithInput(input),
// 		blackjack.WithIncomingDeck(false),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	number := blackjack.GetNumberOfPlayers(output, input)
// 	g.AddNewPlayer(output, input, number)

// 	g.ResetPlayers()
// 	g.Start()

// 	want := blackjack.Record{
// 		Win:         1,
// 		HandsPlayed: 1,
// 	}
// 	got := g.Players[0].Record

// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}

// 	wantReport := "************** Player Win-Lose-Tie Report **************\nPlayer won: 1, lost: 0 and tied: 0\n"
// 	gotReport := g.Players[0].Record.RecordString()

// 	if wantReport != gotReport {
// 		t.Fatalf("want: %q, got: %q", wantReport, gotReport)
// 	}

// 	wantDealerScore := 15
// 	gotDealerScore := g.Dealer.Hands[0].Score()

// 	if wantDealerScore != gotDealerScore {
// 		t.Fatalf("want: %d, got: %d", wantDealerScore, gotDealerScore)
// 	}
// }

// func TestDisplayStageMessage(t *testing.T) {
// 	t.Parallel()

// 	g := blackjack.Game{}

// 	g.SetStage(blackjack.StageStart)

// 	output := &bytes.Buffer{}

// 	blackjack.DisplayStageMessage(output, "LET'S PLAY BLACKJACK!")

// 	want := "LET'S PLAY BLACKJACK!\n"
// 	got := output.String()

// 	if want != got {
// 		t.Fatalf("want: %q, got:%q", want, got)
// 	}

// }

// func TestPlayerMessage(t *testing.T) {
// 	t.Parallel()

// 	p := blackjack.Player{
// 		Message: "TEST MESSAGE",
// 	}

// 	output := &bytes.Buffer{}

// 	blackjack.DisplayMessage(output, p)

// 	want := "TEST MESSAGE\n"
// 	got := output.String()

// 	if want != got {
// 		t.Fatalf("want: %q, got: %q", want, got)
// 	}

// }

// // func TestCliDashboard(t *testing.T) {
// // 	t.Parallel()

// // 	output := &bytes.Buffer{}

// // 	g, err := blackjack.NewBlackjackGame(

// // 		blackjack.WithOutput(output),
// // 	)
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}

// // 	g.DisplayCliDashboard(output)

// // 	fmt.Print(output.String())
// // }

// // func TestDealMessage(t *testing.T) {
// // 	t.Parallel()

// // 	output := &bytes.Buffer{}
// // 	g, err := blackjack.NewBlackjackGame(
// // 		blackjack.WithOutput(output),
// // 	)
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}

// // 	p := &blackjack.Player{
// // 		Name:   "j",
// // 		Cash:   99,
// // 		Decide: blackjack.HumanAction,
// // 		Bet:    blackjack.HumanBet,
// // 		Hands: []*blackjack.Hand{
// // 			{
// // 				Cards: []cards.Card{
// // 					{Rank: cards.Queen, Suit: cards.Club},
// // 				},
// // 				Id:  1,
// // 				Bet: 1,
// // 			},
// // 		},
// // 	}

// // 	g.AddPlayer(p)

// // 	card := g.Deal(output)
// // 	g.Players[0].Hands[0].Cards = append(g.Players[0].Hands[0].Cards, card)
// // 	g.Players[0].Message = p.Name + " is dealt a " + card.String() + "\n"

// // }

// func TestAddBlackjackPlayers(t *testing.T) {
// 	t.Parallel()

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("Megatron")

// 	g, err := blackjack.NewBlackjackGame(
// 		blackjack.WithOutput(output),
// 		blackjack.WithInput(input),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	g.NumberAiPlayers = 1
// 	g.NumberHumanPlayers = 1

// 	g.AddBlackjackPlayers()

// 	want := 2
// 	got := len(g.Players)

// 	if want != got {
// 		t.Fatalf("wanted:%d, got:%d", want, got)
// 	}

// }

// func TestStartStage(t *testing.T) {
// 	t.Parallel()

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("")
// 	g, err := blackjack.NewBlackjackGame(
// 		blackjack.WithOutput(output),
// 		blackjack.WithInput(input),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	p := &blackjack.Player{
// 		Name:   "j",
// 		Cash:   100,
// 		Decide: blackjack.HumanAction,
// 		Bet:    blackjack.HumanBet,
// 		Hands: []*blackjack.Hand{
// 			{
// 				Id: 1,
// 			},
// 		},
// 	}

// 	g.AddPlayer(p)

// 	g.Start()

// 	want := "*************************\n* LET'S PLAY BLACKJACK! *\n*************************\n"
// 	got := output.String()

// 	if want != got {
// 		t.Fatalf("wanted: %q, got: %q", want, got)
// 	}

// }

// func TestPlayerWithAction(t *testing.T) {
// 	t.Parallel()

// 	g := &blackjack.Game{}

// 	player := &blackjack.Player{
// 		Name: "Thundercracker",
// 	}

// 	g.SetActivePlayer(player)

// 	want := &blackjack.Player{
// 		Name: "Thundercracker",
// 	}

// 	got := g.ActivePlayer

// 	if !cmp.Equal(want, got) {
// 		t.Errorf(cmp.Diff(want, got))
// 	}
// }

// func TestRenderGameCliStart(t *testing.T) {
// 	t.Parallel()

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("")

// 	g, err := blackjack.NewBlackjackGame(
// 		blackjack.WithOutput(output),
// 		blackjack.WithInput(input),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	g.SetStage(blackjack.StageStart)

// 	err = blackjack.RenderGameCli(output, input, g)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	want := "*************************\n* LET'S PLAY BLACKJACK! *\n*************************\n"

// 	got := output.String()

// 	if want != got {
// 		t.Fatalf("want: %q, got: %q", want, got)
// 	}

// }

// func TestRenderGameCliBetting(t *testing.T) {
// 	t.Parallel()

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("b\n10")

// 	g, err := blackjack.NewBlackjackGame(
// 		blackjack.WithOutput(output),
// 		blackjack.WithInput(input),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	p := &blackjack.Player{
// 		Name:   "j",
// 		Cash:   100,
// 		Decide: blackjack.HumanAction,
// 		Bet:    blackjack.HumanBet,
// 		Hands: []*blackjack.Hand{
// 			{
// 				Id: 1,
// 			},
// 		},
// 	}

// 	g.AddPlayer(p)

// 	g.SetStage(blackjack.StageBetting)

// 	g.Betting()

// 	want := p

// 	got := g.Players[0]

// 	if !cmp.Equal(want, got) {
// 		t.Errorf(cmp.Diff(want, got))
// 	}

// }

// func TestSetDialog(t *testing.T) {
// 	t.Parallel()

// 	p := blackjack.Player{
// 		Name: "Soundwave",
// 	}

// 	err := p.SetDialog(blackjack.DialogBetOrQuit)

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	want := blackjack.Player{
// 		Name:    "Soundwave",
// 		Dialog:  blackjack.DialogBetOrQuit,
// 		Message: "Soundwave would you like to (b)et or (q)uit? ",
// 	}

// 	got := p

// 	if !cmp.Equal(want, got) {
// 		t.Errorf(cmp.Diff(want, got))
// 	}

// }

func TestNewBlackjackGameWithArgs(t *testing.T) {
	t.Parallel()

	g, err := blackjack.NewBlackjackGameWithArgs(1, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	g.AddBlackjackPlayers()

}

func TestIsValid(t *testing.T) {
	t.Parallel()

	answer := "100"

	p := &blackjack.Player{
		Cash:   100,
		Dialog: blackjack.DialogPlaceYourBet,
	}

	ok, err := blackjack.IsInputValid(answer, p)
	if err != nil {
		t.Fatal(err)
	}

	want := true

	got := ok

	if want != got {
		t.Fatalf("want: %v, got: %v", want, got)
	}

}

func TestSetActionBet(t *testing.T) {
	t.Parallel()

	answer := "b"

	p := blackjack.Player{
		Dialog: blackjack.DialogBetOrQuit,
	}

	err := p.SetAction(answer)
	if err != nil {
		t.Fatal(err)
	}

	want := blackjack.Player{
		Dialog: blackjack.DialogBetOrQuit,
		Action: blackjack.ActionBet,
	}

	got := p

	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}

}

func TestSetActionQuit(t *testing.T) {
	t.Parallel()

	answer := "q"

	g := &blackjack.Game{}

	p := &blackjack.Player{
		Dialog: blackjack.DialogBetOrQuit,
	}

	g.AddPlayer(p)

	err := g.Players[0].SetAction(answer)
	if err != nil {
		t.Fatal(err)
	}

	//g.RemoveQuitPlayers()

	g.Players = g.RemoveQuitPlayers()

	want := 0

	got := len(g.Players)

	if want != got {
		t.Fatalf("want: %d, got: %d", want, got)
	}

	// want := &blackjack.Player{
	// 	Dialog: blackjack.DialogBetOrQuit,
	// 	Action: blackjack.ActionQuit,
	// }

	// got := g.Players[0]

	// if !cmp.Equal(want, got) {
	// 	t.Errorf(cmp.Diff(want, got))
	//}

}
