package main

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	Help key.Binding
	Quit key.Binding
}

func (keys keymap) ShortHelp() []key.Binding {
	return []key.Binding{keys.Help, keys.Quit}
}

func (keys keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keys.Help, keys.Quit},
	}
}

var activeKeyMap = keymap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
