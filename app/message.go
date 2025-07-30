package app

import tea "github.com/charmbracelet/bubbletea"

type ClearErrorMsg struct{}

type CloseSplashMsg struct{}

type ErrorMsg struct {
	err error
}

type ProgramMsg struct {
	program *tea.Program
}
