package blackjack

import (
	"cards"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
)

var ActionMap = map[string]Action{
	"h": ActionHit,
	"s": ActionStand,
	"q": ActionQuit,
	"n": None,
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

type PlayerType int

const (
	PlayerTypeHuman PlayerType = iota
	PlayerTypeAiStandOnly
)

var PlayerTypeMap = map[PlayerType]func(io.Writer, io.Reader) Action{
	PlayerTypeHuman:       GetHumanAction,
	PlayerTypeAiStandOnly: GetAiActionStandOnly,
}

var PlayerTypeInputMap = map[string]PlayerType{
	"h": PlayerTypeHuman,
	"a": PlayerTypeAiStandOnly,
}

var PlayerTypeBetMap = map[PlayerType]func(io.Writer, io.Reader, Player) Player{
	PlayerTypeHuman:       HumanBet,
	PlayerTypeAiStandOnly: AiBet,
}

type Game struct {
	Players []Player
	Dealer  Player
	Shoe    cards.Deck
	Random  rand.Rand
	output  io.Writer
	input   io.Reader
}

type Option func(*Game) error

func WithCustomDeck(deck cards.Deck) Option {
	return func(g *Game) error {
		g.Shoe = deck
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

func NewBlackjackGame(opts ...Option) (*Game, error) {

	deck := cards.NewDeck(
		cards.WithNumberOfDecks(3),
	)

	game := &Game{
		Shoe:   deck,
		output: os.Stdout,
		input:  os.Stdin,
	}

	for _, o := range opts {
		o(game)
	}

	return game, nil
}

func RunCLI() error {

	g, err := NewBlackjackGame()
	if err != nil {
		return fmt.Errorf("cannot create new blackjack game, %s", err)
	}

	g.PlayerSetup(g.output, g.input)

	for g.Continue() {
		g.Betting()
		g.Players = g.RemoveQuitPlayers()
		if len(g.Players) == 0 {
			RunCLI()
		}
		g.ResetPlayers()
		g.Start()
	}

	return nil
}

func (g *Game) AddPlayer(player Player) {
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

		g.AddPlayer(player)
	}

	return nil
}

func (g Game) Continue() bool {

	for _, player := range g.Players {
		if player.Action != ActionQuit {
			return true
		}
	}
	return false
}

func (g *Game) Betting() {

	for index := range g.Players {
		g.Players[index] = g.Players[index].Bet(g.output, g.input, g.Players[index])
	}

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

		for g.Players[index].Continue() {
			if g.Players[index].Action == None {
				g.Players[index].Action = g.Players[index].Decide(g.output, g.input)
			}
			if g.Players[index].Action == ActionHit {

				card := g.Deal()
				g.Players[index].Hand = append(g.Players[index].Hand, card)
				g.Players[index].Action = None

				fmt.Fprintln(g.output, g.Players[index].Name+" has "+g.Players[index].PlayerString())
				if g.Players[index].Score() > 21 {
					g.Players[index].HandOutcome = OutcomeBust
				}
			}
		}
	}

	fmt.Fprintln(g.output, "")
	g.DealerPlay()

	fmt.Fprintln(g.output, "Dealer has "+g.Dealer.PlayerString())
	g.ShowPlayerCards(g.output)
	fmt.Fprintln(g.output, "")

	g.Outcome(g.output)

}

func (g *Game) DealerPlay() {
	fmt.Fprintln(g.output, "****** DEALER'S TURN ******")

	for g.Dealer.Score() <= 16 || (g.Dealer.Score() == 17 && g.Dealer.MinScore() != 17) {
		card := g.Deal()
		g.Dealer.Hand = append(g.Dealer.Hand, card)
	}
	g.Dealer.Action = ActionStand
}

func (g Game) ShowPlayerCards(output io.Writer) {
	for _, player := range g.Players {
		fmt.Fprintln(output, player.Name+" has "+player.PlayerString())
	}
}

func (g *Game) Deal() cards.Card {

	var card cards.Card
	card, g.Shoe.Cards = g.Shoe.Cards[0], g.Shoe.Cards[1:]
	g.Shoe.Cards = append(g.Shoe.Cards, card)
	return card
}

func (g *Game) OpeningDeal() {
	for i := 0; i < 2; i++ {
		for index := range g.Players {
			card := g.Deal()
			g.Players[index].Hand = append(g.Players[index].Hand, card)
		}
		card := g.Deal()
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

		g.Players[index].SetPlayerWinLoseTie(outcome)

		if g.Players[index].Logic == "a" && g.Players[index].AiHandsToPlay == g.Players[index].Record.HandsPlayed {

			g.Players[index].Action = ActionQuit
		}

	}
}

func (g *Game) ResetPlayers() {

	for index := range g.Players {

		g.Players[index].Hand = []cards.Card{}
		g.Players[index].Action = None
		g.Players[index].HandOutcome = OutcomeNone
	}
	g.Dealer.Hand = []cards.Card{}

}

func (g *Game) RemoveQuitPlayers() []Player {
	newPlayers := []Player{}

	for _, player := range g.Players {
		if player.Action != ActionQuit {
			newPlayers = append(newPlayers, player)
		}
	}

	return newPlayers
}

type Player struct {
	Name          string
	Hand          []cards.Card
	Action        Action
	HandOutcome   Outcome
	Bet           func(io.Writer, io.Reader, Player) Player
	Decide        func(io.Writer, io.Reader) Action
	AiHandsToPlay int
	Record        Record
	Logic         string
}

func (p Player) Continue() bool {

	return p.Action != ActionQuit && p.Action != ActionStand && p.HandOutcome != OutcomeBlackjack && p.HandOutcome != OutcomeBust
}

func (p *Player) SetPlayerWinLoseTie(outcome Outcome) {

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
	for _, card := range p.Hand {
		builder.WriteString("[" + card.String() + "]")
	}
	return fmt.Sprint(p.Score()) + ": " + builder.String()
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

type Record struct {
	Win         int
	Lose        int
	Tie         int
	HandsPlayed int
}

func (r Record) String() string {

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

func NewPlayer(output io.Writer, input io.Reader) (Player, error) {
	var name string
	var playerTypeInput string
	var aiHandsInput string

	fmt.Fprintln(output, "Enter your name: ")
	fmt.Fscanln(input, &name)

	fmt.Fprintln(output, "Select (H)uman or (A)i for player type: ")
	fmt.Fscanln(input, &playerTypeInput)

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
		Logic:         strings.ToLower(playerTypeInput),
		Bet:           playerTypeBet,
	}

	return player, nil
}

func GetHumanAction(output io.Writer, input io.Reader) Action {

	var answer string

	fmt.Fprintln(output, "Please choose (H)it or (S)tand")
	fmt.Fscanln(input, &answer)

	return ActionMap[strings.ToLower(answer)]
}

func GetAiActionStandOnly(output io.Writer, input io.Reader) Action {

	return ActionStand
}

func HumanBet(output io.Writer, input io.Reader, player Player) Player {

	var answer string

	fmt.Fprintln(output, "")
	fmt.Fprintln(output, "****** BET or QUIT ******")
	fmt.Fprintln(output, "Enter (b)et or (q)uit:")

	fmt.Fscanln(input, &answer)

	if strings.ToLower(answer) == "q" {
		player.Action = ActionQuit
	}

	return player

}

func AiBet(output io.Writer, input io.Reader, player Player) Player {

	if player.Record.HandsPlayed == player.AiHandsToPlay {
		player.Action = ActionQuit
	}
	return player
}

// additional features
// difficult to fake concreate...take interface instead
// burn cards
// 4. basic strategy ai
// 3. reshuffle
// 2. betting - push outcome...just bet $10, handle no money player and dealer
// 1. multi players
// multi hands for split

// fold
// card counter
