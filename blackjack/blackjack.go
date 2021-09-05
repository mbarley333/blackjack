package blackjack

import (
	"cards"
	"fmt"
	"math/rand"
	"strings"
)

type Action int

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

type Game struct {
	Player Player
	Dealer Player
	Shoe   cards.Deck
	Random rand.Rand
}

type Option func(*Game) error

func WithCustomDeck(deck cards.Deck) Option {
	return func(g *Game) error {
		g.Shoe = deck
		return nil
	}
}

func NewBlackjackGame(opts ...Option) (*Game, error) {

	deck := cards.NewDeck(
		cards.WithNumberOfDecks(3),
	)

	game := &Game{
		Shoe: deck,
	}

	for _, o := range opts {
		o(game)
	}

	return game, nil
}

func (g *Game) RunCLI() {

	for g.Player.Action != ActionQuit {
		g.Player.Hand.Cards = []cards.Card{}
		g.Dealer.Hand.Cards = []cards.Card{}
		g.Player.Action = None
		g.Player.HandOutcome = OutcomeNone
		g.Start()
		outcome := g.Outcome()
		fmt.Println(ReportMap[outcome])

	}

}

func (g *Game) Start() {

	fmt.Println("")
	fmt.Println("****** NEW GAME ******")

	g.Player.Hand.Deal(&g.Shoe)
	g.Dealer.Hand.Deal(&g.Shoe)
	g.Player.Hand.Deal(&g.Shoe)
	g.Dealer.Hand.Deal(&g.Shoe)

	fmt.Println(g.Dealer.DealerString())
	fmt.Println("Player has " + g.Player.String())

	if g.Player.Score() == 21 {
		g.Player.HandOutcome = OutcomeBlackjack
	}

	for g.Player.HandOutcome != OutcomeBlackjack && g.Player.HandOutcome != OutcomeBust && g.Player.Action != ActionStand {

		g.SetPlayerAction()

		if g.Player.Action == ActionHit {

			g.Player.Hand.Deal(&g.Shoe)
			g.Player.Action = None

			fmt.Println("Player has " + g.Player.String())
			if g.Player.Score() > 21 {
				g.Player.HandOutcome = OutcomeBust
			}

		}
	}

	if g.Player.HandOutcome <= OutcomeBlackjack && g.Player.HandOutcome <= OutcomeBust {
		fmt.Println("")
		fmt.Println("****** FINAL ROUND ******")

		for g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {

			g.Dealer.Hand.Deal(&g.Shoe)

			fmt.Println("Dealer has " + g.Dealer.String())
			fmt.Println("Player has " + g.Player.String())
			fmt.Println("")

		}
		g.Dealer.Action = ActionStand
	}

}

func (g *Game) SetPlayerAction() {

	if g.Player.Action == None {
		answer := GetPlayerAction()
		g.Player.Action = ActionMap[strings.ToLower(answer)]
	}
}

func GetPlayerAction() string {

	var answer string
	fmt.Println("Please choose (H)it, (S)tand or (Q)uit")
	fmt.Scanf("%s\n", &answer)
	return answer
}

func (g *Game) Outcome() Outcome {

	var outcome Outcome

	if g.Dealer.Score() > 21 {
		outcome = OutcomeWin
	} else if g.Player.Score() > g.Dealer.Score() {
		outcome = OutcomeWin
	} else if g.Player.Score() < g.Dealer.Score() {
		outcome = OutcomeLose
	} else {
		outcome = OutcomeTie
	}

	return outcome
}

type Hand struct {
	Cards []cards.Card
}

func (h *Hand) Deal(shoe *cards.Deck) {
	var card cards.Card

	card, shoe.Cards = shoe.Cards[0], shoe.Cards[1:]
	shoe.Cards = append(shoe.Cards, card)

	h.Cards = append(h.Cards, card)

}

type Player struct {
	Hand        Hand
	Action      Action
	HandOutcome Outcome
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
