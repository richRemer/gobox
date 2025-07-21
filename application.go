package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type application struct {
	term      string
	width     int
	height    int
	time      time.Time
	bg        string
	keys      keymap
	mainStyle lipgloss.Style
	infoStyle lipgloss.Style
}

func (app application) Init() tea.Cmd {
	return nil
}

func (app application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case time.Time:
		app.time = time.Time(msg)
	case tea.WindowSizeMsg:
		app.height = msg.Height
		app.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, app.keys.Quit):
			return app, tea.Quit
		}
	}

	return app, nil
}

func (app application) View() string {
	text := "Your term is %s\n"
	text += "Your window size is x: %d, y: %d\n"
	text += "Background: %s\n"
	text += "Time: " + app.time.Format(time.RFC1123) + "\n"

	main := fmt.Sprintf(text, app.term, app.width, app.height, app.bg)
	quit := "Press 'q' to quit\n"
	info := lipgloss.Place(app.width, app.height-6, lipgloss.Center, lipgloss.Bottom, quit)

	return app.mainStyle.Render(main) + "\n\n" + app.infoStyle.Render(info)
}
