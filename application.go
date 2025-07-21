package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
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
	help      help.Model
	mainStyle lipgloss.Style
	infoStyle lipgloss.Style
	helpStyle lipgloss.Style
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
		app.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, app.keys.Help):
			app.help.ShowAll = !app.help.ShowAll
		case key.Matches(msg, app.keys.Quit):
			return app, tea.Quit
		}
	}

	return app, nil
}

func (app application) View() string {
	main := app.statusView()
	help := app.help.View(app.keys)

	mainView := app.mainStyle.Render(main)
	helpView := app.helpStyle.Render(lipgloss.Place(app.width, 1, 0.5, 0.5, help))
	space := app.height - strings.Count(mainView, "\n") - strings.Count(helpView, "\n") - 1
	spacing := strings.Repeat("\n", space)

	return mainView + spacing + helpView
}

func (app application) statusView() string {
	text := "Your term is %s\n"
	text += "Your window size is x: %d, y: %d\n"
	text += "Background: %s\n"
	text += "Time: " + app.time.Format(time.DateTime) + "\n"

	return fmt.Sprintf(text, app.term, app.width, app.height, app.bg)
}
