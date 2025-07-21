package main

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	Help key.Binding
	Load key.Binding
	New  key.Binding
	Quit key.Binding
}

func (keys keymap) ShortHelp() []key.Binding {
	return []key.Binding{keys.New, keys.Load, keys.Quit}
}

func (keys keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keys.New, keys.Load, keys.Quit},
		{keys.Help},
	}
}

var activeKeyMap = keymap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Load: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "load"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
