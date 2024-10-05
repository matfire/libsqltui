package screens

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type InitScreenState int

const (
	loading InitScreenState = iota
	initialized
)

type InitScreen struct {
	state          InitScreenState
	sqldStatus     int
	loadingSpinner spinner.Model
	error          error
}

type InitEndMsg struct {
	Valid bool
}

func sendInitEndMsg(valid bool) tea.Cmd {
	return func() tea.Msg {
		return InitEndMsg{Valid: valid}
	}
}

func checkSqldServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get("http://127.0.0.1:8080/health")
	if err != nil {
		// There was an error making our request. Wrap the error we received
		// in a message and return it.
		return errMsg{err}
	}
	// We received a response from the server. Return the HTTP status code
	// as a message.
	return statusMsg(res.StatusCode)
}

func (s InitScreen) Init() tea.Cmd {
	return tea.Batch(s.loadingSpinner.Tick, checkSqldServer)
}

func (s InitScreen) View() string {
	if s.error != nil {
		return fmt.Sprintf("got this error %v", s.error.Error())
	}
	switch s.state {
	case loading:
		return fmt.Sprintf("%s checking sqld is running...", s.loadingSpinner.View())
	case initialized:
		return fmt.Sprintf("got status code: %d", s.sqldStatus)
	default:
		return "this should not happen"
	}

}

func (s InitScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case statusMsg:
		log.Printf("got status message value %v", msg)
		s.sqldStatus = int(msg)
		s.state = initialized
		cmd = sendInitEndMsg(msg == 200)
	case errMsg:
		s.error = msg
	case spinner.TickMsg:
		s.loadingSpinner, cmd = s.loadingSpinner.Update(msg)
	}
	cmds = append(cmds, cmd)
	return s, tea.Batch(cmds...)
}

func NewInitScreen() InitScreen {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return InitScreen{
		state:          loading,
		loadingSpinner: s,
	}
}
