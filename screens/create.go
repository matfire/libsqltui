package screens

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matfire/libsqltui/constants"
)

type CreateScreen struct {
	input textinput.Model
}

func createDatabase(value string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("create database called %s", value)
		res, err := http.Post(fmt.Sprintf("http://127.0.0.1:8081/v1/namespaces/%s/create", value), "application/json", bytes.NewBuffer([]byte("{}")))
		if err != nil {
			log.Printf("Got error %v", err)
			return constants.CreatedMsg{}
		}
		log.Printf("result is %d", res.StatusCode)
		return constants.CreatedMsg{}
	}
}

func (s CreateScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s CreateScreen) View() string {
	return fmt.Sprintf("%s", s.input.View())
}

func (s CreateScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return s, createDatabase(s.input.Value())
		} else {
			break
		}
	case constants.CreatedMsg:
		log.Printf("database has been created")
	}
	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

func NewCreateScreen() CreateScreen {
	input := textinput.New()
	input.Placeholder = "Database Name"
	input.Focus()
	return CreateScreen{
		input: input,
	}
}
