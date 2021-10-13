package blackjack

import (
	"github.com/mbarley333/cards"
)

func CountHiLo(card cards.Card, count int, numberCardsDealt int, deckCount int) (int, float64) {

	if card.Rank >= 2 && card.Rank <= 6 {
		count += 1
	} else if card.Rank >= 10 {
		count -= 1
	}

	fNumberCardsDealt, fDeckCount, fCount := float64(numberCardsDealt), float64(deckCount), float64(count)
	trueCount := fCount / (fDeckCount - (fNumberCardsDealt / 52.0))

	return count, trueCount
}
