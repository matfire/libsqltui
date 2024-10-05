package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matfire/libsqltui/constants"
	"github.com/matfire/libsqltui/screens"
)

type sessionState int

const (
	introView sessionState = iota
	mainView
	createView
)

type model struct {
	state        sessionState
	initScreen   screens.InitScreen
	mainScreen   screens.MainScreen
	createScreen screens.CreateScreen
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.initScreen.Init(), m.createScreen.Init())
}

func (m model) View() string {
	switch m.state {
	case introView:
		return m.initScreen.View()
	case mainView:
		return m.mainScreen.View()
	case createView:
		return m.createScreen.View()
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
		switch msg := msg.(type) {
		case screens.InitEndMsg:
			if msg.Valid {
				m.state = mainView
				return m, m.mainScreen.Init()
			}
		}
		newIntroView, newCmd := m.initScreen.Update(msg)
		introModel, ok := newIntroView.(screens.InitScreen)
		if !ok {
			panic("could not perform assertion on init model ui")
		}
		m.initScreen = introModel
		cmd = newCmd
	case mainView:
		switch msg := msg.(type) {
		case constants.ActionSelectMsg:
			if msg.Item.Id == 1 {
				m.state = createView
				return m, nil
			}
		}
		m.mainScreen, cmd = m.mainScreen.Update(msg)
	case createView:
		switch msg.(type) {
		case constants.BackMsg:
			m.state = mainView
			return m, nil
		}
		newCreateView, newCmd := m.createScreen.Update(msg)
		createModel, ok := newCreateView.(screens.CreateScreen)
		if !ok {
			panic("could not perform assertion on create model ui")
		}
		m.createScreen = createModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func initialModel() model {
	return model{
		state:        introView,
		initScreen:   screens.NewInitScreen(),
		mainScreen:   screens.NewMainScreen(),
		createScreen: screens.NewCreateScreen(),
	}
}

func main() {
	tea.LogToFile("debug.log", "debug")
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
