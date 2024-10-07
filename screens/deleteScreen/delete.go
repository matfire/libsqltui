package deletescreen

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matfire/libsqltui/constants"
)

type screenState int

const (
	inputState screenState = iota
	confirmState
	loadingState
	successState
	errorState
)

type DeleteScreen struct {
	state screenState
	input textinput.Model
}

func (s DeleteScreen) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}

func (s DeleteScreen) View() string {
	switch s.state {
	case inputState:
		return fmt.Sprintf("Enter the name of the database you want to delete\n\n%s\n\n%s", s.input.View(), "{esc} to go back")
	case confirmState:
		return "are you sure"
	case successState:
		return "db deleted"
	case errorState:
		return "could not delete"
	}
	return "you should not be seeing this"
}

func (s DeleteScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" && s.state == inputState {
			s.state = confirmState
			return s, nil
		}
		if msg.String() == "esc" {
			if s.state == successState || s.state == errorState {
				s.state = inputState
				s.input.Reset()
				return s, nil
			}
			return s, constants.SendBackMsg()
		}
		break
	}
	switch s.state {
	case inputState:
		var cmd tea.Cmd
		s.input, cmd = s.input.Update(msg)
		return s, cmd
	}
	return s, nil
}

func NewDeleteScreen() DeleteScreen {
	i := textinput.New()
	i.Placeholder = "enter database name to delete"
	i.Focus()
	return DeleteScreen{state: inputState, input: i}
}
