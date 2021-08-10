package card

import "math/rand"

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Suit
	Rank
}

type Deck []Card

func NewDeck(random *rand.Rand) Deck {
	// for loop to generate same cards
	// gen all cards/suits
	
	for i:=0

	return nil
}

func (d *Deck) Deal() (Card, error) {
	return Card{}, nil
}
