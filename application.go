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

type appview int

const (
	splash appview = 1
	status appview = 2
)

type application struct {
	version     string
	term        string
	width       int
	height      int
	time        time.Time
	bg          string
	view        appview
	keys        keymap
	help        help.Model
	splashTime  int
	mainStyle   lipgloss.Style
	infoStyle   lipgloss.Style
	actionStyle lipgloss.Style
	helpStyle   lipgloss.Style
}

func (app application) Init() tea.Cmd {
	return nil
}

func (app application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case time.Time:
		app.time = time.Time(msg)
		app.splashTime--

		if app.splashTime == 0 {
			app.view = status
		}
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
	switch app.view {
	case splash:
		return app.splashView()
	case status:
		return app.layoutView(app.statusView())
	default:
		return "missing view"
	}
}

func (app application) layoutView(inner string) string {
	help := app.helpView()
	padSize := app.height - strings.Count(inner, "\n") - strings.Count(help, "\n") - 1
	padding := strings.Repeat("\n", padSize)

	return inner + padding + help
}

func (app application) helpView() string {
	help := app.help.View(app.keys)

	return app.helpStyle.Render(lipgloss.Place(app.width, 1, 0.5, 0.5, help))
}

func (app application) splashView() string {
	title := app.actionStyle.Render("GOBOX")
	version := app.infoStyle.Render("v" + app.version)

	return lipgloss.Place(app.width, app.height, 0.5, 0.5, title+" "+version)
}

func (app application) statusView() string {
	text := "Your term is %s\n"
	text += "Your window size is x: %d, y: %d\n"
	text += "Background: %s\n"
	text += "Time: " + app.time.Format(time.DateTime) + "\n"

	view := fmt.Sprintf(text, app.term, app.width, app.height, app.bg)

	return app.mainStyle.Render(view)
}
