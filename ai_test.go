package blackjack_test

// func TestAiBasicAction(t *testing.T) {

// 	g, err := blackjack.NewBlackjackGame()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	type testCase struct {
// 		playerCards []cards.Card
// 		dealerCard  []cards.Card
// 		action      blackjack.Action
// 		bet         int
// 		cash        int
// 		description string
// 	}
// 	tcs := []testCase{
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Seven, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionDoubleDown,
// 			bet:         1,
// 			cash:        10,
// 			description: "10 or 11, dealer < 10 or 11, enough to double",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Seven, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionHit,
// 			bet:         10,
// 			cash:        1,
// 			description: "10 or 11, dealer < 10 or 11, not enough to double",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Six, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionDoubleDown,
// 			bet:         1,
// 			cash:        10,
// 			description: "9, dealer between 3 and 6, enough to double",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Six, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionHit,
// 			bet:         10,
// 			cash:        1,
// 			description: "9, dealer between 3 and 6, not enough to double",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ace, Suit: cards.Club}, {Rank: cards.Jack, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionStand,
// 			description: "Blackjack",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Five, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionHit,
// 			description: "Eleven or less",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ace, Suit: cards.Club}, {Rank: cards.Four, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionHit,
// 			description: "Soft 15 or less",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ace, Suit: cards.Club}, {Rank: cards.Eight, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionStand,
// 			description: "Soft 19 or higher",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ace, Suit: cards.Club}, {Rank: cards.Six, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Seven, Suit: cards.Club}},
// 			action:      blackjack.ActionHit,
// 			description: "Soft 16 to 18, dealer >= 7",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ace, Suit: cards.Club}, {Rank: cards.Six, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Six, Suit: cards.Club}},
// 			action:      blackjack.ActionDoubleDown,
// 			description: "Soft 16 to 18, dealer <= 6, enough to double",
// 			bet:         1,
// 			cash:        10,
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ace, Suit: cards.Club}, {Rank: cards.Six, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Six, Suit: cards.Club}},
// 			action:      blackjack.ActionHit,
// 			description: "Soft 16 to 18, dealer <= 6, not enough to double",
// 			bet:         10,
// 			cash:        1,
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ten, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}, {Rank: cards.Five, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionStand,
// 			description: "Hard 17 to 21",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ten, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}, {Rank: cards.Two, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Four, Suit: cards.Club}},
// 			action:      blackjack.ActionStand,
// 			description: "Hard 12 to 16 w/ Dealer <= 6",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Six, Suit: cards.Club}, {Rank: cards.Two, Suit: cards.Club}, {Rank: cards.Four, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Three, Suit: cards.Club}},
// 			action:      blackjack.ActionHit,
// 			description: "Hard 12 w/ Dealer <= 3",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ten, Suit: cards.Club}, {Rank: cards.Three, Suit: cards.Club}, {Rank: cards.Two, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Seven, Suit: cards.Club}},
// 			action:      blackjack.ActionHit,
// 			description: "Hard 12 to 16 w/ Dealer >= 7",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ace, Suit: cards.Club}, {Rank: cards.Ace, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Seven, Suit: cards.Club}},
// 			action:      blackjack.ActionSplit,
// 			description: "Split Aces",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Eight, Suit: cards.Club}, {Rank: cards.Eight, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Seven, Suit: cards.Club}},
// 			action:      blackjack.ActionSplit,
// 			description: "Split Eights",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Ten, Suit: cards.Club}, {Rank: cards.Ten, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Seven, Suit: cards.Club}},
// 			action:      blackjack.ActionStand,
// 			description: "No Split Tens",
// 		},
// 		{
// 			playerCards: []cards.Card{{Rank: cards.Seven, Suit: cards.Club}, {Rank: cards.Seven, Suit: cards.Club}},
// 			dealerCard:  []cards.Card{{Rank: cards.Six, Suit: cards.Club}},
// 			action:      blackjack.ActionSplit,
// 			description: "Split Seven with Dealer showing Six",
// 		},
// 	}

// 	output := &bytes.Buffer{}
// 	input := strings.NewReader("")

// 	for _, tc := range tcs {
// 		p := &blackjack.Player{
// 			Hands: []*blackjack.Hand{
// 				{
// 					Cards: tc.playerCards,
// 					Bet:   tc.bet,
// 				},
// 			},
// 			Decide: blackjack.AiActionBasic,
// 			Cash:   tc.cash,
// 		}
// 		g.AddPlayer(p)
// 		g.Dealer = &blackjack.Player{
// 			Hands: []*blackjack.Hand{
// 				{
// 					Cards: tc.dealerCard,
// 					Bet:   tc.bet,
// 				},
// 			},
// 			Decide: blackjack.AiActionBasic,
// 			Cash:   tc.cash,
// 		}

// 		want := tc.action

// 		index := 0
// 		got := g.Players[0].Decide(output, input, g.Players[0], g.Dealer.Hands[0].Cards[0], index, g.CardCounter)

// 		if want != got {
// 			t.Fatalf("%q: wanted: %q, got: %q", tc.description, want.String(), got.String())
// 		}
// 		g.Players = []*blackjack.Player{}
// 		g.Dealer.Hands = nil

// 	}

// }
