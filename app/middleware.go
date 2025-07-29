package app

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"
)

func Middleware(users UserRepo) wish.Middleware {
	program := func(model tea.Model, options ...tea.ProgramOption) *tea.Program {
		program := tea.NewProgram(model, options...)

		go func() {
			for {
				<-time.After(1 * time.Second)
				program.Send(time.Time(time.Now()))
			}
		}()

		go func() {
			<-time.After(2 * time.Second)
			program.Send(CloseSplashMsg{})
		}()

		return program
	}

	handler := func(session ssh.Session) *tea.Program {
		pty, _, _ := session.Pty()
		renderer := bubbletea.MakeRenderer(session)
		mainStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
		actionStyle := renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
		helpStyle := renderer.NewStyle().Border(lipgloss.NormalBorder(), true, false, false)
		inputStyle := renderer.NewStyle().Border(lipgloss.NormalBorder(), true)
		bg := "light"

		if renderer.HasDarkBackground() {
			bg = "dark"
		}

		model := Model{
			version:     "0.0.1-alpha1",
			term:        pty.Term,
			width:       pty.Window.Width,
			height:      pty.Window.Height,
			time:        time.Now(),
			bg:          bg,
			user:        session.Context().Value("user").(User),
			publicKey:   session.Context().Value("publicKey").(string),
			help:        help.New(),
			err:         nil,
			users:       users,
			mainStyle:   mainStyle,
			infoStyle:   infoStyle,
			actionStyle: actionStyle,
			helpStyle:   helpStyle,
			inputStyle:  inputStyle,
		}.WithView(SplashView)

		model.help.Styles.Ellipsis = infoStyle
		model.help.Styles.FullDesc = infoStyle
		model.help.Styles.FullKey = actionStyle
		model.help.Styles.FullSeparator = infoStyle
		model.help.Styles.ShortDesc = infoStyle
		model.help.Styles.ShortKey = actionStyle
		model.help.Styles.ShortSeparator = infoStyle

		options := append(bubbletea.MakeOptions(session), tea.WithAltScreen())

		return program(model, options...)
	}

	return bubbletea.MiddlewareWithProgramHandler(handler, termenv.ANSI256)
}
