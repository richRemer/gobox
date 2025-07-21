package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type terminal struct {
	term      string
	width     int
	height    int
	time      time.Time
	bg        string
	mainStyle lipgloss.Style
	infoStyle lipgloss.Style
}

func (term terminal) Init() tea.Cmd {
	return nil
}

func (term terminal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case time.Time:
		term.time = time.Time(msg)
	case tea.WindowSizeMsg:
		term.height = msg.Height
		term.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, activeKeyMap.Quit):
			return term, tea.Quit
		}
	}

	return term, nil
}

func (term terminal) View() string {
	text := "Your term is %s\n"
	text += "Your window size is x: %d, y: %d\n"
	text += "Background: %s\n"
	text += "Time: " + term.time.Format(time.RFC1123) + "\n"

	main := fmt.Sprintf(text, term.term, term.width, term.height, term.bg)
	quit := "Press 'q' to quit\n"
	info := lipgloss.Place(term.width, term.height-6, lipgloss.Center, lipgloss.Bottom, quit)

	return term.mainStyle.Render(main) + "\n\n" + term.infoStyle.Render(info)
}
