package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type ViewMode int

const (
	SplashView   ViewMode = 1
	StatusView   ViewMode = 2
	GuestView    ViewMode = 3
	RegisterView ViewMode = 4
)

type Model struct {
	program     *tea.Program
	version     string
	term        string
	width       int
	height      int
	time        time.Time
	bg          string
	view        ViewMode
	keys        KeyMap
	user        User
	publicKey   string
	help        help.Model
	input       textinput.Model
	errors      ErrorModel
	users       UserRepo
	helpHeight  int
	mainStyle   lipgloss.Style
	infoStyle   lipgloss.Style
	actionStyle lipgloss.Style
	helpStyle   lipgloss.Style
	inputStyle  lipgloss.Style
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd = nil

	switch msg := msg.(type) {
	case time.Time:
		model.time = time.Time(msg)
	case CloseSplashMsg:
		model = model.WithDefaultView()
	case ProgramMsg:
		model.program = msg.program
	case tea.WindowSizeMsg:
		model.height = msg.Height
		model.width = msg.Width
		model.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, model.keys.Cancel):
			model = model.WithDefaultView()
			return model, nil
		case key.Matches(msg, model.keys.Confirm):
			nick := model.input.Value()
			user, err := model.users.RegisterWithKey(nick, model.publicKey)

			if err == nil {
				log.Info("registration", "nick", nick)
				model.user = user
				model = model.WithDefaultView()
				return model, nil
			} else {
				log.Error("registration failed", "nick", nick, "err", err)
				return model, func() tea.Msg { return ErrorMsg{err: err} }
			}
		case key.Matches(msg, model.keys.Help):
			model.help.ShowAll = !model.help.ShowAll
			return model, nil
		case key.Matches(msg, model.keys.Quit):
			return model, tea.Quit
		case key.Matches(msg, model.keys.Register):
			model = model.WithView(RegisterView)
			return model, nil
		}
	}

	model.errors, cmd = model.errors.Update(msg)
	cmds = append(cmds, cmd)

	model.input, cmd = model.input.Update(msg)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model Model) View() string {
	switch model.view {
	case SplashView:
		return model.splashView()
	case StatusView:
		return model.layoutView(model.statusView())
	case GuestView:
		return model.layoutView(model.guestView())
	case RegisterView:
		return model.layoutView(model.registerView())
	default:
		return "missing view"
	}
}

func (model Model) WithDefaultView() Model {
	if model.user.Role == RoleGuest {
		return model.WithView(GuestView)
	} else {
		return model.WithView(StatusView)
	}
}

func (model Model) WithView(view ViewMode) Model {
	model.view = view

	switch model.view {
	case GuestView:
		model.keys = GuestKeyMap
	case SplashView:
		model.keys = SplashKeyMap
	case RegisterView:
		model.keys = RegisterKeyMap
		model.input = textinput.New()
		model.input.Placeholder = "Nickname"
		model.input.CharLimit = 32
		model.input.Width = 24
		model.input.Cursor.Blink = true
		model.input.Focus()
	default:
		model.keys = DefaultKeyMap
	}

	fullHelp := model.help
	fullHelp.ShowAll = true

	model.helpHeight = lipgloss.Height(fullHelp.View(model.keys)) + 1

	return model
}

func (model Model) layoutView(inner string) string {
	help := model.helpView()
	errors := model.errorView()
	height := model.height - lipgloss.Height(help) - lipgloss.Height(errors)
	main := lipgloss.Place(model.width, height, 0.5, 0.5, inner)

	return lipgloss.JoinVertical(lipgloss.Center, errors, main, help)
}

func (model Model) helpView() string {
	help := model.help.View(model.keys)
	view := model.helpStyle.Render(lipgloss.PlaceHorizontal(model.width, 0.5, help))

	return lipgloss.PlaceVertical(model.helpHeight, 1.0, view)
}

func (model Model) errorView() string {
	errors := model.errors.View()
	view := lipgloss.Place(model.width, 1, 0.5, 0.5, errors)

	return view
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

func (model Model) registerView() string {
	return model.inputStyle.Render(model.input.View())
}

func (model Model) statusView() string {
	text := "Term: %s (%d x %d) [%s]\n"
	text += "Name: " + model.user.Name + "\n"
	text += "Role: " + model.user.Role.String() + "\n"
	text += "Time: " + model.time.Format(time.DateTime)

	view := fmt.Sprintf(text, model.term, model.width, model.height, model.bg)

	return model.mainStyle.Render(view)
}
