package card

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
	Suit
	Rank
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

type Deck []Card

func NewDeck() Deck {

	deck := Deck{}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}

	}
	return deck

}

func (d *Deck) Shuffle(random *rand.Rand) Deck {

	shuffled_deck := make([]Card, len(*d))
	perm := random.Perm(len(*d))

	for i, j := range perm {
		shuffled_deck[i] = (*d)[j]
	}
	return shuffled_deck

}

func (d *Deck) Deal(numberCards int) ([]Card, error) {
	var card Card
	var cards []Card

	for i := 0; i < numberCards; i++ {
		card, *d = (*d)[0], (*d)[1:]
		cards = append(cards, card)

	}
	return cards, nil

}

func EvaluateAceOrNothing(hand []Card) (string, error) {

	if hand[0].Rank == Ace {
		return "Ace: WIN", nil
	}
	result := hand[0].String() + ": LOSE"

	return result, nil
}

func NewAceOrNothing() string {
	deck := NewDeck()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffledDeck := deck.Shuffle(random)
	hand, err := shuffledDeck.Deal(1)
	if err != nil {
		fmt.Printf("unable to deal card, %s", err)
	}
	result, err := EvaluateAceOrNothing(hand)
	if err != nil {
		fmt.Printf("unable to evaluate hand,%s", err)
	}
	return result
}
