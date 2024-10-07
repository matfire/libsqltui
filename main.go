package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matfire/libsqltui/constants"
	createscreen "github.com/matfire/libsqltui/screens/createScreen"
	deletescreen "github.com/matfire/libsqltui/screens/deleteScreen"
	initscreen "github.com/matfire/libsqltui/screens/initScreen"
	mainscreen "github.com/matfire/libsqltui/screens/mainScreen"
)

type sessionState int

const (
	introView sessionState = iota
	mainView
	createView
	deleteView
)

type model struct {
	state        sessionState
	initScreen   initscreen.InitScreen
	mainScreen   mainscreen.MainScreen
	createScreen createscreen.CreateScreen
	deleteScreen deletescreen.DeleteScreen
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.initScreen.Init(), m.createScreen.Init(), m.deleteScreen.Init())
}

func (m model) View() string {
	switch m.state {
	case introView:
		return m.initScreen.View()
	case mainView:
		return m.mainScreen.View()
	case createView:
		return m.createScreen.View()
	case deleteView:
		return m.deleteScreen.View()
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
		case initscreen.InitEndMsg:
			if msg.Valid {
				m.state = mainView
				return m, m.mainScreen.Init()
			}
		}
		newIntroView, newCmd := m.initScreen.Update(msg)
		introModel, ok := newIntroView.(initscreen.InitScreen)
		if !ok {
			panic("could not perform assertion on init model ui")
		}
		m.initScreen = introModel
		cmd = newCmd
	case mainView:
		switch msg := msg.(type) {
		case constants.ActionSelectMsg:
			switch msg.Item.Id {
			case 1:
				m.state = createView
			case 3:
				m.state = deleteView
			}
			return m, nil
		}
		m.mainScreen, cmd = m.mainScreen.Update(msg)
	case createView:
		switch msg.(type) {
		case constants.BackMsg:
			m.state = mainView
			return m, nil
		}
		newCreateView, newCmd := m.createScreen.Update(msg)
		createModel, ok := newCreateView.(createscreen.CreateScreen)
		if !ok {
			panic("could not perform assertion on create model ui")
		}
		m.createScreen = createModel
		cmd = newCmd
	case deleteView:
		switch msg.(type) {
		case constants.BackMsg:
			m.state = mainView
			return m, nil
		}
		newDeleteView, newCmd := m.deleteScreen.Update(msg)
		deleteModel, ok := newDeleteView.(deletescreen.DeleteScreen)
		if !ok {
			panic("could not perform assertion on delete model ui")
		}
		m.deleteScreen = deleteModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func initialModel() model {
	return model{
		state:        introView,
		initScreen:   initscreen.NewInitScreen(),
		mainScreen:   mainscreen.NewMainScreen(),
		createScreen: createscreen.NewCreateScreen(),
		deleteScreen: deletescreen.NewDeleteScreen(),
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
