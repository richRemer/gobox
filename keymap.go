package main

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	Quit key.Binding
}

var activeKeyMap = keymap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
