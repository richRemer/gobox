package app

import (
	"fmt"
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
	helpHeight  int
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

	fullHelp := model.help
	fullHelp.ShowAll = true
	model.helpHeight = lipgloss.Height(fullHelp.View(model.keys)) + 1
}

func (model Model) layoutView(inner string) string {
	help := model.helpView()
	height := model.height - lipgloss.Height(help)
	main := lipgloss.Place(model.width, height, 0.5, 0.5, inner)

	return lipgloss.JoinVertical(lipgloss.Center, main, help)
}

func (model Model) helpView() string {
	help := model.help.View(model.keys)
	view := model.helpStyle.Render(lipgloss.PlaceHorizontal(model.width, 0.5, help))

	return lipgloss.PlaceVertical(model.helpHeight, 1.0, view)
}

func (model Model) splashView() string {
	title := model.actionStyle.Render("GOBOX")
	version := model.infoStyle.Render("v" + model.version)

	return lipgloss.Place(model.width, model.height, 0.5, 0.5, title+" "+version)
}

func (model Model) guestView() string {
	text := "Welcome to my house!\n"
	text += "Enter freely.\n"
	text += "Go safely, and leave something of the happiness you bring."

	return lipgloss.Place(model.width, 3, 0.5, 0.5, text)
}

func (model Model) statusView() string {
	text := "Term: %s (%d x %d) [%s]\n"
	text += "Name: " + model.user.Name + "\n"
	text += "Role: " + model.user.Role.String() + "\n"
	text += "Time: " + model.time.Format(time.DateTime)

	view := fmt.Sprintf(text, model.term, model.width, model.height, model.bg)

	return model.mainStyle.Render(view)
}
