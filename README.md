# Blackjack

Blackjack is a command line version of the classic game written in Go

Built with Aloha in Hawaii 🌊

Thank you to @bitfield for his mentoring on my Go journey!


# Installation

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
* Card counting allowed and provided via HiLo method
* AI Players (Basic Strategy or Stand Only)


# Getting started
```bash
./blackjack
```
* Enter number of players
* Enter name for player
* Select either AI or Human player
* For AI, enter number of hands to play
* From any command line, as a human player, enter "c" to get the count and true count


# Adding custom AI
* Clone repo to desktop and cd
* Add new AI player algo to blackjack/ai.go
* Open blackjack/blackjack.go








