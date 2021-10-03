package blackjack

import (
	"cards"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Action int

var ActionStringMap = map[Action]string{
	None:             "Invalid Action",
	ActionHit:        "Hit",
	ActionStand:      "Stand",
	ActionQuit:       "Quit",
	ActionDoubleDown: "Double Down",
	ActionSplit:      "Split",
}

func (a Action) String() string {
	return ActionStringMap[a]
}

const (
	None Action = iota
	ActionStand
	ActionHit
	ActionQuit
	ActionDoubleDown
	ActionSplit
)

var ActionMap = map[string]Action{
	"h": ActionHit,
	"s": ActionStand,
	"q": ActionQuit,
	"n": None,
	"d": ActionDoubleDown,
	"p": ActionSplit,
}

type Outcome int

var OutcomeStringMap = map[Outcome]string{
	OutcomeNone:      "Invalid Outcome",
	OutcomeBlackjack: "Blackjack",
	OutcomeWin:       "Win",
	OutcomeLose:      "Lose",
	OutcomeTie:       "Tie",
	OutcomeBust:      "Bust",
}

func (o Outcome) String() string {
	return OutcomeStringMap[o]
}

const (
	OutcomeNone Outcome = iota
	OutcomeBlackjack
	OutcomeWin
	OutcomeLose
	OutcomeTie
	OutcomeBust
)

var ReportMap = map[Outcome]string{
	OutcomeBlackjack: "***** Blackjack!  Player wins *****",
	OutcomeWin:       "***** Player wins! *****",
	OutcomeLose:      "***** Player loses *****",
	OutcomeTie:       "***** Player and Dealer tie *****",
	OutcomeBust:      "***** Bust!  Player loses *****",
}

var BalanceReportMap = map[Outcome]string{
	OutcomeBlackjack: " won $",
	OutcomeWin:       " won $",
	OutcomeLose:      " lost $",
	OutcomeTie:       " push",
	OutcomeBust:      " lost $",
}

type PlayerType int

const (
	PlayerTypeHuman PlayerType = iota
	PlayerTypeAiStandOnly
	PlayerTypeAiBasic
)

var PlayerTypeMap = map[PlayerType]func(io.Writer, io.Reader, *Player, cards.Card, int) Action{
	PlayerTypeHuman:       HumanAction,
	PlayerTypeAiStandOnly: AiActionStandOnly,
	PlayerTypeAiBasic:     AiActionBasic,
}

var PlayerTypeInputMap = map[string]PlayerType{
	"h": PlayerTypeHuman,
	"b": PlayerTypeAiBasic,
	"s": PlayerTypeAiStandOnly,
}

var PlayerTypeBetMap = map[PlayerType]func(io.Writer, io.Reader, *Player, int) error{
	PlayerTypeHuman:       HumanBet,
	PlayerTypeAiStandOnly: AiBet,
	PlayerTypeAiBasic:     AiBet,
}

type Game struct {
	Players              []*Player
	Dealer               *Player
	Shoe                 cards.Deck
	Random               rand.Rand
	output               io.Writer
	input                io.Reader
	IsIncomingDeck       bool
	CardsDealt           int
	IncomingDeckPosition int
	DeckCount            int
	random               *rand.Rand
}

type Option func(*Game) error

func WithCustomDeck(deck cards.Deck) Option {
	return func(g *Game) error {
		g.Shoe = deck
		g.IsIncomingDeck = false
		return nil
	}
}

func WithOutput(output io.Writer) Option {
	return func(g *Game) error {
		g.output = output
		return nil
	}
}

func WithInput(input io.Reader) Option {
	return func(g *Game) error {
		g.input = input
		return nil
	}
}

func WithIncomingDeck(r bool) Option {
	return func(g *Game) error {
		g.IsIncomingDeck = r
		return nil
	}
}

func WithDeckCount(d int) Option {
	return func(g *Game) error {
		g.DeckCount = d
		return nil
	}
}

func WithRandom(random *rand.Rand) Option {
	return func(g *Game) error {
		g.random = random
		return nil
	}
}

func NewBlackjackGame(opts ...Option) (*Game, error) {

	game := &Game{
		output:         os.Stdout,
		input:          os.Stdin,
		IsIncomingDeck: true,
		DeckCount:      6,
		random:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	deck := cards.NewDeck(
		cards.WithNumberOfDecks(game.DeckCount),
	)
	game.Shoe = deck

	for _, o := range opts {
		o(game)
	}

	if game.IsIncomingDeck {
		// randomly determine number between 1 and 17 and
		// covert to percent.  use percentage to figure out
		// how many cards must be dealt before new incoming deck
		max := 0.17
		min := 0.01
		random := min + rand.Float64()*(max-min)
		count := len(game.Shoe.Cards)
		fcount := float64(count)
		val := int(fcount * random)

		game.IncomingDeckPosition = count - val
	}

	game.Dealer = &Player{
		Hands: []*Hand{
			{
				Id: 1,
			},
		},
	}

	return game, nil
}

func RunCLI() {

	g, err := NewBlackjackGame()
	if err != nil {
		fmt.Println(fmt.Errorf("cannot create new blackjack game, %s", err))
	}

	g.PlayerSetup(g.output, g.input)

	for g.PlayAgain() {

		g.ResetPlayers()
		g.Betting()
		g.Players = g.RemoveQuitPlayers()
		if len(g.Players) == 0 {
			break
		}
		g.Start()
	}

	fmt.Println("No players left in game.  Exiting...")
}

func (g *Game) AddPlayer(player *Player) {
	g.Players = append(g.Players, player)
}

func (g *Game) PlayerSetup(output io.Writer, input io.Reader) error {

	var answer int

	for answer < 1 {
		fmt.Fprintln(output, "Please enter number of Blackjack players:")
		fmt.Fscanln(input, &answer)
	}

	for i := 1; i <= answer; i++ {

		player, err := NewPlayer(output, input)
		if err != nil {
			return fmt.Errorf("unable to setup players, %s", err)
		}

		g.AddPlayer(&player)
	}

	return nil
}

func (g *Game) PlayAgain() bool {

	for _, player := range g.Players {
		if player.Action != ActionQuit {
			return true
		}
	}
	return false
}

func (g *Game) Betting() error {

	var err error
	for index, player := range g.Players {

		err = player.Bet(g.output, g.input, player, index)
		if err != nil {
			return fmt.Errorf("unable to place bet for player: %s", player.Name)
		}

	}
	return nil

}

func (g *Game) Start() {

	fmt.Fprintln(g.output, "")
	fmt.Fprintln(g.output, "****** NEW GAME ******")

	g.OpeningDeal()

	fmt.Fprintln(g.output, g.Dealer.DealerString())
	g.ShowPlayerCards(g.output)

	for _, player := range g.Players {
		g.PlayHand(player)
	}

	fmt.Fprintln(g.output, "")
	g.DealerPlay()

	fmt.Fprintln(g.output, "Dealer"+g.Dealer.PlayerString())
	g.ShowPlayerCards(g.output)

	fmt.Fprintln(g.output, "")

	g.Outcome(g.output)

}

func (g *Game) DealerPlay() {

	dealerOk := g.IsDealerDraw()

	if dealerOk {
		fmt.Fprintln(g.output, "****** DEALER'S TURN ******")

		for g.Dealer.Hands[0].Score() <= 16 || (g.Dealer.Hands[0].Score() == 17 && g.Dealer.Hands[0].MinScore() != 17) {
			card := g.Deal(g.output)
			g.Dealer.Hands[0].Cards = append(g.Dealer.Hands[0].Cards, card)
		}
		g.Dealer.Hands[0].Action = ActionStand
	}
}

func (g Game) IsDealerDraw() bool {

	result := false
	allNotBustOrBlackjack := false

	for _, player := range g.Players {
		for _, hand := range player.Hands {
			if hand.Outcome != OutcomeBust && hand.Outcome != OutcomeBlackjack {
				allNotBustOrBlackjack = true
			}
		}
	}

	// verifies if game conditions warrant dealer drawing a card
	if g.Dealer.Hands[0].Score() <= 16 || (g.Dealer.Hands[0].Score() == 17 && g.Dealer.Hands[0].MinScore() != 17) {
		if allNotBustOrBlackjack {
			result = true
		}
	}

	return result
}

func (g Game) ShowPlayerCards(output io.Writer) {
	for _, player := range g.Players {
		fmt.Fprintln(output, player.PlayerString())
	}
}

func (g *Game) Deal(output io.Writer) cards.Card {

	var card cards.Card
	card, g.Shoe.Cards = g.Shoe.Cards[0], g.Shoe.Cards[1:]
	g.CardsDealt += 1

	if g.IsIncomingDeck {
		if g.CardsDealt >= g.IncomingDeckPosition {
			fmt.Fprintln(output, "\n******************************")
			fmt.Fprintln(output, "***** New Deck Incoming ******")
			fmt.Fprintln(output, "\n******************************")
			g.Shoe = g.IncomingDeck()
			g.CardsDealt = 0
		}

	} else {
		g.Shoe.Cards = append(g.Shoe.Cards, card)
	}
	return card
}

func (g *Game) OpeningDeal() {
	for i := 0; i < 2; i++ {
		for _, player := range g.Players {
			card := g.Deal(g.output)
			player.Hands[player.HandIndex].Cards = append(player.Hands[player.HandIndex].Cards, card)
		}
		card := g.Deal(g.output)
		g.Dealer.Hands[g.Dealer.HandIndex].Cards = append(g.Dealer.Hands[g.Dealer.HandIndex].Cards, card)
	}
}

func (g *Game) Outcome(output io.Writer) {

	var outcome Outcome
	for _, player := range g.Players {
		for _, hand := range player.Hands {
			if hand.Outcome == OutcomeBlackjack || hand.Outcome == OutcomeBust {
				outcome = hand.Outcome
			} else if g.Dealer.Hands[0].Score() > 21 {
				outcome = OutcomeWin
			} else if hand.Score() > g.Dealer.Hands[0].Score() {
				outcome = OutcomeWin
			} else if hand.Score() < g.Dealer.Hands[0].Score() {
				outcome = OutcomeLose
			} else if hand.Score() == g.Dealer.Hands[0].Score() {
				outcome = OutcomeTie
			}

			hand.Outcome = outcome

		}

		player.SetWinLoseTie()

		player.Payout()

		player.Broke()

		player.OutcomeReport(output)
	}

	g.RemoveQuitPlayers()
}

func (g *Game) ResetPlayers() {

	for _, player := range g.Players {
		player.HandIndex = 0
		player.Hands = []*Hand{}
		hand := NewHand(1)
		player.AddHand(hand)
		player.Action = None

	}
	g.Dealer.Hands = []*Hand{}
	hand := NewHand(1)
	g.Dealer.AddHand(hand)

}

func (g *Game) RemoveQuitPlayers() []*Player {
	newPlayers := []*Player{}
	if len(g.Players) > 0 {
		for _, player := range g.Players {
			if player.Action != ActionQuit {
				newPlayers = append(newPlayers, player)
			}
		}
	}

	return newPlayers
}

func (g Game) IncomingDeck() cards.Deck {
	deck := cards.NewDeck(
		cards.WithNumberOfDecks(g.DeckCount),
	)

	return deck

}

type Player struct {
	Name          string
	Action        Action
	Bet           func(io.Writer, io.Reader, *Player, int) error
	Decide        func(io.Writer, io.Reader, *Player, cards.Card, int) Action
	AiHandsToPlay int
	Record        Record
	Cash          int
	Hands         []*Hand
	HandIndex     int
}

func (p *Player) Payout() {

	for _, hand := range p.Hands {
		if hand.Outcome == OutcomeWin {
			hand.Payout = hand.Bet
			p.Cash += hand.Bet + hand.Payout
			hand.Bet = 0
		} else if hand.Outcome == OutcomeLose || hand.Outcome == OutcomeBust {
			hand.Payout = -1 * hand.Bet
			hand.Bet = 0
		} else if hand.Outcome == OutcomeTie {
			hand.Payout = 0
			p.Cash += hand.Bet
			hand.Bet = 0
		} else if hand.Outcome == OutcomeBlackjack {
			hand.Payout = 2 * hand.Bet
			p.Cash += hand.Bet + hand.Payout
			hand.Bet = 0
		}
	}
}

func (p *Player) Broke() {
	if p.Cash == 0 {
		p.Action = ActionQuit
	}

}

func (p *Player) Split(output io.Writer, card1, card2 cards.Card) {

	id := p.NextHandId()
	hand := NewHand(id)
	p.AddHand(hand)
	indexNewHand := len(p.Hands) - 1

	// take last card in original hand and append to the new split hand
	card := p.Hands[p.HandIndex].Cards[1]
	p.Hands[indexNewHand].Cards = append(p.Hands[indexNewHand].Cards, card)

	// reset slice on original hand to only have the first card
	p.Hands[p.HandIndex].Cards = p.Hands[p.HandIndex].Cards[:len(p.Hands[p.HandIndex].Cards)-1]

	// mirror bet on new hand
	p.Hands[indexNewHand].Bet += p.Hands[p.HandIndex].Bet
	p.Cash -= p.Hands[p.HandIndex].Bet

	// add cards to each hand
	p.Hands[p.HandIndex].Cards = append(p.Hands[p.HandIndex].Cards, card1)
	p.Hands[indexNewHand].Cards = append(p.Hands[indexNewHand].Cards, card2)
	p.Hands[p.HandIndex].Action = None
	p.Hands[indexNewHand].Action = None

	fmt.Fprintln(output, p.Hands[p.HandIndex].HandString(p.Name))
	fmt.Fprintln(output, p.Hands[indexNewHand].HandString(p.Name))

}

func (p *Player) AddHand(hand *Hand) {
	p.Hands = append(p.Hands, hand)
}

func NewHand(id int) *Hand {
	hand := Hand{
		Id: id,
	}
	return &hand
}

func (p *Player) SetWinLoseTie() {

	for _, hand := range p.Hands {
		if hand.Outcome == OutcomeWin || hand.Outcome == OutcomeBlackjack {
			p.Record.Win += 1
		} else if hand.Outcome == OutcomeTie {
			p.Record.Tie += 1
		} else {
			p.Record.Lose += 1
		}
		p.Record.HandsPlayed += 1
	}
}

func (p *Player) OutcomeReport(output io.Writer) {

	var payout string

	for _, hand := range p.Hands {
		if hand.Outcome == OutcomeTie {
			payout = ""
		} else {
			payout = strconv.Itoa(absInt(hand.Payout))
		}

		str := []string{
			p.Name,
			BalanceReportMap[hand.Outcome],
			payout,
			".  Cash available: $",
			strconv.Itoa(p.Cash),
		}

		fmt.Fprintln(output, strings.Join(str, ""))
	}
}

func (p Player) DealerString() string {

	return "Dealer has: [" + p.Hands[0].Cards[0].String() + "]" + "[???]"
}

func (p Player) PlayerString() string {

	builder := strings.Builder{}
	var response string
	for index, hand := range p.Hands {
		if hand.Action == ActionDoubleDown {
			builder.WriteString(p.Name + " has ???: " + "[" + hand.Cards[0].String() + "]" + "[" + hand.Cards[1].String() + "]" + "[???]\n")
			response += builder.String()
		} else {
			for _, card := range p.Hands[index].Cards {
				builder.WriteString("[" + card.String() + "]")
			}
			str := []string{p.Name, " has ", fmt.Sprint(hand.Score()), ": ", builder.String(), "\n"}
			response = strings.Join(str, "")
		}

	}

	return response
}

func NewPlayer(output io.Writer, input io.Reader) (Player, error) {
	var name string
	var playerTypeInput string
	var aiHandsInput string

	fmt.Fprintln(output, "Enter your name: ")
	fmt.Fscanln(input, &name)

	for strings.ToLower(playerTypeInput) != "h" && strings.ToLower(playerTypeInput) != "a" {
		fmt.Fprintln(output, "Select (H)uman or (A)i")
		fmt.Fscanln(input, &playerTypeInput)
	}

	if strings.ToLower(playerTypeInput) == "a" {
		for strings.ToLower(playerTypeInput) != "b" && strings.ToLower(playerTypeInput) != "s" {
			fmt.Fprintln(output, "Select AI Type: (B)asic or (S)tandOnly")
			fmt.Fscanln(input, &playerTypeInput)
		}
	}

	playerTypeInputValue := PlayerTypeInputMap[strings.ToLower(playerTypeInput)]
	playerType := PlayerTypeMap[playerTypeInputValue]
	playerTypeBet := PlayerTypeBetMap[playerTypeInputValue]

	aiHands := 0
	var err error
	if strings.ToLower(playerTypeInput) != "h" {
		fmt.Fprintln(output, "Enter number of hands the AI plays: ")
		fmt.Fscanln(input, &aiHandsInput)
		aiHands, err = strconv.Atoi(aiHandsInput)
		if err != nil {
			return Player{}, fmt.Errorf("unable to create new player,%s", err)
		}
	}

	player := Player{

		Name:          name,
		Decide:        playerType,
		AiHandsToPlay: aiHands,
		Bet:           playerTypeBet,
		Cash:          100,
	}

	return player, nil
}

type Hand struct {
	Id      int
	Cards   []cards.Card
	Bet     int
	Action  Action
	Outcome Outcome
	Payout  int
}

func (h *Hand) Hit(output io.Writer, card cards.Card, name string) {

	h.Cards = append(h.Cards, card)
	h.Action = None

	if h.Score() > 21 {
		h.Outcome = OutcomeBust
	}

}

func (h *Hand) DoubleDown(output io.Writer, card cards.Card, name string) {

	h.Bet += h.Bet
	h.Cards = append(h.Cards, card)
	fmt.Fprintln(output, h.HandString(name))
	h.Action = ActionStand

}

func (p Player) NextHandId() int {
	return len(p.Hands) + 1
}

func (h Hand) ChooseAction() bool {
	return h.Action != ActionQuit && h.Action != ActionStand && h.Outcome != OutcomeBlackjack && h.Outcome != OutcomeBust
}

func (h Hand) HandString(name string) string {

	builder := strings.Builder{}
	var response string

	if h.Action == ActionDoubleDown {
		builder.WriteString(name + " hand #" + strconv.Itoa(h.Id) + " has ???: " + "[" + h.Cards[0].String() + "]" + "[" + h.Cards[1].String() + "]" + "[???]\n")
		response += builder.String()
	} else {
		for _, card := range h.Cards {
			builder.WriteString("[" + card.String() + "]")
		}
		str := []string{name, " hand #", strconv.Itoa(h.Id), " has ", fmt.Sprint(h.Score()), ": ", builder.String(), "\n"}
		response = strings.Join(str, "")
	}

	return response
}

func (h Hand) Score() int {
	minScore := h.MinScore()

	if minScore > 11 {
		return minScore
	}
	for _, c := range h.Cards {
		if c.Rank == cards.Ace {
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

type Record struct {
	Win         int
	Lose        int
	Tie         int
	HandsPlayed int
}

func (r Record) RecordString() string {

	str := []string{
		"************** Player Win-Lose-Tie Report **************\nPlayer won: ",
		strconv.Itoa(r.Win),
		", lost: ",
		strconv.Itoa(r.Lose),
		" and tied: ",
		strconv.Itoa(r.Tie),
		"\n",
	}

	return strings.Join(str, "")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func absInt(x int) int {
	return absDiffInt(x, 0)
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func ScoreDealerHoleCard(card cards.Card) int {

	var score int
	min := min(int(card.Rank), 10)

	if card.Rank == cards.Ace {
		score = min + 10
	} else {
		score = min
	}

	return score
}

func HumanAction(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card, index int) Action {

	var answer string

	// check to see if not enough to split or double
	if player.Hands[index].Bet > player.Cash {
		for strings.ToLower(answer) != "h" && strings.ToLower(answer) != "s" {
			fmt.Fprintln(output, "Please choose (H)it or (S)tand")
			fmt.Fscanln(input, &answer)
		}
	} else {
		// check if split is ok
		if player.Hands[index].Cards[0].Rank == player.Hands[index].Cards[1].Rank {
			for strings.ToLower(answer) != "h" && strings.ToLower(answer) != "p" && strings.ToLower(answer) != "d" && strings.ToLower(answer) != "s" {
				fmt.Fprintln(output, "Please choose (H)it, S(P)lit, (D)ouble or (S)tand")
				fmt.Fscanln(input, &answer)
			}
		} else {
			// double ok
			for strings.ToLower(answer) != "h" && strings.ToLower(answer) != "d" && strings.ToLower(answer) != "s" {
				fmt.Fprintln(output, "Please choose (H)it, (D)ouble or (S)tand")
				fmt.Fscanln(input, &answer)
			}
		}
	}

	return ActionMap[strings.ToLower(answer)]
}

func HumanBet(output io.Writer, input io.Reader, player *Player, index int) error {

	var answer string
	var bet int = 0

	fmt.Fprintln(output, "")
	fmt.Fprintln(output, "****** BET or QUIT ******")

	for answer != "b" && answer != "q" {
		fmt.Fprintln(output, "Enter (b)et or (q)uit:")
		fmt.Fscanln(input, &answer)
	}

	switch strings.ToLower(answer) {
	case "q":
		player.Action = ActionQuit
	case "b":
		for bet < 1 || bet > player.Cash {
			fmt.Fprintf(output, "Enter bet amount between 1 and %d: ", player.Cash)
			fmt.Fscanln(input, &bet)
		}
		player.Cash -= bet
		player.Hands[player.HandIndex].Bet += bet

	}

	return nil

}

func (g *Game) PlayHand(player *Player) {
	for index, hand := range player.Hands {
		if hand.Score() == 21 {
			hand.Outcome = OutcomeBlackjack
		}

		for hand.ChooseAction() {
			if hand.Action == None {
				fmt.Fprintln(g.output, hand.HandString(player.Name))
				hand.Action = player.Decide(g.output, g.input, player, g.Dealer.Hands[0].Cards[0], index)
			}

			if hand.Action == ActionHit {
				card := g.Deal(g.output)
				hand.Hit(g.output, card, player.Name)
			} else if hand.Action == ActionDoubleDown {
				player.Cash -= hand.Bet
				card := g.Deal(g.output)
				hand.DoubleDown(g.output, card, player.Name)
			} else if hand.Action == ActionSplit {
				card1 := g.Deal(g.output)
				card2 := g.Deal(g.output)
				player.Split(g.output, card1, card2)
				g.PlayHand(player)
			}
		}
	}

}

// additional features
// difficult to fake concreate...take interface instead
// burn cards

// 17. ui
// 16. client/server
// 15. card counting ai
// 14. ai betting - inc or dec depending on last outcome
// 13. betting limits (max - min)
// 12. card counting
// 11. split
// 10. Double down - done
// 9. basic strategy ai - done
// 8. reshuffle - done
// 7. betting - push outcome...just bet $10, handle no money player and dealer - done
// 6. multi players - done
// 5. deal - done
// 4. player - done
// 3. game - done
// 2. deck - done
// 1. card - done
