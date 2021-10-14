# Blackjack

Blackjack is a command line version of the classic game written in Go

Built with Aloha in Hawaii 🌊

Thank you to @bitfield for all of his mentoring on my Go journey!


# Installation Options
1) For Mac or Linux, execute the install.sh file
```bash
curl https://raw.githubusercontent.com/mbarley333/blackjack/main/install.sh | sh
```

2) For all OS
* Download the prebuilt binaries for your OS from the Releases section
* Unzip
* cd to folder
```bash
./blackjack
```

# Blackjack features
* All Blackjacks pay 2:1
* Split
* Double down
* Minimum bet $1
* Minimum 83% deck penetration before reshuffle
* Six deck shoe
* Card counting allowed and provided via HiLo method
* AI Players (Basic Strategy or Stand Only)
* Hints for Hit, Stand, Double and Split decisions


# Getting started
```bash
./blackjack -help

        Parameters:
          humanPlayers     Number of human players.  Default is 1
          aiPlayers        Number of Ai players.  Default is 0
          deckCount        Number of decks in shoe.  Default is 6

        Usage:
        ./blackjack -humanPlayers 1 -aiPlayers 1 -deckCount 7
```
* Enter number of players
* Enter name of player
* Select either AI or Human player
* For AI, enter number of hands to play
* From any command line, as a human player, enter "c" to get the card count and the true count
* In game hint available when "?" is displayed from command line
```bash
	Player1 please choose (H)it, (D)ouble, (S)tand (?)Hint: 
```


# Adding a custom AI
* Clone repo to desktop and cd
* Add new AI player algo to blackjack/ai.go
* Open blackjack/blackjack.go
* Update the x value for the PlayerTypeInput map
* Update the PlayerTypeAiCustom value for the PlayerTypeBetMap
* Build

```bash

var PlayerTypeInputMap = map[string]PlayerType{
	"h": PlayerTypeHuman,
	"b": PlayerTypeAiBasic,
	"s": PlayerTypeAiStandOnly,
	"x": PlayerTypeAiBasic,
}

var PlayerTypeBetMap = map[PlayerType]func(io.Writer, io.Reader, *Player, int, CardCounter) error{
	PlayerTypeHuman:       HumanBet,
	PlayerTypeAiStandOnly: AiBet,
	PlayerTypeAiBasic:     AiBet,
	PlayerTypeAiCustom:    AiBet,
  
```








