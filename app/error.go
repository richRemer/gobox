package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrorModel struct {
	program *tea.Program
	width   int
	height  int
	style   lipgloss.Style
	current error
	pending chan error
	// TODO: keep past errors for user to inspect
}

func (model ErrorModel) Init() tea.Cmd {
	return nil
}

func (model ErrorModel) Update(msg tea.Msg) (ErrorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ClearErrorMsg:
		select {
		case err := <-model.pending:
			model.current = err
			go model.clearErrors()
		default:
			model.current = nil
		}
	case ErrorMsg:
		if model.current == nil {
			model.current = msg.err
			go model.clearErrors()
		} else {
			model.pending <- msg.err
		}
	case ProgramMsg:
		model.program = msg.program
	}

	return model, nil
}

func (model ErrorModel) View() string {
	if model.current == nil {
		return lipgloss.Place(model.width, model.height, 0.5, 0.5, "")
	} else {
		view := model.style.Render(model.current.Error())
		return lipgloss.Place(model.width, model.height, 0.5, 0.5, view)
	}
}

func (model ErrorModel) clearErrors() {
	<-time.After(3 * time.Second)
	model.program.Send(ClearErrorMsg{})
}
