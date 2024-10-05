package constants

import (
	"github.com/charmbracelet/bubbles/list"
)

type ActionItem struct {
	title, desc string
	Id          int
}

func (i ActionItem) Title() string       { return i.title }
func (i ActionItem) Description() string { return i.desc }
func (i ActionItem) FilterValue() string { return i.title }

func GetActionItems() []list.Item {
	return []list.Item{
		ActionItem{title: "Create", desc: "Create a new Database", Id: 1},
		ActionItem{title: "Fork", desc: "Fork a database to another db", Id: 2},
		ActionItem{title: "Delete", desc: "Delete an existing Database", Id: 3},
	}
}
