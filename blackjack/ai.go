package blackjack

import (
	"cards"
	"fmt"
	"io"
)

func AiActionBasic(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card, index int, c CardCounter) Action {

	var action Action
	handValue := player.Hands[index].Score()
	dealerCardValue := ScoreDealerHoleCard(dealerCard)

	isSoft := false
	isSplitable := false

	if player.Hands[index].Cards[0].Rank == player.Hands[index].Cards[1].Rank {
		isSplitable = true
	}

	for _, card := range player.Hands[index].Cards {
		length := len(player.Hands[index].Cards)
		if card.Rank == cards.Ace && length == 2 {
			isSoft = true
		}
	}

	// split aces and eights
	if isSplitable && player.Hands[index].Cards[0].Rank == cards.Ace || player.Hands[index].Cards[0].Rank == cards.Eight {
		action = ActionSplit
		// split all pairs when dealer showing 6 or less AND pair != 4,5,10
	} else if isSplitable && player.Hands[index].Cards[0].Rank != cards.Five && player.Hands[index].Cards[0].Rank != cards.Four && player.Hands[index].Cards[0].Rank <= 9 && dealerCardValue <= 6 {
		action = ActionSplit
	} else if (handValue == 10 && dealerCardValue < handValue || handValue == 11 && dealerCardValue < handValue) && player.Cash > player.Hands[index].Bet {
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

func AiActionStandOnly(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card, index int, c CardCounter) Action {

	return ActionStand
}

func AiBet(output io.Writer, input io.Reader, player *Player, index int, c CardCounter) error {

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
