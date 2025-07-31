package app

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"modernc.org/sqlite"
)

const MaxErrors = 8

type WrappedError struct {
	message string
	wrapped error
}

func (err *WrappedError) Error() string {
	return err.message
}

func (err *WrappedError) Unwrap() error {
	return err.wrapped
}

func wrap(err error) *WrappedError {
	if err == nil {
		return nil
	} else if _, ok := err.(*WrappedError); ok {
		return err.(*WrappedError)
	}

	var message string = "unknown error"

	if sqlerr, ok := err.(*sqlite.Error); ok {
		message = "database error"

		switch sqlerr.Code() {
		case 2067:
			if strings.Contains(sqlerr.Error(), "user.name") {
				message = "nick already registered"
			} else if strings.Contains(sqlerr.Error(), "public_key.pem") {
				message = "key in use"
			}
		}
	} else if errors.Is(err, sql.ErrConnDone) {
		message = "database connection closed"
	} else if errors.Is(err, sql.ErrNoRows) {
		message = "no results"
	} else if errors.Is(err, sql.ErrTxDone) {
		message = "database transaction closed"
	}

	return &WrappedError{message: message, wrapped: err}
}

type ErrorModel struct {
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
			return model, tick()
		default:
			model.current = nil
		}
	case ErrorMsg:
		if model.current == nil {
			model.current = wrap(msg.err)
			return model, tick()
		} else {
			model.pending <- wrap(msg.err)
		}
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

func tick() tea.Cmd {
	return func() tea.Msg {
		<-time.After(3 * time.Second)
		return ClearErrorMsg{}
	}
}
