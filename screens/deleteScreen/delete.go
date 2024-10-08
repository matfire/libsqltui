package deletescreen

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matfire/libsqltui/constants"
)

type DeleteMsg struct{}

type DeletedMsg struct {
	success bool
}

type screenState int

const (
	inputState screenState = iota
	confirmState
	loadingState
	successState
	errorState
)

type DeleteScreen struct {
	state    screenState
	input    textinput.Model
	loader   spinner.Model
	adminUrl string
}

func deleteDatabase(adminUrl string, name string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/namespaces/%s", adminUrl, name), nil)
		if err != nil {
			return DeletedMsg{success: false}
		}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode > 299 {
			return DeletedMsg{success: false}
		}
		return DeletedMsg{success: true}
	}
}

func (s DeleteScreen) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, s.loader.Tick)
}

func (s DeleteScreen) View() string {
	switch s.state {
	case inputState:
		return fmt.Sprintf("Enter the name of the database you want to delete\n\n%s\n\n%s", s.input.View(), "{esc} to go back")
	case confirmState:
		return fmt.Sprintf("Are you sure you want to delete the database named %s\n\nType {Y}es or {N]o", s.input.Value())
	case successState:
		return "db deleted"
	case errorState:
		return "could not delete"
	case loadingState:
		return fmt.Sprintf("%s Sending request to sqld server...", s.loader.View())
	}
	return "you should not be seeing this"
}

func (s DeleteScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" && s.state == inputState {
			s.state = confirmState
			return s, nil
		}
		if key := strings.ToLower(msg.String()); (key == "y" || key == "n") && s.state == confirmState {
			if key == "y" {
				s.state = loadingState
				return s, deleteDatabase(s.adminUrl, s.input.Value())
			} else {
				s.state = inputState
				return s, nil
			}
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
	case spinner.TickMsg:
		s.loader, cmd = s.loader.Update(msg)
		cmds = append(cmds, cmd)
	}
	switch s.state {
	case inputState:
		s.input.Focus()
		var cmd tea.Cmd
		s.input, cmd = s.input.Update(msg)
		cmds = append(cmds, cmd)
	}
	return s, tea.Batch(cmds...)
}

func NewDeleteScreen(adminUrl string) DeleteScreen {
	i := textinput.New()
	i.Placeholder = "enter database name to delete"
	loader := spinner.New()
	loader.Spinner = spinner.Dot
	loader.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return DeleteScreen{state: inputState, input: i, loader: loader, adminUrl: adminUrl}
}
