package deletescreen

import tea "github.com/charmbracelet/bubbletea"

type DeleteScreen struct {
}

func (s DeleteScreen) Init() tea.Cmd {
	return nil
}

func (s DeleteScreen) View() string {
	return "in delete screen"
}

func (s DeleteScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func NewDeleteScreen() DeleteScreen {
	return DeleteScreen{}
}
