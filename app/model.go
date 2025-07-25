package app

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

type Model struct {
	version     string
	term        string
	width       int
	height      int
	time        time.Time
	bg          string
	view        appview
	keys        KeyMap
	user        User
	help        help.Model
	mainStyle   lipgloss.Style
	infoStyle   lipgloss.Style
	actionStyle lipgloss.Style
	helpStyle   lipgloss.Style
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case time.Time:
		model.time = time.Time(msg)
	case CloseSplashMsg:
		model.view = status
		model.keys = DefaultKeyMap
	case tea.WindowSizeMsg:
		model.height = msg.Height
		model.width = msg.Width
		model.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, model.keys.Help):
			model.help.ShowAll = !model.help.ShowAll
		case key.Matches(msg, model.keys.Quit):
			return model, tea.Quit
		}
	}

	return model, nil
}

func (model Model) View() string {
	switch model.view {
	case splash:
		return model.splashView()
	case status:
		return model.layoutView(model.statusView())
	default:
		return "missing view"
	}
}

func (model Model) layoutView(inner string) string {
	help := model.helpView()
	padSize := model.height - strings.Count(inner, "\n") - strings.Count(help, "\n") - 1
	padding := strings.Repeat("\n", padSize)

	return inner + padding + help
}

func (model Model) helpView() string {
	help := model.help.View(model.keys)

	return model.helpStyle.Render(lipgloss.Place(model.width, 1, 0.5, 0.5, help))
}

func (model Model) splashView() string {
	title := model.actionStyle.Render("GOBOX")
	version := model.infoStyle.Render("v" + model.version)

	return lipgloss.Place(model.width, model.height, 0.5, 0.5, title+" "+version)
}

func (model Model) statusView() string {
	text := "Term: %s (%d x %d) [%s]\n"
	text += "Name: " + model.user.Name + "\n"
	text += "Role: " + model.user.Role.String() + "\n"
	text += "Time: " + model.time.Format(time.DateTime) + "\n"

	view := fmt.Sprintf(text, model.term, model.width, model.height, model.bg)

	return model.mainStyle.Render(view)
}
