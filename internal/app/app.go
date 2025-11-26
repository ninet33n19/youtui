package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ninet33n19/youtui/internal/tui"
)

func Run() error {
	model := tui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()

	return err
}
