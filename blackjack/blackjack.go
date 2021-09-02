package blackjack

import (
	"cards"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Action int

const (
	None Action = iota
	Stand
	Hit
)

type Outcome int

const (
	PlayerBlackjack Outcome = iota
	PlayerWin
	PlayerLose
	PlayerTie
	PlayerBust
)

type Game struct {
	Player Player
	Dealer Player
	Shoe   cards.Deck
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
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	game := &Game{
		Shoe: deck.Shuffle(random),
	}

	for _, o := range opts {
		o(game)
	}

	return game, nil
}

func (g *Game) Deal() cards.Card {
	var card cards.Card

	card, g.Shoe.Cards = g.Shoe.Cards[0], g.Shoe.Cards[1:]
	g.Shoe.Cards = append(g.Shoe.Cards, card)

	return card

}

func (g *Game) ResetBlackjack() {
	g.Player.Hand = []cards.Card{}
	g.Dealer.Hand = []cards.Card{}
	g.Player.Action = None
	g.Start()
	g.DealerStart()
	outcome := g.Outcome()
	fmt.Println(ReportWinner(outcome))
}

func (g *Game) Start() {

	g.Player.GetCard(g.Deal())
	g.Dealer.GetCard(g.Deal())
	g.Player.GetCard(g.Deal())
	g.Dealer.GetCard(g.Deal())

	fmt.Println(g.Dealer.DealerString())
	fmt.Println("Player has " + g.Player.String())

	if g.Player.Score() == 21 {
		g.Player.Action = Stand
		g.Outcome()

	}

	for g.Player.Action != Stand {

		g.SetPlayerAction()

		if g.Player.Action == Hit {
			g.Player.GetCard(g.Deal())
			g.Player.Action = None

			fmt.Println(g.Player.String())
			if g.Player.Score() > 21 {
				g.Player.Action = Stand
				g.Outcome()
			}

		}
	}

}

func (g *Game) DealerStart() {

	fmt.Println("")
	fmt.Println("******FINAL ROUND******")

	for g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {

		g.Dealer.GetCard(g.Deal())

		fmt.Println("Dealer has " + g.Dealer.String())
		fmt.Println("Player has " + g.Player.String())
		fmt.Println("")

	}
	g.Dealer.Action = Stand
}

func (g *Game) SetPlayerAction() {

	var answer string

	if g.Player.Action == None {
		fmt.Println("Please choose (H)it or (S)tand")
		fmt.Scanf("%s\n", &answer)

		if strings.ToLower(answer) == "h" {
			g.Player.Action = Hit
		} else {
			g.Player.Action = Stand
		}

	}

}

func (g *Game) Outcome() Outcome {

	if g.Player.Score() == 21 && len(g.Player.Hand) == 2 {
		return PlayerBlackjack
	} else if g.Dealer.Score() > 21 {
		return PlayerWin
	} else if g.Player.Score() > 21 {
		return PlayerBust
	} else if g.Player.Score() > g.Dealer.Score() {
		return PlayerWin
	}

	return PlayerLose

}

func ReportWinner(outcome Outcome) string {
	reportMap := make(map[Outcome]string)
	reportMap[PlayerBlackjack] = "Blackjack!  Player wins"
	reportMap[PlayerWin] = "Player wins!"
	reportMap[PlayerLose] = "Player loses"
	reportMap[PlayerTie] = "Player and Dealer tie"
	reportMap[PlayerBust] = "Bust!  Player loses"

	return "***** " + reportMap[outcome] + " *****"
}

type Player struct {
	Hand   []cards.Card
	Action Action
}

func (p Player) String() string {

	var showCards string
	for _, card := range p.Hand {
		showCards = showCards + "[" + card.String() + "]"
	}
	return fmt.Sprint(p.Score()) + ": " + showCards
}

func (p Player) DealerString() string {

	return "Dealer has: [" + p.Hand[0].String() + "]" + "[???]"

}

func (p *Player) GetCard(card cards.Card) {
	p.Hand = append(p.Hand, card)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
