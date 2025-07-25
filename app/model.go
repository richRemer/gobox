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

type ViewMode int

const (
	SplashView ViewMode = 1
	StatusView ViewMode = 2
	GuestView  ViewMode = 3
)

type Model struct {
	version     string
	term        string
	width       int
	height      int
	time        time.Time
	bg          string
	view        ViewMode
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
		model.SelectDefaultView()
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
	case SplashView:
		return model.splashView()
	case StatusView:
		return model.layoutView(model.statusView())
	case GuestView:
		return model.layoutView(model.guestView())
	default:
		return "missing view"
	}
}

func (model *Model) SelectDefaultView() {
	if model.user.Role == RoleGuest {
		model.SelectView(GuestView)
	} else {
		model.SelectView(StatusView)
	}
}

func (model *Model) SelectView(view ViewMode) {
	model.view = view

	switch model.view {
	case GuestView:
		model.keys = GuestKeyMap
	case SplashView:
		model.keys = SplashKeyMap
	default:
		model.keys = DefaultKeyMap
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

func (model Model) guestView() string {
	text := "Welcome to my house!\n"
	text += "Enter freely.\n"
	text += "Go safely, and leave something of the happiness you bring.\n"

	return lipgloss.Place(model.width, 3, 0.5, 0.5, text)
}

func (model Model) statusView() string {
	text := "Term: %s (%d x %d) [%s]\n"
	text += "Name: " + model.user.Name + "\n"
	text += "Role: " + model.user.Role.String() + "\n"
	text += "Time: " + model.time.Format(time.DateTime) + "\n"

	view := fmt.Sprintf(text, model.term, model.width, model.height, model.bg)

	return model.mainStyle.Render(view)
}
