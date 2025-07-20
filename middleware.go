package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"
)

func middleware() wish.Middleware {
	program := func(m tea.Model, options ...tea.ProgramOption) *tea.Program {
		program := tea.NewProgram(m, options...)

		go func() {
			for {
				<-time.After(1 * time.Second)
				program.Send(time.Time(time.Now()))
			}
		}()

		return program
	}

	teaHandler := func(session ssh.Session) *tea.Program {
		pty, _, active := session.Pty()
		renderer := bubbletea.MakeRenderer(session)
		mainStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
		infoStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
		bg := "light"

		if renderer.HasDarkBackground() {
			bg = "dark"
		}

		if !active {
			wish.Fatalln(session, "no active terminal, skipping")
		}

		m := model{
			term:      pty.Term,
			width:     pty.Window.Width,
			height:    pty.Window.Height,
			time:      time.Now(),
			bg:        bg,
			mainStyle: mainStyle,
			infoStyle: infoStyle,
		}

		return program(m, append(bubbletea.MakeOptions(session), tea.WithAltScreen())...)
	}

	return bubbletea.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}
