package blackjack

import (
	"cards"
	"fmt"
	"io"
)

func AiActionBasic(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card, index int) Action {

	var action Action
	handValue := player.Hands[index].Score()
	dealerCardValue := ScoreDealerHoleCard(dealerCard)

	var isSoft bool

	for _, card := range player.Hands[index].Cards {

		length := len(player.Hands[index].Cards)
		if card.Rank == cards.Ace && length == 2 {
			isSoft = true
		}
	}

	if (handValue == 10 && dealerCardValue < handValue || handValue == 11 && dealerCardValue < handValue) && player.Cash > player.Hands[index].Bet {
		action = ActionDoubleDown
	} else if handValue == 9 && dealerCardValue >= 3 && dealerCardValue <= 6 && player.Cash > player.Hands[index].Bet {
		action = ActionDoubleDown
	} else if handValue <= 11 {
		action = ActionHit
	} else if handValue <= 15 && isSoft {
		action = ActionHit
	} else if handValue >= 19 && isSoft {
		action = ActionStand
	} else if handValue >= 16 && handValue <= 18 && isSoft && dealerCardValue >= 7 {
		action = ActionHit
	} else if handValue >= 16 && handValue <= 18 && isSoft && dealerCardValue <= 6 && player.Cash > player.Hands[index].Bet {
		action = ActionDoubleDown
	} else if handValue >= 16 && handValue <= 18 && isSoft && dealerCardValue <= 6 && player.Cash < player.Hands[index].Bet {
		action = ActionHit
	} else if handValue >= 17 && handValue <= 21 {
		action = ActionStand
	} else if handValue == 12 && dealerCardValue <= 3 {
		action = ActionHit
	} else if handValue >= 12 && handValue <= 16 && dealerCardValue <= 6 {
		action = ActionStand
	} else if handValue >= 12 && handValue <= 16 && dealerCardValue >= 7 {
		action = ActionHit
	}
	return action
}

func AiActionStandOnly(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card, index int) Action {

	return ActionStand
}

func AiBet(output io.Writer, input io.Reader, player *Player, index int) error {

	if player.Record.HandsPlayed == player.AiHandsToPlay {
		player.Action = ActionQuit
	} else {
		fmt.Println(player)
		bet := 1
		player.Cash -= bet
		player.Hands[index].Bet += bet

	}
	return nil
}
