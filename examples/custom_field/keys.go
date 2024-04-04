package main

import "github.com/charmbracelet/bubbles/key"

type PickerKeyMap struct {
	AcceptSuggestion key.Binding
	NextSuggestion   key.Binding
	PrevSuggestion   key.Binding
	Submit           key.Binding
}

var DefaultPickerKeyMap = PickerKeyMap{
	AcceptSuggestion: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "pick suggestion")),
	NextSuggestion:   key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("up", "prev")),
	PrevSuggestion:   key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("down", "next")),
	Submit:           key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
}

func (k PickerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.AcceptSuggestion,
		k.NextSuggestion,
		k.PrevSuggestion,
		k.Submit,
	}
}

func (k PickerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
