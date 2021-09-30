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
	ActionHit:   "Hit",
	ActionStand: "Stand",
	ActionQuit:  "Quit",
	None:        "Invalid Action",
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
)

var ActionMap = map[string]Action{
	"h": ActionHit,
	"s": ActionStand,
	"q": ActionQuit,
	"n": None,
	"d": ActionDoubleDown,
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

var PlayerTypeMap = map[PlayerType]func(io.Writer, io.Reader, *Player, cards.Card) Action{
	PlayerTypeHuman:       HumanAction,
	PlayerTypeAiStandOnly: AiActionStandOnly,
	PlayerTypeAiBasic:     AiActionBasic,
}

var PlayerTypeInputMap = map[string]PlayerType{
	"h": PlayerTypeHuman,
	"b": PlayerTypeAiBasic,
	"s": PlayerTypeAiStandOnly,
}

var PlayerTypeBetMap = map[PlayerType]func(io.Writer, io.Reader, *Player) error{
	PlayerTypeHuman:       HumanBet,
	PlayerTypeAiStandOnly: AiBet,
	PlayerTypeAiBasic:     AiBet,
}

type Game struct {
	Players              []*Player
	Dealer               Player
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
		DeckCount:      3,
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
	for index := range g.Players {

		err = g.Players[index].Bet(g.output, g.input, g.Players[index])
		if err != nil {
			return fmt.Errorf("unable to place bet for player: %s", g.Players[index].Name)
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

	for index := range g.Players {

		if g.Players[index].Score() == 21 {
			g.Players[index].HandOutcome = OutcomeBlackjack
		}

		for g.Players[index].ChooseAction() {
			if g.Players[index].Action == None {
				g.Players[index].Action = g.Players[index].Decide(g.output, g.input, g.Players[index], g.Dealer.Hand[0])
			}
			if g.Players[index].Action == ActionHit {

				card := g.Deal(g.output)
				g.Players[index].Hit(g.output, card)
			} else if g.Players[index].Action == ActionDoubleDown {
				card := g.Deal(g.output)
				g.Players[index].DoubleDown(g.output, card)

			}
		}
	}

	fmt.Fprintln(g.output, "")
	g.DealerPlay()

	fmt.Fprintln(g.output, "Dealer"+g.Dealer.PlayerString())
	g.ShowPlayerCards(g.output)

	fmt.Fprintln(g.output, "")

	g.Outcome(g.output)

}

func (g *Game) DealerPlay() {
	fmt.Fprintln(g.output, "****** DEALER'S TURN ******")

	for g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {
		card := g.Deal(g.output)
		g.Dealer.Hand = append(g.Dealer.Hand, card)
	}
	g.Dealer.Action = ActionStand
}

func (g Game) IsDealerDraw() bool {

	result := false
	allNotBustOrBlackjack := false

	for index := range g.Players {
		if g.Players[index].HandOutcome != OutcomeBust && g.Players[index].HandOutcome != OutcomeBlackjack {
			allNotBustOrBlackjack = true
		}
	}

	// verifies if game conditions warrant dealer drawing a card
	if g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {
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
		for index := range g.Players {
			card := g.Deal(g.output)
			g.Players[index].Hand = append(g.Players[index].Hand, card)
		}
		card := g.Deal(g.output)
		g.Dealer.Hand = append(g.Dealer.Hand, card)
	}
}

func (g *Game) Outcome(output io.Writer) {

	var outcome Outcome
	for index := range g.Players {
		if g.Players[index].HandOutcome == OutcomeBlackjack || g.Players[index].HandOutcome == OutcomeBust {
			outcome = g.Players[index].HandOutcome
		} else if g.Dealer.Score() > 21 {
			outcome = OutcomeWin
		} else if g.Players[index].Score() > g.Dealer.Score() {
			outcome = OutcomeWin
		} else if g.Players[index].Score() < g.Dealer.Score() {
			outcome = OutcomeLose
		} else {
			outcome = OutcomeTie
		}

		g.Players[index].HandOutcome = outcome

		g.Players[index].SetWinLoseTie(outcome)

		g.Players[index].Payout()

		g.Players[index].Broke()

		g.Players[index].OutcomeReport(output)

	}

	g.RemoveQuitPlayers()
}

func (g *Game) ResetPlayers() {

	for index := range g.Players {

		g.Players[index].Hand = []cards.Card{}
		g.Players[index].Action = None
		g.Players[index].HandOutcome = OutcomeNone
		g.Players[index].HandBet = 0
		g.Players[index].HandPayout = 0
	}
	g.Dealer.Hand = []cards.Card{}

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
	Hand          []cards.Card
	Action        Action
	HandOutcome   Outcome
	Bet           func(io.Writer, io.Reader, *Player) error
	Decide        func(io.Writer, io.Reader, *Player, cards.Card) Action
	AiHandsToPlay int
	Record        Record
	Cash          int
	HandBet       int
	HandPayout    int
	Hands         []Hand
	HandIndex     int
}

func (p *Player) Payout() {

	if p.HandOutcome == OutcomeWin {
		p.HandPayout = p.HandBet
		p.Cash += p.HandBet + p.HandPayout
		p.HandBet = 0
	} else if p.HandOutcome == OutcomeLose || p.HandOutcome == OutcomeBust {
		p.HandPayout = -1 * p.HandBet
		p.HandBet = 0
	} else if p.HandOutcome == OutcomeTie {
		p.HandPayout = 0
		p.Cash += p.HandBet
		p.HandBet = 0
	} else if p.HandOutcome == OutcomeBlackjack {
		p.HandPayout = 2 * p.HandBet
		p.Cash += p.HandBet + p.HandPayout
		p.HandBet = 0
	}
}

func (p *Player) Broke() {
	if p.Cash == 0 {
		p.Action = ActionQuit
	}

}

func (p *Player) Hit(output io.Writer, card cards.Card) {

	p.Hand = append(p.Hand, card)
	p.Action = None

	fmt.Fprintln(output, p.PlayerString())
	if p.Score() > 21 {
		p.HandOutcome = OutcomeBust
	}

}

func (p *Player) DoubleDown(output io.Writer, card cards.Card) {

	p.HandBet += p.HandBet
	p.Hand = append(p.Hand, card)
	fmt.Fprintln(output, p.PlayerString())
	p.Action = ActionStand

}

func (p Player) ChooseAction() bool {

	return p.Action != ActionQuit && p.Action != ActionStand && p.HandOutcome != OutcomeBlackjack && p.HandOutcome != OutcomeBust
}

func (p *Player) SetWinLoseTie(outcome Outcome) {

	if outcome == OutcomeWin || outcome == OutcomeBlackjack {
		p.Record.Win += 1
	} else if outcome == OutcomeTie {
		p.Record.Tie += 1
	} else {
		p.Record.Lose += 1
	}

	p.Record.HandsPlayed += 1

}

func (p Player) PlayerString() string {

	builder := strings.Builder{}
	var response string
	if p.Action == ActionDoubleDown {
		builder.WriteString(p.Name + " has ???: " + "[" + p.Hand[0].String() + "]" + "[" + p.Hand[1].String() + "]" + "[???]")
		response = builder.String()
	} else {
		for _, card := range p.Hand {
			builder.WriteString("[" + card.String() + "]")
		}
		str := []string{p.Name, " has ", fmt.Sprint(p.Score()), ": ", builder.String()}
		response = strings.Join(str, "")
	}
	return response
}

func (p *Player) OutcomeReport(output io.Writer) {

	var payout string
	if p.HandOutcome == OutcomeTie {
		payout = ""
	} else {
		payout = strconv.Itoa(absInt(p.HandPayout))
	}

	str := []string{
		p.Name,
		BalanceReportMap[p.HandOutcome],
		payout,
		".  Cash available: $",
		strconv.Itoa(p.Cash),
	}

	fmt.Fprintln(output, strings.Join(str, ""))
}

func (p Player) DealerString() string {

	return "Dealer has: [" + p.Hand[0].String() + "]" + "[???]"
}

func (p Player) Score() int {
	minScore := p.MinScore()

	if minScore > 11 {
		return minScore
	}
	for _, c := range p.Hand {
		if c.Rank == cards.Ace {
			return minScore + 10
		}
	}
	return minScore
}

func (p Player) MinScore() int {
	score := 0
	for _, c := range p.Hand {
		score += min(int(c.Rank), 10)
	}
	return score
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

type Hand struct {
	Id      int
	Cards   []cards.Card
	Bet     int
	Action  Action
	Outcome Outcome
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

func HumanAction(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card) Action {

	var answer string

	// check if double is possible
	if player.HandBet > player.Cash {
		for strings.ToLower(answer) != "h" && strings.ToLower(answer) != "s" {
			fmt.Fprintln(output, "Please choose (H)it or (S)tand")
			fmt.Fscanln(input, &answer)
		}
	} else {
		for strings.ToLower(answer) != "h" && strings.ToLower(answer) != "d" && strings.ToLower(answer) != "s" {
			fmt.Fprintln(output, "Please choose (H)it, (D)ouble or (S)tand")
			fmt.Fscanln(input, &answer)
		}
	}

	return ActionMap[strings.ToLower(answer)]
}

func HumanBet(output io.Writer, input io.Reader, player *Player) error {

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
		player.HandBet += bet
	}

	return nil

}

func (p Player) NextHandId() int {
	return len(p.Hands) + 1
}

func (p Player) NewHand() Hand {
	id := p.NextHandId()

	hand := Hand{
		Id: id,
	}
	return hand
}

func (p *Player) AddHand(hand Hand) {
	p.Hands = append(p.Hands, hand)

}

func (p *Player) Split(card1, card2 cards.Card) {
	hand := p.NewHand()
	p.AddHand(hand)
	indexNewHand := len(p.Hands) - 1

	// take last card in original hand and append to the new split hand
	// reset slice on original hand to only have the first card
	var card cards.Card
	card, p.Hands[p.HandIndex].Cards = p.Hands[p.HandIndex].Cards[1], p.Hands[p.HandIndex].Cards[0:0]
	p.Hands[indexNewHand].Cards = append(p.Hands[indexNewHand].Cards, card)

	p.Hands[indexNewHand].Bet = p.Hands[p.HandIndex].Bet
	p.Cash -= p.Hands[p.HandIndex].Bet

	p.Hands[indexNewHand].Cards = append(p.Hands[indexNewHand].Cards, card1)
	p.Hands[indexNewHand].Cards = append(p.Hands[indexNewHand].Cards, card2)

}

// additional features
// difficult to fake concreate...take interface instead
// burn cards

// 7. split
// 6. Double down
// 5. card counting
// 4. basic strategy ai
// 3. reshuffle
// 2. betting - push outcome...just bet $10, handle no money player and dealer
// 1. multi players
// multi hands for split

// fold
// card counter
