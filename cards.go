package cards

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit int

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

type Rank int

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

var ranks = []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

type Card struct {
	Rank
	Suit
}

func (c Card) String() string {
	suitMap := make(map[int]string)
	suitMap[int(Spade)] = "Spade"
	suitMap[int(Diamond)] = "Diamond"
	suitMap[int(Club)] = "Club"
	suitMap[int(Heart)] = "Heart"

	rankMap := make(map[int]string)
	rankMap[int(Ace)] = "Ace"
	rankMap[int(Two)] = "Two"
	rankMap[int(Three)] = "Three"
	rankMap[int(Four)] = "Four"
	rankMap[int(Five)] = "Five"
	rankMap[int(Six)] = "Six"
	rankMap[int(Seven)] = "Seven"
	rankMap[int(Eight)] = "Eight"
	rankMap[int(Nine)] = "Nine"
	rankMap[int(Ten)] = "Ten"
	rankMap[int(Jack)] = "Jack"
	rankMap[int(Queen)] = "Queen"
	rankMap[int(King)] = "King"

	cardName := fmt.Sprintf("%s of %ss", rankMap[int(c.Rank)], suitMap[int(c.Suit)])
	return cardName
}

type Deck struct {
	Cards  []Card
	count  int
	random *rand.Rand
}

type Option func(*Deck) error

func WithNumberOfDecks(number int) Option {
	return func(d *Deck) error {
		d.count = number
		return nil
	}
}

func WithRandom(random *rand.Rand) Option {
	return func(d *Deck) error {
		d.random = random
		return nil
	}
}

func NewDeck(opts ...Option) Deck {

	deck := &Deck{
		count:  1,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, o := range opts {
		o(deck)
	}

	for i := 0; i < deck.count; i++ {
		for _, suit := range suits {
			for _, rank := range ranks {
				deck.Cards = append(deck.Cards, Card{Suit: suit, Rank: rank})
			}

		}
	}
	shuffled_cards := make([]Card, len(deck.Cards))
	perm := deck.random.Perm(len(deck.Cards))

	for i, j := range perm {
		shuffled_cards[i] = deck.Cards[j]
	}
	deck.Cards = shuffled_cards
	return Deck{
		Cards: shuffled_cards,
	}

}
