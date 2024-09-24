package deck

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type deck []string

func NewDeck() deck {
	card := deck{}
	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four"}

	for _, suite := range cardSuits {
		for _, value := range cardValues {
			card = append(card, fmt.Sprintf("%s of %s", suite, value))
		}
	}

	return card
}

func NewDeckFromString(stringDeck string) deck {
	fileDeck := strings.Split(string(stringDeck), ",")
	return deck(fileDeck)
}

func NewDeckFromFile(fileName string) deck {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(errors.New("Something went wrong when writing to file"))
	}
	return NewDeckFromString(string(data))
}

func (d deck) Print() {
	fmt.Println("---- DECK VALUES ----")
	for _, card := range d {
		fmt.Println(card)
	}
}

func (d deck) Deal(handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func (d deck) String() string {
	return strings.Join(d, ",")
}

func (d deck) SaveToFile(fileName string) {
	catdsStr := d.String()
	err := os.WriteFile(fileName, []byte(catdsStr), 0666)
	if err != nil {
		panic(errors.New("Something went wrong when writing to file"))
	}
}

func (d deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	r.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}
