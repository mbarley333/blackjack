package blackjack

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mbarley333/cards"
)

func RunCLI() {

	flag.Usage = help

	humanPlayersPtr := flag.Int("humanPlayers", 1, "Number of human players.  Default is 1")
	aiPlayersPtr := flag.Int("aiPlayers", 0, "Number of AI players.  Default is 0")
	deckCountPtr := flag.Int("deckCount", 6, "Number of decks in shoe.  Default is 6")

	flag.Parse()

	g, err := NewBlackjackGameWithArgs(*humanPlayersPtr, *aiPlayersPtr, *deckCountPtr)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot create new blackjack game, %s", err))
	}

	g.AddBlackjackPlayers()

	for g.PlayAgain() {

		g.ResetPlayers()
		g.Betting()
		g.Players = g.RemoveQuitPlayers()
		if g.PlayAgain() {

			g.OpeningDeal()
			g.Deciding()
			g.DealerPlay()
			g.Outcome(g.output)
		}
	}
	fmt.Fprintln(g.output, "No players left in game.  Exiting...")
}

func NewBlackjackGameWithArgs(humanPlayers, aiPlayers, deckCount int) (*Game, error) {

	g, err := NewBlackjackGame(
		WithNumberOfHumanPlayers(humanPlayers),
		WithNumberOfAiPlayers(aiPlayers),
		WithDeckCount(deckCount),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create game, %s", err)
	}

	return g, nil

}

func (g *Game) AddBlackjackPlayers() {

	for i := 0; i < g.NumberHumanPlayers; i++ {
		player := NewHumanPlayer(g.output, g.input, i)
		g.AddPlayer(player)
	}

	for i := 0; i < g.NumberAiPlayers; i++ {
		player := NewAiPlayer(g.output, g.input, i)
		g.AddPlayer(player)
	}

}

func (g *Game) Betting() error {

	g.SetStage(StageBetting)

	var err error

	for _, player := range g.Players {
		g.ActivePlayer = player

		err = player.Bet(g)

		if err != nil {
			return fmt.Errorf("unable to place bet for player: %s", player.Name)
		}
	}
	return nil
}

func (g *Game) OpeningDeal() {

	g.SetStage(StageOpeningDeal)
	RenderStageMessage(g.output, g.StageMessage)

	for i := 0; i < 2; i++ {
		for _, player := range g.Players {
			g.SetActivePlayer(player)
			card := g.Deal(g.output)
			player.Hands[0].Cards = append(player.Hands[0].Cards, card)

			player.Message = player.Name + " is dealt the [" + card.String() + "]\n"

			RenderPlayerMessage(g.output, player)
			time.Sleep(750 * time.Millisecond)

		}
		card := g.Deal(g.output)
		g.Dealer.Hands[0].Cards = append(g.Dealer.Hands[0].Cards, card)

		if i == 0 {
			g.Dealer.Message = "Dealer is dealt a [???]\n" + "\n"
		} else {
			g.Dealer.Message = "Dealer is dealt a [" + card.String() + "]\n"
		}
		RenderPlayerMessage(g.output, g.Dealer)
		time.Sleep(750 * time.Millisecond)

	}

	RenderPlayerAndDealerCards(g.output, g.input, g.Players, g.Dealer, g.Stage)
}

func NewHumanPlayer(output io.Writer, input io.Reader, index int) *Player {

	defaultName := "Player" + strconv.Itoa(index+1)

	fmt.Fprintf(output, "%s enter your name [%s]: ", defaultName, defaultName)

	reader := bufio.NewReader(input)
	name, _ := reader.ReadString('\n')
	name = strings.Replace(name, "\n", "", -1)

	if len(name) == 0 {
		name = defaultName
	}

	player := &Player{

		Name:       name,
		Decide:     HumanAction,
		Bet:        HumanBet,
		CurrentBet: 1,
		Cash:       100,
		Hands: []*Hand{
			{Id: 1},
		},
	}

	return player

}

func NewAiPlayer(output io.Writer, input io.Reader, index int) *Player {
	var name string
	var playerTypeInput string
	var aiRoundsToPlay string

	defaultName := "AiPlayer" + strconv.Itoa(index+1)

	reader := bufio.NewReader(input)

	fmt.Fprintf(output, "%s enter your name [%s]: ", defaultName, defaultName)
	name, _ = reader.ReadString('\n')
	name = strings.Replace(name, "\n", "", -1)

	if name == "" {
		name = defaultName
	}

	for strings.ToLower(playerTypeInput) != "b" && strings.ToLower(playerTypeInput) != "s" && strings.ToLower(playerTypeInput) != "x" {
		fmt.Fprintf(output, "Select AI Type (B)asic Strategy, (S)tandOnly or (X)custom [B]: ")
		//fmt.Fscanln(input, &playerTypeInput)
		playerTypeInput, _ = reader.ReadString('\n')
		playerTypeInput = strings.Replace(playerTypeInput, "\n", "", -1)
		if playerTypeInput == "" {
			playerTypeInput = "b"
		}

	}

	playerTypeInputValue := PlayerTypeInputMap[strings.ToLower(playerTypeInput)]
	playerType := PlayerTypeMap[playerTypeInputValue]
	playerTypeBet := PlayerTypeBetMap[playerTypeInputValue]

	aiHands := 0
	var err error

	fmt.Fprint(output, "Enter number of rounds the AI plays: ")
	//fmt.Fscan(input, &aiRoundsToPlay)
	aiRoundsToPlay, _ = reader.ReadString('\n')
	aiRoundsToPlay = strings.Replace(aiRoundsToPlay, "\n", "", -1)
	aiHands, err = strconv.Atoi(aiRoundsToPlay)
	if err != nil {
		return nil
	}

	player := &Player{

		Name:           name,
		Decide:         playerType,
		Bet:            playerTypeBet,
		AiRoundsToPlay: aiHands,
		Cash:           100,
		Hands: []*Hand{
			{Id: 1},
		},
	}

	return player

}

func (g *Game) Deciding() error {

	g.SetStage(StageDeciding)

	var err error

	for _, player := range g.Players {

		g.StageMessage = strings.ToUpper(player.Name) + " MAKE YOUR CHOICE"
		RenderStageMessage(g.output, g.StageMessage)

		g.ActivePlayer = player

		err = g.PlayHand(player)
		if err != nil {
			return err
		}

		if err != nil {
			return fmt.Errorf("unable to place bet for player: %s", player.Name)
		}
	}
	return nil
}

func (g *Game) PlayHand(player *Player) error {
	for index, hand := range player.Hands {
		if hand.Score() == 21 {
			hand.Outcome = OutcomeBlackjack
		}

		err := RenderPlayerAndDealerCards(g.output, g.input, g.Players, g.Dealer, g.Stage)
		if err != nil {
			return err
		}

		for hand.ChooseAction() {
			if hand.Action == None {
				player.Message = hand.HandString(player.Name)
				RenderPlayerMessage(g.output, player)

				hand.Action = player.Decide(g.output, g.input, player, g.Dealer.Hands[0].Cards[0], index, g.CardCounter, g.Stage)
			}
			if hand.Action == ActionHit {
				card := g.Deal(g.output)
				player.Message = player.Name + " is dealt the [" + card.String() + "]\n\n"
				hand.Hit(g.output, card, player.Name)
				RenderPlayerMessage(g.output, player)
			} else if hand.Action == ActionDoubleDown {
				player.Cash -= hand.Bet
				card := g.Deal(g.output)
				player.Message = player.Name + " is dealt [???]\n\n"
				hand.DoubleDown(g.output, card, player.Name)
				RenderPlayerMessage(g.output, player)
			} else if hand.Action == ActionSplit {
				card1 := g.Deal(g.output)
				card2 := g.Deal(g.output)
				player.Split(g.output, card1, card2, index)
				err = g.PlayHand(player)
				if err != nil {
					return err
				}

			}
		}
	}
	return nil

}

func (g *Game) DealerPlay() {

	g.SetStage(StageDealerPlay)
	RenderStageMessage(g.output, g.StageMessage)
	dealerOk := g.IsDealerDraw()

	if dealerOk {

		for g.Dealer.Hands[0].Score() <= 16 || (g.Dealer.Hands[0].Score() == 17 && g.Dealer.Hands[0].MinScore() != 17) {
			card := g.Deal(g.output)
			g.Dealer.Hands[0].Cards = append(g.Dealer.Hands[0].Cards, card)
			g.Dealer.Message = "Dealer is dealt a [" + card.String() + "]\n"
			RenderPlayerMessage(g.output, g.Dealer)
			time.Sleep(2 * time.Second)
		}
		g.Dealer.Hands[0].Action = ActionStand
		RenderPlayerAndDealerCards(g.output, g.input, g.Players, g.Dealer, g.Stage)
	}
}

func (g *Game) Outcome(output io.Writer) {

	g.SetStage(StageOutcome)
	RenderStageMessage(g.output, g.StageMessage)
	RenderPlayerAndDealerCards(g.output, g.input, g.Players, g.Dealer, g.Stage)

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

func RenderGameCli(output io.Writer, input io.Reader, g *Game) error {

	RenderStageMessage(output, g.StageMessage)

	err := RenderPlayerAndDealerCards(output, input, g.Players, g.Dealer, g.Stage)
	if err != nil {
		return err
	}

	return nil
}

func RenderPlayerAndDealerCards(output io.Writer, input io.Reader, players []*Player, dealer *Player, stage Stage) error {

	fmt.Fprint(output, "\n\n")
	if stage != StageStart {
		if len(players) == 0 {
			return fmt.Errorf("game struct does not have any players set")
		}

		for _, player := range players {
			for _, hand := range player.Hands {
				if player.HandIndex == 0 {
					fmt.Fprint(output, hand.HandString(player.Name))
				} else {
					fmt.Fprint(output, hand.HandStringMulti(player.Name))
				}

			}
		}

		if len(dealer.Hands[0].Cards) == 0 {
			return fmt.Errorf("dealer does not have any cards.  requires a cards[0] for string")
		}
		if stage != StageOutcome {
			for _, h := range dealer.Hands {
				fmt.Fprint(output, h.DealerHandString())
			}
		} else {
			for _, h := range dealer.Hands {
				fmt.Fprint(output, h.HandString(dealer.Name))
			}
		}

	}
	fmt.Fprint(output, "\n\n")
	return nil
}

func RenderPlayerMessage(output io.Writer, player *Player) {

	fmt.Fprint(output, player.Message)

}

func RenderPlayerInput(output io.Writer, input io.Reader, player *Player, stage Stage, c CardCounter, dealerCard cards.Card) error {

	ok := false
	var err error

	for !ok {

		reader := bufio.NewReader(input)
		answer, _ := reader.ReadString('\n')
		answer = strings.Replace(answer, "\n", "", -1)

		// show count
		if answer == "c" {
			fmt.Fprintln(output, c.String())
		}

		if answer == "?" && stage == StageDeciding {
			index := 0
			action := GetHint(output, input, player, dealerCard, index, c, stage)
			hint := "The suggested action is to " + action.String() + "\n"
			fmt.Fprintln(output, hint)
		}

		ok, err = IsInputValid(answer, player)
		if err != nil {
			return err
		}

		if ok {
			err = player.SetAction(answer)
			if err != nil {
				return err
			}
		} else {
			RenderPlayerMessage(output, player)
		}

	}

	return nil
}

func RenderStageMessage(output io.Writer, message string) {
	messageLength := len(message)

	builder := strings.Builder{}

	for i := 0; i < messageLength+4; i++ {
		builder.WriteString("*")
	}

	// game stage message
	fmt.Fprint(output, "\n\n\n")
	fmt.Fprintln(output, builder.String())
	fmt.Fprintln(output, "* "+message+" *")
	fmt.Fprintln(output, builder.String())
	fmt.Fprint(output, "\n\n\n")

}

func (p *Player) SetAction(answer string) error {

	var err error
	switch p.Dialog {
	case DialogBetOrQuit:
		if answer == "" {
			p.Action = ActionBet
		} else {
			p.Action = ActionMap[strings.ToLower(answer)]
		}
	case DialogPlaceYourBet:
		var bet int
		if answer == "" {
			bet = p.CurrentBet
		} else {
			bet, err = strconv.Atoi(answer)
		}
		if err != nil {
			return fmt.Errorf("unable to set bet amount, %s", err)
		}

		p.CurrentBet = bet
		p.Cash -= bet
		p.Hands[p.HandIndex].Bet += bet
	case DialogHitOrStand:
		p.Action = ActionMap[strings.ToLower(answer)]
	case DialogHitDoubleStand:
		p.Action = ActionMap[strings.ToLower(answer)]
	case DialogHitSplitDoubleStand:
		p.Action = ActionMap[strings.ToLower(answer)]

	default:
		return fmt.Errorf("missing Dialog in switch, %s", p.Dialog.String())
	}
	return nil

}

func IsInputValid(answer string, player *Player) (bool, error) {

	ok := false

	switch player.Dialog {
	case DialogBetOrQuit:
		if strings.ToLower(answer) == "b" || strings.ToLower(answer) == "q" || answer == "" {
			ok = true
		}
	case DialogPlaceYourBet:

		bet, err := strconv.Atoi(answer)
		if answer == "" {
			ok = true
		} else if err != nil {
			ok = false
		} else if bet < 1 || bet > player.Cash {
			ok = false
		} else {
			ok = true
		}
	case DialogHitOrStand:
		if strings.ToLower(answer) == "h" || strings.ToLower(answer) == "s" {
			ok = true
		}
	case DialogHitDoubleStand:
		if strings.ToLower(answer) == "h" || strings.ToLower(answer) == "s" || strings.ToLower(answer) == "d" {
			ok = true
		}
	case DialogHitSplitDoubleStand:
		if strings.ToLower(answer) == "h" || strings.ToLower(answer) == "s" || strings.ToLower(answer) == "d" || strings.ToLower(answer) == "p" {
			ok = true
		}
	default:
		return ok, fmt.Errorf("missing Dialog in switch, %s", player.Dialog.String())
	}

	return ok, nil
}

// func GetCliInput(input io.Reader) string {

// 	reader := bufio.NewReader(input)
// 	answer, _ := reader.ReadString('\n')
// 	answer = strings.Replace(answer, "\n", "", -1)

// 	return answer
// }

func help() {
	fmt.Fprintln(os.Stderr, `
	Parameters:
	  humanPlayers     Number of human players.  Default is 1
	  aiPlayers        Number of Ai players.  Default is 0
	  deckCount        Number of decks in shoe.  Default is 6
	
	Usage:
	./blackjack -humanPlayers 1 -aiPlayers 1 -deckCount 7
	`)
}
