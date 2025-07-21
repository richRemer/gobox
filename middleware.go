package main

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

func middleware() wish.Middleware {
	program := func(model tea.Model, options ...tea.ProgramOption) *tea.Program {
		program := tea.NewProgram(model, options...)

		go func() {
			for {
				<-time.After(1 * time.Second)
				program.Send(time.Time(time.Now()))
			}
		}()

		return program
	}

	handler := func(session ssh.Session) *tea.Program {
		pty, _, _ := session.Pty()
		renderer := bubbletea.MakeRenderer(session)
		keyStyle := renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
		mainStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
		helpStyle := renderer.NewStyle().Border(lipgloss.NormalBorder(), true, false, false)
		bg := "light"

		if renderer.HasDarkBackground() {
			bg = "dark"
		}

		model := application{
			term:      pty.Term,
			width:     pty.Window.Width,
			height:    pty.Window.Height,
			time:      time.Now(),
			bg:        bg,
			keys:      activeKeyMap,
			help:      help.New(),
			mainStyle: mainStyle,
			infoStyle: infoStyle,
			helpStyle: helpStyle,
		}

		model.help.Styles.Ellipsis = infoStyle
		model.help.Styles.FullDesc = infoStyle
		model.help.Styles.FullKey = keyStyle
		model.help.Styles.FullSeparator = infoStyle
		model.help.Styles.ShortDesc = infoStyle
		model.help.Styles.ShortKey = keyStyle
		model.help.Styles.ShortSeparator = infoStyle

		options := append(bubbletea.MakeOptions(session), tea.WithAltScreen())

		return program(model, options...)
	}

	return bubbletea.MiddlewareWithProgramHandler(handler, termenv.ANSI256)
}
