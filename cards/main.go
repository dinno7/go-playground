package main

import (
	"card/deck"
)

func main() {
	cards := deck.NewDeck()
	cards.Print()

	cards.Shuffle()

	cards.Print()
}
