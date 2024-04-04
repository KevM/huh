package main

import (
	"math/rand"
	"strings"
)

type TickerFinder struct {
	tickers []Ticker
}

func NewTickerFinder() *TickerFinder {
	tickers := generateTickers(1000)
	return &TickerFinder{tickers: tickers}
}

func (f *TickerFinder) Find(fragment string) ([]Ticker, error) {
	var tickers []Ticker
	for _, t := range f.tickers {
		if strings.Contains(strings.ToLower(t.Symbol), fragment) || strings.Contains(strings.ToLower(t.Name), fragment) {
			tickers = append(tickers, t)
		}
	}

	return tickers, nil
}

func generateTickers(count int) []Ticker {
	tickers := make([]Ticker, count)
	for i := range tickers {
		tickers[i] = Ticker{
			Symbol: randomString(4),
			Name:   randomCompanyName(4, 9),
		}
	}
	return tickers
}

func randomString(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randomCompanyName(min, max int) string {
	var nameAtoms = []string{
		"Meter",
		"Science",
		"Tech",
		"Data",
		"AI",
		"Robotics",
		"Automation",
		"Systems",
		"Analytics",
		"Cloud",
		"Charm",
		"Bracelet",
		"Software",
		"Hardware",
		"Basil",
		"Ganglia",
		"Kale",
		"Kelp",
		"Corp",
		"Inc",
		"Co",
		"LLC",
		"Group",
	}

	nameLength := rand.Intn(max-min+1) + min
	companyName := make([]string, nameLength)
	for i := range companyName {
		companyName[i] = nameAtoms[rand.Intn(len(nameAtoms))]
	}
	return strings.Join(companyName, " ")
}
