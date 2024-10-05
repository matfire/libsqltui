package screens

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matfire/libsqltui/constants"
)

type screenState int

const (
	inputState screenState = iota
	loadingState
	successState
	errorState
)

type CreateScreen struct {
	input          textinput.Model
	loadingSpinner spinner.Model
	state          screenState
}

func createDatabase(value string) tea.Cmd {
	return func() tea.Msg {
		res, err := http.Post(fmt.Sprintf("http://127.0.0.1:8081/v1/namespaces/%s/create", value), "application/json", bytes.NewBuffer([]byte("{}")))
		if err != nil {
			return constants.CreatedMsg{Status: res.StatusCode}
		}
		return constants.CreatedMsg{Status: res.StatusCode}
	}
}

func (s CreateScreen) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, s.loadingSpinner.Tick)
}

func (s CreateScreen) View() string {
	switch s.state {
	case inputState:
		return fmt.Sprintf("Enter the name of the database you want to create:\n\n%s\n\n%s", s.input.View(), "{esc} to go back")
	case loadingState:
		return fmt.Sprintf("%s sending request to sqld server...", s.loadingSpinner.View())

	default:
		return fmt.Sprintf("This should not happen")
	}
}

func (s CreateScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return s, createDatabase(s.input.Value())
		} else if msg.String() == "esc" {
			return s, constants.SendBackMsg()
		} else {
			break
		}
	case constants.CreatedMsg:
	}
	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

func NewCreateScreen() CreateScreen {
	input := textinput.New()
	input.Placeholder = "Database Name"
	input.Focus()
	s := spinner.New()
	s.Spinner = spinner.Dot
	return CreateScreen{
		input:          input,
		state:          inputState,
		loadingSpinner: s,
	}
}
