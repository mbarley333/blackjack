package blackjack

import (
	"cards"
	"io"
)

func AiActionBasic(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card) Action {

	var action Action
	handValue := player.Score()
	dealerCardValue := ScoreDealerHoleCard(dealerCard)

	var isSoft bool

	for _, card := range player.Hand {
		if card.Rank == cards.Ace {
			isSoft = true
		}
	}

	if (handValue == 10 && dealerCardValue < handValue || handValue == 11 && dealerCardValue < handValue) && player.Cash > player.HandBet {
		action = ActionDoubleDown
	} else if handValue == 9 && dealerCardValue >= 3 && dealerCardValue <= 6 && player.Cash > player.HandBet {
		action = ActionDoubleDown
	} else if handValue <= 11 {
		action = ActionHit
	} else if handValue <= 15 && isSoft {
		action = ActionHit
	} else if handValue >= 19 && isSoft {
		action = ActionStand
	} else if handValue >= 16 && handValue <= 18 && isSoft && dealerCardValue >= 7 {
		action = ActionHit
	} else if handValue >= 16 && handValue <= 18 && isSoft && dealerCardValue <= 6 && player.Cash > player.HandBet {
		action = ActionDoubleDown
	} else if handValue >= 16 && handValue <= 18 && isSoft && dealerCardValue <= 6 && player.Cash < player.HandBet {
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

func AiActionStandOnly(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card) Action {

	return ActionStand
}

func AiBet(output io.Writer, input io.Reader, player *Player) error {

	if player.Record.HandsPlayed == player.AiHandsToPlay {
		player.Action = ActionQuit
	} else {
		bet := 1
		player.Cash -= bet
		player.HandBet += bet

	}
	return nil
}
