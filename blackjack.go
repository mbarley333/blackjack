package card

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Hand struct {
	Player string
	Cards  []Card
	Score  int
}

func (h *Hand) String() string {

	var showCards string
	for _, card := range h.Cards {
		showCards = showCards + "[" + card.String() + "]"
	}
	return h.Player + " has " + fmt.Sprint(h.Score) + ": " + showCards
}

func (h *Hand) AddCards(cards []Card) error {

	for _, card := range cards {
		h.Cards = append(h.Cards, card)
		h.ScoreCards(card)
	}

	return nil
}

func (h Hand) DealerString() string {

	return h.Player + ": [" + h.Cards[0].String() + "]" + "[???]"

}

func (h *Hand) Hit(shuffled *Deck) error {

	card, err := shuffled.Deal(1)
	if err != nil {
		return fmt.Errorf("unable to hit,%s", err)
	}

	h.AddCards(card)
	return nil

}

func (h *Hand) ScoreCards(card Card) {

	var cardScore int

	if card.Rank > Ten {
		cardScore = 10
	} else if card.Rank == Ace {

	} else {
		cardScore = int(card.Rank)
	}
	h.Score += cardScore

}

func (h *Hand) EvaluateAce() int {

	switch len(h.Cards) {
	case 1:
		return 0
	default:
		ace1 := h.Score + 1
		ace11 := h.Score + 11
		fmt.Println("Would you like %d or %d?", ace1, ace11)

	}

}

func NewBlackjackGame() error {
	deck := NewDeck(
		WithNumberOfDecks(3),
	)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffledDeck := deck.Shuffle(random)

	//opening
	player := Hand{Player: "Player1"}
	dealer := Hand{Player: "Dealer"}

	playerCards, err := shuffledDeck.Deal(2)
	if err != nil {
		return fmt.Errorf("unable to deal cards to %s, %s", player.Player, err)
	}
	dealerCards, err := shuffledDeck.Deal(2)
	if err != nil {
		return fmt.Errorf("unable to deal cards to %s, %s", player.Player, err)
	}

	player.AddCards(playerCards)
	dealer.AddCards(dealerCards)

	fmt.Println(dealer.DealerString())
	fmt.Println(player.String())

	var status string

	for strings.ToLower(status) != "s" {
		fmt.Printf("%s would you like to (H)it or (S)tand?", player.Player)
		fmt.Scanln(&status)
		if strings.ToLower(status) == "h" {
			err := player.Hit(&shuffledDeck)
			if err != nil {
				return err
			}
			fmt.Println(player.String())

		}
	}

	//dealer
	fmt.Println(dealer.String())
	fmt.Println(player.String())

	return nil
}
