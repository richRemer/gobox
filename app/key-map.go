package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Cancel   key.Binding
	Confirm  key.Binding
	Help     key.Binding
	Quit     key.Binding
	Register key.Binding
}

func (keys KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.Register, keys.Quit}
}

func (keys KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keys.Register, keys.Quit},
		{keys.Confirm, keys.Cancel, keys.Help},
	}
}

var SplashKeyMap = KeyMap{
	Quit: DefaultKeyMap.Quit,
}

var GuestKeyMap = KeyMap{
	Help: DefaultKeyMap.Help,
	Quit: DefaultKeyMap.Quit,
	Register: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "register"),
	),
}

var RegisterKeyMap = KeyMap{
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "done"),
	),
}

var DefaultKeyMap = KeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
