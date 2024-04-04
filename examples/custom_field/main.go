package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func main() {

	var ticker string
	var picked bool

	finder := NewTickerFinder() // this would normally be a web service call...

	huh.NewForm(
		huh.NewGroup(NewTickerPicker(finder.Find).
			Title("Pick a ticker").
			Value(&ticker)),

		huh.NewGroup(huh.NewConfirm().
			Title("Did you pick a ticker?").
			Value(&picked)),
	).Run()

	if picked {
		fmt.Printf("Ticker %s was picked.\n", ticker)
	}
}
