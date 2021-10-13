package blackjack_test

import (
	"blackjack"
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewBlackjackGameWithArgs(t *testing.T) {
	t.Parallel()

	g, err := blackjack.NewBlackjackGameWithArgs(1, 0, 6)
	if err != nil {
		t.Fatal(err)
	}

	want := 1
	got := g.NumberHumanPlayers

	if want != got {
		t.Fatalf("want: %d human player, got %d human player", want, got)
	}

	wantAiPlayers := 0
	gotAiPlayers := g.NumberAiPlayers

	if wantAiPlayers != gotAiPlayers {
		t.Fatalf("want: %d ai player, got %d ai player", wantAiPlayers, gotAiPlayers)
	}

	wantDecks := 6
	gotDecks := g.DeckCount

	if wantDecks != gotDecks {
		t.Fatalf("want: %d decks, got %d decks", wantDecks, gotDecks)
	}

}

func TestSetActivePlayer(t *testing.T) {
	t.Parallel()

	g := &blackjack.Game{}

	player := &blackjack.Player{
		Name: "Thundercracker",
	}

	g.SetActivePlayer(player)

	want := &blackjack.Player{
		Name: "Thundercracker",
	}

	got := g.ActivePlayer

	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}

func TestStageBetting(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	input := strings.NewReader("b\n10")

	g, err := blackjack.NewBlackjackGame(
		blackjack.WithOutput(output),
		blackjack.WithInput(input),
	)
	if err != nil {
		t.Fatal(err)
	}

	p := &blackjack.Player{
		Name:   "j",
		Cash:   100,
		Decide: blackjack.HumanAction,
		Bet:    blackjack.HumanBet,
		Hands: []*blackjack.Hand{
			{
				Id: 1,
			},
		},
	}

	g.AddPlayer(p)

	g.SetStage(blackjack.StageBetting)

	g.Betting()

	want := p

	got := g.Players[0]

	if !cmp.Equal(want, got, cmpopts.IgnoreFields(blackjack.Player{}, "Bet", "Decide")) {
		t.Errorf(cmp.Diff(want, got))
	}

}

func TestSetDialog(t *testing.T) {
	t.Parallel()

	p := blackjack.Player{
		Name: "Soundwave",
	}

	err := p.SetDialog(blackjack.DialogBetOrQuit)

	if err != nil {
		t.Fatal(err)
	}
	want := blackjack.Player{
		Name:   "Soundwave",
		Dialog: blackjack.DialogBetOrQuit,
	}

	got := p

	if !cmp.Equal(want, got, cmpopts.IgnoreFields(blackjack.Player{}, "Message")) {
		t.Errorf(cmp.Diff(want, got))
	}

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

func TestNewAiPlayer(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	input := strings.NewReader("aaa\nb\n10")

	index := 0

	want := &blackjack.Player{
		Name:           "aaa",
		Cash:           100,
		AiRoundsToPlay: 10,
		Hands: []*blackjack.Hand{
			{Id: 1},
		},
	}

	got := blackjack.NewAiPlayer(output, input, index)

	if !cmp.Equal(want, got, cmpopts.IgnoreFields(blackjack.Player{}, "Bet", "Decide")) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestSpacesInNameHumanPlayer(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	input := strings.NewReader("James Bond")

	index := 0

	human := blackjack.NewHumanPlayer(output, input, index)

	want := "James Bond"

	got := human.Name

	if want != got {
		t.Fatalf("want: %q, got: %q", want, got)
	}

}

func TestSpacesInNameAiPlayer(t *testing.T) {
	t.Parallel()

	//output := &bytes.Buffer{}
	input := strings.NewReader("James Bond\nb\n10")

	reader := bufio.NewReader(input)
	name, _ := reader.ReadString('\n')
	//reader.ReadString()
	//name = strings.Replace(name, "\n", "", -1)

	fmt.Println(name)

	// index := 0

	// ai := blackjack.NewAiPlayer(output, input, index)

	// want := "James Bond"

	// got := ai.Name

	// if want != got {
	// 	t.Fatalf("want: %q, got: %q", want, got)
	// }

}
