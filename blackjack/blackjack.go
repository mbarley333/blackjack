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

var OutcomeStringMap = map[Outcome]string{
	OutcomeNone:      "Invalid Outcome",
	OutcomeBlackjack: "Blackjack",
	OutcomeWin:       "Win",
	OutcomeLose:      "Lose",
	OutcomeTie:       "Tie",
	OutcomeBust:      "Bust",
}

func (o Outcome) String() string {

	return OutcomeStringMap[o]
}

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
	Players []Player
	Dealer  Player
	Shoe    cards.Deck
	Random  rand.Rand
	output  io.Writer
	input   io.Reader
}

type Option func(*Game) error

func WithCustomDeck(deck cards.Deck) Option {
	return func(g *Game) error {
		g.Shoe = deck
		return nil
	}
}

// func WithAiType(ai Ai) Option {
// 	return func(g *Game) error {
// 		if ai == AiStandOnly {
// 			g.GetPlayerAction = GetAiActionStandOnly
// 		}
// 		return nil
// 	}
// }

// func WithAiHandsToPlay(number int) Option {
// 	return func(g *Game) error {
// 		g.AiHandsToPlay = number
// 		return nil
// 	}
// }

func WithOutput(output io.Writer) Option {
	return func(g *Game) error {
		g.output = output
		return nil
	}
}

func NewBlackjackGame(opts ...Option) (*Game, error) {

	deck := cards.NewDeck(
		cards.WithNumberOfDecks(3),
	)

	// func opt for input
	game := &Game{
		Shoe:   deck,
		output: os.Stdout,
		input:  os.Stdin,
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

		for i := 1; i <= answer; i++ {
			g.Players = append(g.Players, Player{})
		}
	}

}

func (g *Game) RunCLI() {

	// for g.Player.Action != ActionQuit {
	// 	g.ResetPlayers()
	// 	g.Start()
	// }
	// fmt.Fprintln(g.output, "")
	// fmt.Fprintln(g.output, "***** Player Report *****")
	// fmt.Fprintln(g.output, g.GetPlayerReport())

	for len(g.Players) > 0 {
		g.ResetPlayers()
		g.Start()
	}

}

func (g *Game) Start() {

	fmt.Fprintln(g.output, "")
	fmt.Fprintln(g.output, "****** NEW GAME ******")

	g.OpeningDeal()

	fmt.Fprintln(g.output, g.Dealer.DealerString())
	g.ShowPlayerCards(g.output)

	for index := range g.Players {

		if g.Players[index].Score() == 21 {
			g.Players[index].HandOutcome = OutcomeBlackjack
			//ok = false
		}

		for g.Players[index].GetStatus() {
			g.Players[index].SetPlayerAction(g.output, g.input)

			if g.Players[index].Action == ActionHit {

				card := g.Deal()
				g.Players[index].Hand = append(g.Players[index].Hand, card)
				g.Players[index].Action = None

				fmt.Fprintln(g.output, g.Players[index].Name+" has "+g.Players[index].String())
				if g.Players[index].Score() > 21 {
					g.Players[index].HandOutcome = OutcomeBust
					//ok = false
				}

			}
		}
	}

	fmt.Fprintln(g.output, "")
	g.DealerPlay()

	fmt.Fprintln(g.output, "Dealer has "+g.Dealer.String())
	g.ShowPlayerCards(g.output)
	fmt.Fprintln(g.output, "")

	g.Outcome(g.output)

	// // more status logic
	// if g.Player.Action != ActionQuit {
	// 	if g.Player.HandOutcome <= OutcomeBlackjack && g.Player.HandOutcome <= OutcomeBust {
	// 		fmt.Fprintln(g.output, "")
	// 		fmt.Fprintln(g.output, "****** FINAL ROUND ******")

	// 		for g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {
	// 			card := g.Deal()
	// 			g.Dealer.Hand = append(g.Dealer.Hand, card)
	// 		}
	// 		g.Dealer.Action = ActionStand
	// 	}
	// 	fmt.Fprintln(g.output, "Dealer has "+g.Dealer.String())
	// 	fmt.Fprintln(g.output, "Player has "+g.Player.String())
	// 	fmt.Fprintln(g.output, "")

	// 	g.Outcome()
	// 	fmt.Fprintln(g.output, ReportMap[g.Player.HandOutcome])
	// 	g.SetPlayerWinLoseTie(g.Player.HandOutcome)
	// 	g.HandsPlayed += 1
	// 	g.SetPlayerActionForAiHandsPlayed()
	// }

}

func (g *Game) DealerPlay() {
	fmt.Fprintln(g.output, "****** DEALER'S TURN ******")

	for g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {
		card := g.Deal()
		g.Dealer.Hand = append(g.Dealer.Hand, card)
	}
	g.Dealer.Action = ActionStand
}

func (g Game) ShowPlayerCards(output io.Writer) {
	for _, player := range g.Players {
		fmt.Fprintln(g.output, player.Name+" has "+player.String())
	}
}

func (g *Game) Deal() cards.Card {

	var card cards.Card

	card, g.Shoe.Cards = g.Shoe.Cards[0], g.Shoe.Cards[1:]

	return card
}

func (g *Game) OpeningDeal() {
	for i := 0; i < 2; i++ {
		for index := range g.Players {
			card := g.Deal()
			g.Players[index].Hand = append(g.Players[index].Hand, card)
		}
		card := g.Deal()
		g.Dealer.Hand = append(g.Dealer.Hand, card)
	}
}

func (g *Game) Outcome(output io.Writer) {

	var outcome Outcome
	for index := range g.Players {
		if g.Players[index].HandOutcome == OutcomeBlackjack || g.Players[index].HandOutcome == OutcomeBust {
			outcome = g.Players[index].HandOutcome
		} else if g.Dealer.Score() > 21 {
			outcome = OutcomeWin
		} else if g.Players[index].Score() > g.Dealer.Score() {
			outcome = OutcomeWin
		} else if g.Players[index].Score() < g.Dealer.Score() {
			outcome = OutcomeLose
		} else {
			outcome = OutcomeTie
		}

		g.Players[index].HandOutcome = outcome
		fmt.Fprintln(output, g.Players[index].Name+": "+ReportMap[g.Players[index].HandOutcome])
	}
}

func (g *Game) ResetPlayers() {

	for _, player := range g.Players {
		if player.Action == ActionQuit {
			g.Players = g.Players[1:]
		}
	}

	for _, player := range g.Players {
		player.Hand = []cards.Card{}
		player.Action = None
		player.HandOutcome = OutcomeNone
	}
	g.Dealer.Hand = []cards.Card{}

}

type Player struct {
	Name            string
	Hand            []cards.Card
	Action          Action
	HandOutcome     Outcome
	GetPlayerAction func(io.Writer, io.Reader) Action
	AiHandsToPlay   int
	HandsPlayed     int
	PlayerWin       int
	PlayerLose      int
	PlayerTie       int
}

func (p Player) GetStatus() bool {
	return p.Action != ActionQuit && p.Action != ActionStand && p.HandOutcome != OutcomeBlackjack && p.HandOutcome != OutcomeBust
}

func (p *Player) GetPlayerReport() string {
	return "Player won: " + fmt.Sprint(p.PlayerWin) + ", lost: " + fmt.Sprint(p.PlayerLose) + " and tied: " + fmt.Sprint(p.PlayerTie)
}

func (p *Player) SetPlayerAction(output io.Writer, input io.Reader) {

	if p.Action == None {
		p.Action = p.GetPlayerAction(output, input)
	}
}

func (p *Player) SetPlayerActionForAiHandsPlayed() {
	if p.HandsPlayed == p.AiHandsToPlay {
		p.Action = ActionQuit
	}
}

func (p *Player) SetPlayerWinLoseTie(outcome Outcome) {
	if outcome == OutcomeWin || outcome == OutcomeBlackjack {
		p.PlayerWin += 1
	} else if outcome == OutcomeTie {
		p.PlayerTie += 1
	} else {
		p.PlayerLose += 1
	}
}

func (p Player) String() string {

	builder := strings.Builder{}
	for _, card := range p.Hand {
		builder.WriteString("[" + card.String() + "]")
	}
	return fmt.Sprint(p.Score()) + ": " + builder.String()
}

func (p Player) DealerString() string {

	return "Dealer has: [" + p.Hand[0].String() + "]" + "[???]"

}

func (p Player) Score() int {
	minScore := p.MinScore()

	if minScore > 11 {
		return minScore
	}
	for _, c := range p.Hand {
		if c.Rank == cards.Ace {
			return minScore + 10
		}
	}
	return minScore
}

func (p Player) MinScore() int {
	score := 0
	for _, c := range p.Hand {
		score += min(int(c.Rank), 10)
	}
	return score
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
