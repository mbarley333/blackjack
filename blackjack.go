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
}

func (h *Hand) String() string {

	var showCards string
	for _, card := range h.Cards {
		showCards = showCards + "[" + card.String() + "]"
	}
	return h.Player + " has " + fmt.Sprint(h.Score()) + ": " + showCards
}

func (h Hand) DealerString() string {

	return h.Player + ": [" + h.Cards[0].String() + "]" + "[???]"

}

func (h *Hand) Hit(shuffled *Deck) error {

	card, err := shuffled.Deal(1)
	h.Cards = append(h.Cards, card...)
	if err != nil {
		return fmt.Errorf("unable to hit,%s", err)
	}

	return nil

}

func (h Hand) Score() int {
	minScore := h.MinScore()

	if minScore > 11 {
		return minScore
	}
	for _, c := range h.Cards {
		if c.Rank == Ace {
			return minScore + 10
		}
	}
	return minScore
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h.Cards {
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

func NewBlackjackGame() error {
	deck := NewDeck(
		WithNumberOfDecks(3),
	)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffledDeck := deck.Shuffle(random)

	player := Hand{Player: "Player1"}
	dealer := Hand{Player: "Dealer"}

	var err error

	player.Cards, err = shuffledDeck.Deal(2)
	if err != nil {
		return fmt.Errorf("unable to deal cards to %s, %s", player.Player, err)
	}
	dealer.Cards, err = shuffledDeck.Deal(2)
	if err != nil {
		return fmt.Errorf("unable to deal cards to %s, %s", dealer.Player, err)
	}

	fmt.Println(dealer.DealerString())
	fmt.Println(player.String())

	var status string

	for strings.ToLower(status) != "s" && strings.ToLower(status) != "b" {
		fmt.Printf("%s would you like to (H)it or (S)tand? ", player.Player)
		fmt.Scanln(&status)
		if strings.ToLower(status) == "h" {
			err := player.Hit(&shuffledDeck)
			if err != nil {
				return err
			}
			fmt.Println(player.String())
			if player.Score() > 21 {
				status = "b"
			}

		}
	}

	if status != "b" {
		fmt.Printf("\n****************FINAL ROUND***************\n")

		for dealer.Score() <= 16 || dealer.MinScore() < 17 {
			err := dealer.Hit(&shuffledDeck)
			if err != nil {
				return err
			}
		}
		fmt.Println(dealer.String())
		fmt.Println(player.String())

		playerWin := IsPlayerWinner(player.Score(), dealer.Score())
		if playerWin {
			fmt.Println(player.Player + " WINS!!!")
		} else {
			fmt.Println(player.Player + " LOSES")
		}
	}

	return nil
}

func IsPlayerWinner(player, dealer int) bool {
	if dealer > 21 {
		return true
	}
	return player > dealer
}
