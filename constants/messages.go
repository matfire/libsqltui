package constants

import tea "github.com/charmbracelet/bubbletea"

type BackMsg struct{}

type ActionSelectMsg struct {
	Item ActionItem
}

type CreateMsg struct {
	Value string
}

type CreatedMsg struct {
	Status int
}

func SendBackMsg() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}

func SendActionSelectMsg(selectedItem ActionItem) tea.Cmd {
	return func() tea.Msg {
		return ActionSelectMsg{Item: selectedItem}
	}
}
