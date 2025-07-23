package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Help key.Binding
	Load key.Binding
	New  key.Binding
	Quit key.Binding
}

func (keys KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.New, keys.Load, keys.Quit}
}

func (keys KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keys.New, keys.Load, keys.Quit},
		{keys.Help},
	}
}

var ActiveKeyMap = KeyMap{
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
