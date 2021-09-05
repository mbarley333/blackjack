package aceornothing

import (
	"cards"
	"fmt"
	"strings"
)

type Hand struct {
	Cards []cards.Card
}

func (h *Hand) Deal(deck *cards.Deck) {
	var card cards.Card

	card, deck.Cards = deck.Cards[0], deck.Cards[1:]
	deck.Cards = append(deck.Cards, card)

	h.Cards = append(h.Cards, card)
}

type Game struct {
	Hand Hand
	Shoe cards.Deck
}
type Option func(*Game) error

func WithCustomDeck(deck cards.Deck) Option {
	return func(g *Game) error {
		g.Shoe = deck
		return nil
	}
}

func NewGame(opts ...Option) *Game {

	deck := cards.NewDeck(
		cards.WithNumberOfDecks(3),
	)

	game := &Game{
		Shoe: deck,
	}

	for _, o := range opts {
		o(game)
	}

	return game

}

func (g *Game) Start() {
	var response string
	for {
		fmt.Println("Would you like to play Ace Or Nothing? Please enter (Y)es or (N)o):")
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" {
			g.Hand.Deal(&g.Shoe)

			result, err := g.Evaluate()
			if err != nil {
				fmt.Printf("unable to evaluate hand,%s", err)
			}
			fmt.Println(result)
		} else if strings.ToLower(response) == "n" {
			fmt.Println("Thank you for playing!")
			break

		}
	}
}

func (g Game) Evaluate() (string, error) {

	if g.Hand.Cards[0].Rank == cards.Ace {
		result := g.Hand.Cards[0].String() + ": WIN"
		return result, nil
	}
	result := g.Hand.Cards[0].String() + ": LOSE"

	return result, nil
}
