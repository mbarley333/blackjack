package blackjack

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mbarley333/cards"
)

type Action int

var ActionStringMap = map[Action]string{
	None:             "Invalid Action",
	ActionHit:        "Hit",
	ActionStand:      "Stand",
	ActionQuit:       "Quit",
	ActionDoubleDown: "Double Down",
	ActionSplit:      "Split",
	ActionBet:        "Bet",
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
	ActionBet
)

var ActionMap = map[string]Action{
	"h": ActionHit,
	"s": ActionStand,
	"q": ActionQuit,
	"n": None,
	"d": ActionDoubleDown,
	"p": ActionSplit,
	"b": ActionBet,
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
	PlayerTypeAiCustom
)

var PlayerTypeMap = map[PlayerType]func(io.Writer, io.Reader, *Player, cards.Card, int, CardCounter, Stage) Action{
	PlayerTypeHuman:       HumanAction,
	PlayerTypeAiStandOnly: AiActionStandOnly,
	PlayerTypeAiBasic:     AiActionBasic,
	PlayerTypeAiCustom:    AiActionBasic,
}

var PlayerTypeInputMap = map[string]PlayerType{
	"h": PlayerTypeHuman,
	"b": PlayerTypeAiBasic,
	"s": PlayerTypeAiStandOnly,
	"x": PlayerTypeAiBasic,
}

var PlayerTypeBetMap = map[PlayerType]func(*Game) error{
	PlayerTypeHuman:       HumanBet,
	PlayerTypeAiStandOnly: AiBet,
	PlayerTypeAiBasic:     AiBet,
	PlayerTypeAiCustom:    AiBet,
}

type Stage int

const (
	StageStart Stage = iota
	StageOpeningDeal
	StageBetting
	StageDeciding
	StageDealerPlay
	StageOutcome
)

var StageMap = map[Stage]string{
	StageStart:       "Start",
	StageOpeningDeal: "Opening Deal",
	StageBetting:     "Betting",
	StageDeciding:    "Deciding",
	StageDealerPlay:  "Dealer Play",
	StageOutcome:     "Outcome",
}

var StageDisplayMessageMap = map[Stage]string{
	StageStart:       "LET'S PLAY BLACKJACK!",
	StageOpeningDeal: "OPENING DEAL",
	StageBetting:     "PLACE YOUR BETS!",
	StageDeciding:    "PLAYERS MAKE YOUR CHOICE",
	StageDealerPlay:  "DEALER PLAY",
	StageOutcome:     "OUTCOME",
}

func (s Stage) String() string {
	return StageMap[s]
}

func (s Stage) Message() string {
	return StageDisplayMessageMap[s]
}

type Dialog int

const (
	DialogNone Dialog = iota
	DialogBetOrQuit
	DialogPlaceYourBet
	DialogHitOrStand
	DialogHitSplitDoubleStand
	DialogHitDoubleStand
)

var DialogMap = map[Dialog]string{
	DialogNone:                "Invalid Dialog",
	DialogBetOrQuit:           "BetOrQuit",
	DialogPlaceYourBet:        "PlaceYourBet",
	DialogHitOrStand:          "HitOrStand",
	DialogHitSplitDoubleStand: "HitSplitDoubleStand",
	DialogHitDoubleStand:      "HitDoubleStand",
}

var DialogPlayerMessage = map[Dialog]string{
	DialogNone:                "Invalid Dialog",
	DialogBetOrQuit:           "enter (B)et or (Q)uit [b]:",
	DialogPlaceYourBet:        "place your bet",
	DialogHitOrStand:          "please choose (H)it, (S)tand or (?)Hint: ",
	DialogHitSplitDoubleStand: "please choose (H)it, S(P)lit, (D)ouble, (S)tand or (?)Hint: ",
	DialogHitDoubleStand:      "please choose (H)it, (D)ouble, (S)tand (?)Hint: ",
}

func (d Dialog) String() string {
	return DialogMap[d]
}

type Game struct {
	Players              []*Player
	Dealer               *Player
	Shoe                 cards.Deck
	output               io.Writer
	input                io.Reader
	IsIncomingDeck       bool
	CardsDealt           int
	IncomingDeckPosition int
	DeckCount            int
	random               *rand.Rand
	CountCards           func(cards.Card, int, int, int) (int, float64)
	CardCounter          CardCounter
	Stage                Stage
	StageMessage         string
	NumberHumanPlayers   int
	NumberAiPlayers      int
	ActivePlayer         *Player
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

func WithNumberOfHumanPlayers(human int) Option {
	return func(g *Game) error {
		g.NumberHumanPlayers = human

		return nil
	}
}

func WithNumberOfAiPlayers(ai int) Option {
	return func(g *Game) error {
		g.NumberAiPlayers = ai
		return nil
	}
}

func NewBlackjackGame(opts ...Option) (*Game, error) {

	game := &Game{
		output:             os.Stdout,
		input:              os.Stdin,
		IsIncomingDeck:     true,
		DeckCount:          6,
		random:             rand.New(rand.NewSource(time.Now().UnixNano())),
		CountCards:         CountHiLo,
		NumberHumanPlayers: 1,
		NumberAiPlayers:    0,
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
		Name: "Dealer",
		Hands: []*Hand{
			{
				Id: 1,
			},
		},
	}

	return game, nil
}

func (g *Game) PlayAgain() bool {
	response := false
	if len(g.Players) > 0 {
		response = true
	}

	return response
}

func (g *Game) SetActivePlayer(player *Player) {
	g.ActivePlayer = player
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

	g.CardCounter.Count, g.CardCounter.TrueCount = g.CountCards(card, g.CardCounter.Count, g.CardsDealt, g.DeckCount)

	if g.IsIncomingDeck {
		if g.CardsDealt >= g.IncomingDeckPosition {
			g.StageMessage = "NEW DECK INCOMING"
			RenderStageMessage(g.output, g.StageMessage)

			g.Shoe = g.IncomingDeck()
			g.ResetFieldsAfterIncomingDeck()
		}

	} else {
		g.Shoe.Cards = append(g.Shoe.Cards, card)
	}
	return card
}

func (g *Game) ResetFieldsAfterIncomingDeck() {
	g.CardsDealt = 0
	g.CardCounter.Count = 0
	g.CardCounter.TrueCount = 0
}

func (g *Game) ResetPlayers() {

	for _, player := range g.Players {
		player.HandIndex = 0
		player.Hands = []*Hand{}
		hand := NewHand(1)
		player.AddHand(hand)
		player.Action = None
		player.Message = ""

	}
	g.Dealer.Hands = []*Hand{}
	hand := NewHand(1)
	g.Dealer.AddHand(hand)
	g.Dealer.Message = ""

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

func (g *Game) SetStage(stage Stage) {
	g.Stage = stage
	g.StageMessage = StageDisplayMessageMap[stage]
}

func (g *Game) AddPlayer(player *Player) {
	g.Players = append(g.Players, player)
}

type Player struct {
	Name           string
	Action         Action
	Bet            func(*Game) error
	Decide         func(io.Writer, io.Reader, *Player, cards.Card, int, CardCounter, Stage) Action
	AiRoundsToPlay int
	Record         Record
	Cash           int
	Hands          []*Hand
	HandIndex      int
	Message        string
	Dialog         Dialog
	CurrentBet     int
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

func (p *Player) Split(output io.Writer, card1, card2 cards.Card, index int) {

	id := p.NextHandId()
	hand := NewHand(id)
	p.AddHand(hand)
	indexNewHand := len(p.Hands) - 1
	//indexNewHand := index + 1

	// take last card in original hand and append to the new split hand
	card := p.Hands[index].Cards[1]
	p.Hands[indexNewHand].Cards = append(p.Hands[indexNewHand].Cards, card)

	// reset slice on original hand to only have the first card
	p.Hands[index].Cards = p.Hands[index].Cards[:len(p.Hands[index].Cards)-1]

	// mirror bet on new hand
	p.Hands[indexNewHand].Bet += p.Hands[index].Bet
	p.Cash -= p.Hands[index].Bet

	// add cards to each hand
	p.Hands[index].Cards = append(p.Hands[index].Cards, card1)
	p.Hands[indexNewHand].Cards = append(p.Hands[indexNewHand].Cards, card2)
	p.Hands[index].Action = None

	p.Hands[indexNewHand].Action = None

	p.Message = p.Hands[index].HandStringMulti(p.Name)
	RenderPlayerMessage(output, p)

	p.Message = p.Hands[indexNewHand].HandStringMulti(p.Name)
	RenderPlayerMessage(output, p)

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

func (p Player) NextHandId() int {
	return len(p.Hands) + 1
}

func (p *Player) SetDialog(dialog Dialog) error {

	switch dialog {
	case DialogBetOrQuit:
		str := []string{p.Name, " has $", strconv.Itoa(p.Cash), " ", DialogPlayerMessage[dialog], " "}
		p.Message = strings.Join(str, "")
	case DialogPlaceYourBet:
		str := []string{p.Name, " has $", strconv.Itoa(p.Cash), " ", DialogPlayerMessage[dialog], " ($1 to $", strconv.Itoa(p.Cash), " [$", strconv.Itoa(p.CurrentBet), "]): $"}
		p.Message = strings.Join(str, "")

	default:
		return fmt.Errorf("missing Dialog value switch, %v", dialog.String())
	}
	p.Dialog = dialog
	return nil
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

func (h Hand) ChooseAction() bool {
	return h.Action != ActionQuit && h.Action != ActionStand && h.Outcome != OutcomeBlackjack && h.Outcome != OutcomeBust
}

func (h Hand) HandStringMulti(name string) string {

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

func (h Hand) HandString(name string) string {

	builder := strings.Builder{}
	var response string

	if h.Action == ActionDoubleDown {
		builder.WriteString(name + " has ??: " + h.Cards[0].Render() + h.Cards[1].Render() + "[??]\n")
		response += builder.String()
	} else {
		for _, card := range h.Cards {
			builder.WriteString(card.Render())
		}
		str := []string{name, " has ", fmt.Sprint(h.Score()), ": ", builder.String(), "\n"}
		response = strings.Join(str, "")
	}

	return response
}

func (h Hand) DealerHandString() string {

	builder := strings.Builder{}
	var response string

	builder.WriteString("Dealer has ??: " + "[??]" + h.Cards[1].Render() + "\n")
	response += builder.String()

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

type CardCounter struct {
	Count     int
	TrueCount float64
}

func (c CardCounter) String() string {
	return "Count: " + strconv.Itoa(c.Count) + ", True Count: " + strconv.FormatFloat(c.TrueCount, 'f', -1, 64)
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

func HumanAction(output io.Writer, input io.Reader, player *Player, dealerCard cards.Card, index int, c CardCounter, stage Stage) Action {

	// check to see if not enough to split or double
	if player.Hands[index].Bet > player.Cash || len(player.Hands[index].Cards) > 2 {

		player.Dialog = DialogHitOrStand

	} else {
		// check if split is ok
		if player.Hands[index].Cards[0].Rank == player.Hands[index].Cards[1].Rank {
			player.Dialog = DialogHitSplitDoubleStand

		} else {
			// double ok
			player.Dialog = DialogHitDoubleStand

		}
	}
	str := []string{
		player.Name,
		" ",
		DialogPlayerMessage[player.Dialog],
	}
	player.Message = strings.Join(str, "")
	RenderPlayerMessage(output, player)
	RenderPlayerInput(output, input, player, stage, c, dealerCard)

	return player.Action
}

func HumanBet(g *Game) error {

	g.ActivePlayer.SetDialog(DialogBetOrQuit)

	// hack
	card := cards.Card{}

	RenderPlayerMessage(g.output, g.ActivePlayer)
	RenderPlayerInput(g.output, g.input, g.ActivePlayer, g.Stage, g.CardCounter, card)

	if g.ActivePlayer.Action != ActionQuit {
		g.ActivePlayer.SetDialog(DialogPlaceYourBet)
		RenderPlayerMessage(g.output, g.ActivePlayer)
		RenderPlayerInput(g.output, g.input, g.ActivePlayer, g.Stage, g.CardCounter, card)
	}

	return nil

}

// additional features
// difficult to fake concreate...take interface instead
// burn cards

// 17. ui
// 16. client/server
// 15. card counting ai
// 14. ai betting - inc or dec depending on last outcome
// 13. betting limits (max - min)
// 12. card counting - done
// 11. split - done
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
