package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matfire/libsqltui/screens"
	"os"
)

type sessionState int

const (
	introView sessionState = iota
	mainView
	addView
)

type model struct {
	state      sessionState
	initScreen screens.InitScreen
	mainScreen screens.MainScreen
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	switch m.state {
	case introView:
		return m.initScreen.View()
	}
	return "no view defined for this state"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	switch m.state {
	case introView:
		newIntroView, newCmd := m.initScreen.Update(msg)
		introModel, ok := newIntroView.(screens.InitScreen)
		if !ok {
			panic("could not perform assertion on init model ui")
		}
		m.initScreen = introModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func initialModel() model {
	return model{
		state:      introView,
		initScreen: screens.NewInitScreen(),
		mainScreen: screens.MainScreen{},
	}
}

func main() {
	tea.LogToFile("debug.log", "debug")
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
