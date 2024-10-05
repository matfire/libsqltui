package screens

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matfire/libsqltui/constants"
)

type MainScreen struct {
	actionList list.Model
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m MainScreen) Init() tea.Cmd {
	return nil
}

func (m MainScreen) View() string {
	return m.actionList.View()
}

func (m MainScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return m, constants.SendActionSelectMsg(m.actionList.SelectedItem().(constants.ActionItem))
		}
	case tea.WindowSizeMsg:
		m.actionList.SetSize(msg.Width, msg.Height)
	}
	m.actionList, cmd = m.actionList.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func NewMainScreen() MainScreen {

	var m = MainScreen{
		actionList: list.New(constants.GetActionItems(), list.NewDefaultDelegate(), 80, 24),
	}
	m.actionList.Title = "What do you want to do"
	return m
}
