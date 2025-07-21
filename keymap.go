package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding
}

var activeKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
