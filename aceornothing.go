package card

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func EvaluateAceOrNothing(hand []Card) (string, error) {

	if hand[0].Rank == Ace {
		result := hand[0].String() + ": WIN"
		return result, nil
	}
	result := hand[0].String() + ": LOSE"

	return result, nil
}

func NewAceOrNothing() {
	deck := NewDeck(
		WithNumberOfDecks(3),
	)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffledDeck := deck.Shuffle(random)

	var response string
	for {
		fmt.Println("Would you like to play Ace Or Nothing? Please enter (Y)es or (N)o):")
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" {
			hand, err := shuffledDeck.Deal(1)
			if err != nil {
				fmt.Printf("unable to deal card, %s", err)
			}
			result, err := EvaluateAceOrNothing(hand)
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
