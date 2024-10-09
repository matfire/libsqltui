package forkscreen

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matfire/libsqltui/constants"
)

type forkScreenStatus int

type ForkedMsg struct {
	success bool
}

const (
	fromInputState forkScreenStatus = iota
	toInputState
	loadingState
	successState
	errorState
)

type ForkScreen struct {
	state          forkScreenStatus
	adminUrl       string
	loadingSpinner spinner.Model
	fromInput      textinput.Model
	toInput        textinput.Model
}

func forkNamespace(adminUrl string, from string, to string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 10 * time.Second}
		res, err := client.Post(fmt.Sprintf("%s/v1/namespaces/%s/fork/%s", adminUrl, from, to), "application/json", bytes.NewBuffer([]byte("{}")))
		if err != nil {
			return ForkedMsg{success: false}
		}
		return ForkedMsg{success: res.StatusCode < 300}
	}
}

func (s ForkScreen) Init() tea.Cmd {
	return tea.Batch(s.loadingSpinner.Tick, textinput.Blink)
}

func (s ForkScreen) View() string {
	switch s.state {
	case fromInputState:
		return fmt.Sprintf("Enter the name of the database you want to fork:\n\n%s\n\n%s", s.fromInput.View(), "{esc} to go back")
	case toInputState:
		return fmt.Sprintf("Enter the name of the to be created fork:\n\n%s\n\n%s", s.toInput.View(), "{esc} to go back")
	case successState:
		return fmt.Sprintf("database forked succesfully")
	case errorState:
		return fmt.Sprintf("database could not be forked")
	case loadingState:
		return fmt.Sprintf("%s sending request to sqld server...", s.loadingSpinner.View())
	}
	return "you should never see this"
}

func (s ForkScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if s.state == fromInputState {
				s.state = toInputState
				s.toInput.Focus()
				return s, nil
			} else if s.state == toInputState {
				s.state = loadingState
				return s, forkNamespace(s.adminUrl, s.fromInput.Value(), s.toInput.Value())
			}
			break
		}
		if msg.String() == "esc" && s.state != loadingState {
			if s.state == fromInputState {
				return s, constants.SendBackMsg()
			}
			s.state = fromInputState
			s.fromInput.Focus()
			return s, nil
		}
		break
	case ForkedMsg:
		if msg.success {
			s.state = successState
		} else {
			s.state = errorState
		}
		return s, nil
	case spinner.TickMsg:
		s.loadingSpinner, cmd = s.loadingSpinner.Update(msg)
		return s, cmd
	}
	switch s.state {
	case fromInputState:
		s.fromInput, cmd = s.fromInput.Update(msg)
		return s, cmd
	case toInputState:
		s.toInput, cmd = s.toInput.Update(msg)
		return s, cmd
	}
	return s, nil
}

func NewForkScreen(adminUrl string) ForkScreen {
	s := spinner.New()
	s.Spinner = spinner.Dot
	fromInput := textinput.New()
	fromInput.Placeholder = "Which database do you want to fork?"
	toInput := textinput.New()
	toInput.Placeholder = "What name should de forked database have?"
	fromInput.Focus()
	return ForkScreen{
		state:          fromInputState,
		adminUrl:       adminUrl,
		loadingSpinner: s,
		fromInput:      fromInput,
		toInput:        toInput,
	}
}
