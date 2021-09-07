package blackjack

import (
	"cards"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

type Action int

var ActionStringMap = map[Action]string{
	ActionHit:   "Hit",
	ActionStand: "Stand",
	ActionQuit:  "Quit",
	None:        "Invalid Action",
}

func (a Action) String() string {

	return ActionStringMap[a]
}

const (
	None Action = iota
	ActionStand
	ActionHit
	ActionQuit
)

var ActionMap = map[string]Action{
	"h": ActionHit,
	"s": ActionStand,
	"q": ActionQuit,
	"n": None,
}

type Outcome int

const (
	OutcomeNone Outcome = iota
	OutcomeBlackjack
	OutcomeWin
	OutcomeLose
	OutcomeTie
	OutcomeBust
)

var ReportMap = map[Outcome]string{
	OutcomeBlackjack: "***** Blackjack!  Player wins *****",
	OutcomeWin:       "***** Player wins! *****",
	OutcomeLose:      "***** Player loses *****",
	OutcomeTie:       "***** Player and Dealer tie *****",
	OutcomeBust:      "***** Bust!  Player loses *****",
}

type Ai int

const (
	AiNone Ai = iota
	AiStandOnly
)

type Game struct {
	Players         []Player
	Player          Player
	Dealer          Player
	Shoe            cards.Deck
	Random          rand.Rand
	GetPlayerAction func(io.Writer, io.Reader) Action
	AiHandsToPlay   int
	HandsPlayed     int
	PlayerWin       int
	PlayerLose      int
	PlayerTie       int
	Output          io.Writer
	input           io.Reader
}

type Option func(*Game) error

func WithCustomDeck(deck cards.Deck) Option {
	return func(g *Game) error {
		g.Shoe = deck
		return nil
	}
}

func WithAiType(ai Ai) Option {
	return func(g *Game) error {
		if ai == AiStandOnly {
			g.GetPlayerAction = GetAiActionStandOnly
		}
		return nil
	}
}

func WithAiHandsToPlay(number int) Option {
	return func(g *Game) error {
		g.AiHandsToPlay = number
		return nil
	}
}

func NewBlackjackGame(opts ...Option) (*Game, error) {

	deck := cards.NewDeck(
		cards.WithNumberOfDecks(3),
	)

	// func opt for input
	game := &Game{
		Shoe:            deck,
		GetPlayerAction: GetPlayerAction,
		Output:          os.Stdout,
		input:           os.Stdin,
	}

	for _, o := range opts {
		o(game)
	}

	return game, nil
}

func (g *Game) PlayerSetup(output io.Writer, input io.Reader) {

	var answer int

	for answer < 1 {
		fmt.Fprintln(output, "Please enter number of Blackjack players:")
		fmt.Fscanln(input, &answer)

		for i := 0; i <= answer; i++ {
			g.Players = append(g.Players, Player{})
		}
	}

}

func (g *Game) RunCLI() {

	for g.Player.Action != ActionQuit {
		g.Player.Hand.Cards = []cards.Card{}
		g.Dealer.Hand.Cards = []cards.Card{}
		g.Player.Action = None
		g.Player.HandOutcome = OutcomeNone
		g.Start()

	}

	fmt.Fprintln(g.Output, "")
	fmt.Fprintln(g.Output, "***** Player Report *****")
	fmt.Fprintln(g.Output, g.GetPlayerReport())

}

func (g *Game) Start() {

	fmt.Fprintln(g.Output, "")
	fmt.Fprintln(g.Output, "****** NEW GAME ******")

	// !!! was this the intent
	g.Player.Hand.Deal(&g.Shoe)
	g.Dealer.Hand.Deal(&g.Shoe)
	g.Player.Hand.Deal(&g.Shoe)
	g.Dealer.Hand.Deal(&g.Shoe)

	fmt.Fprintln(g.Output, g.Dealer.DealerString())
	fmt.Fprintln(g.Output, "Player has "+g.Player.String())

	if g.Player.Score() == 21 {
		g.Player.HandOutcome = OutcomeBlackjack
	}

	// !!! a lot of status type logic
	//
	for g.Player.HandOutcome != OutcomeBlackjack && g.Player.HandOutcome != OutcomeBust && g.Player.Action != ActionStand {

		g.SetPlayerAction()

		// !!! not sure if break is a good design
		if g.Player.Action == ActionQuit {
			break
		}

		if g.Player.Action == ActionHit {

			g.Player.Hand.Deal(&g.Shoe)
			g.Player.Action = None

			fmt.Fprintln(g.Output, "Player has "+g.Player.String())
			if g.Player.Score() > 21 {
				g.Player.HandOutcome = OutcomeBust
			}

		}
	}

	// more status logic
	if g.Player.Action != ActionQuit {
		if g.Player.HandOutcome <= OutcomeBlackjack && g.Player.HandOutcome <= OutcomeBust {
			fmt.Fprintln(g.Output, "")
			fmt.Fprintln(g.Output, "****** FINAL ROUND ******")

			for g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {
				g.Dealer.Hand.Deal(&g.Shoe)
			}
			g.Dealer.Action = ActionStand
		}
		fmt.Fprintln(g.Output, "Dealer has "+g.Dealer.String())
		fmt.Fprintln(g.Output, "Player has "+g.Player.String())
		fmt.Fprintln(g.Output, "")

		g.Outcome()
		fmt.Fprintln(g.Output, ReportMap[g.Player.HandOutcome])
		g.SetPlayerWinLoseTie(g.Player.HandOutcome)
		g.HandsPlayed += 1
		g.SetPlayerActionForAiHandsPlayed()
	}

}

func (g *Game) Outcome() {

	var outcome Outcome

	if g.Player.HandOutcome == OutcomeBlackjack || g.Player.HandOutcome == OutcomeBust {
		outcome = g.Player.HandOutcome
	} else if g.Dealer.Score() > 21 {
		outcome = OutcomeWin
	} else if g.Player.Score() > g.Dealer.Score() {
		outcome = OutcomeWin
	} else if g.Player.Score() < g.Dealer.Score() {
		outcome = OutcomeLose
	} else {
		outcome = OutcomeTie
	}

	g.Player.HandOutcome = outcome
}

func (g *Game) SetPlayerWinLoseTie(outcome Outcome) {
	if outcome == OutcomeWin || outcome == OutcomeBlackjack {
		g.PlayerWin += 1
	} else if outcome == OutcomeTie {
		g.PlayerTie += 1
	} else {
		g.PlayerLose += 1
	}
}

func (g *Game) SetPlayerAction() {

	if g.Player.Action == None {
		g.Player.Action = g.GetPlayerAction(g.Output, g.input)
	}
}

func (g *Game) SetPlayerActionForAiHandsPlayed() {
	if g.HandsPlayed == g.AiHandsToPlay {
		g.Player.Action = ActionQuit
	}
}

func (g *Game) GetPlayerReport() string {
	return "Player won: " + fmt.Sprint(g.PlayerWin) + ", lost: " + fmt.Sprint(g.PlayerLose) + " and tied: " + fmt.Sprint(g.PlayerTie)
}

func GetPlayerAction(output io.Writer, input io.Reader) Action {

	var answer string
	fmt.Fprintln(output, "Please choose (H)it, (S)tand or (Q)uit")
	fmt.Fscanln(input, &answer)

	return ActionMap[strings.ToLower(answer)]
}

func GetAiActionStandOnly(io.Writer, io.Reader) Action {
	return ActionStand
}

type Hand struct {
	Cards []cards.Card
}

// put on deal on player struct
func (h *Hand) Deal(shoe *cards.Deck) {
	var card cards.Card

	card, shoe.Cards = shoe.Cards[0], shoe.Cards[1:]
	shoe.Cards = append(shoe.Cards, card)

	h.Cards = append(h.Cards, card)

}

type Player struct {
	Hand            Hand
	Action          Action
	HandOutcome     Outcome
	GetPlayerAction func(io.Writer, io.Reader) Action
}

func (p Player) String() string {

	builder := strings.Builder{}
	for _, card := range p.Hand.Cards {
		builder.WriteString("[" + card.String() + "]")
	}
	return fmt.Sprint(p.Score()) + ": " + builder.String()
}

func (p Player) DealerString() string {

	return "Dealer has: [" + p.Hand.Cards[0].String() + "]" + "[???]"
}

func (p Player) Score() int {
	minScore := p.MinScore()

	if minScore > 11 {
		return minScore
	}
	for _, c := range p.Hand.Cards {
		if c.Rank == cards.Ace {
			return minScore + 10
		}
	}
	return minScore
}

func (p Player) MinScore() int {
	score := 0
	for _, c := range p.Hand.Cards {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// additional features
// difficult to fake concreate...take interface instead
// burn cards
// 4. basic strategy ai
// 3. reshuffle
// 2. betting - push outcome...just bet $10, handle no money player and dealer
// 1. multi players
// multi hands for split

// fold
// card counter
