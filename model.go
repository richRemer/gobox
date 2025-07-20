package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	term      string
	width     int
	height    int
	time      time.Time
	bg        string
	mainStyle lipgloss.Style
	infoStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case time.Time:
		m.time = time.Time(msg)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	text := "Your term is %s\n"
	text += "Your window size is x: %d, y: %d\n"
	text += "Background: %s\n"
	text += "Time: " + m.time.Format(time.RFC1123) + "\n"

	main := fmt.Sprintf(text, m.term, m.width, m.height, m.bg)
	quit := "Press 'q' to quit\n"
	info := lipgloss.Place(m.width, m.height-6, lipgloss.Center, lipgloss.Bottom, quit)

	return m.mainStyle.Render(main) + "\n\n" + m.infoStyle.Render(info)
}
