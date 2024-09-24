package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func main() {
	// variable for bank balance
	var balance int
	var balanceMx sync.Mutex

	// print out starting values
	fmt.Printf("Initial account balance: $%d.00", balance)
	fmt.Println()

	// define weekly revenue
	incomes := []Income{
		{Source: "Main job", Amount: 500},
		{Source: "Gifts", Amount: 10},
		{Source: "Part time job", Amount: 50},
		{Source: "Investments", Amount: 100},
	}

	wg.Add(len(incomes))

	// loop through 52 weeks and print out how much is made; keep a running total
	for i, income := range incomes {
		go func(i int, income Income) {
			defer wg.Done()

			for week := 1; week <= 52; week++ {
				balanceMx.Lock()
				temp := balance
				temp += income.Amount
				balance = temp
				balanceMx.Unlock()

				fmt.Printf(
					"On week %d, you earned $%d.00 from %s\n",
					week,
					income.Amount,
					income.Source,
				)
			}
		}(i, income)
	}

	wg.Wait()

	// print out final balance
	fmt.Printf("Final bank balance: $%d.00", balance)
	fmt.Println()
}

