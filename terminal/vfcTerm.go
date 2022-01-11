package terminal

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Vmodel struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialModel() Vmodel {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Vmodel{spinner: s}
}

func (m Vmodel) Init() tea.Cmd {
	return spinner.Tick
}

func (m Vmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "m":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m Vmodel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s VFC server is working port \"10015\"...press Double CTR+C to stop service\n\n If you want return to menu ...press m\n", m.spinner.View())

	if m.quitting {
		return str + "\n"
	}
	return str
}

func VFCTerm() tea.Model {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	model, err := p.StartReturningModel()
	if err != nil {
		fmt.Println("Error starting Bubble Tea program:", err)
		os.Exit(1)
	}
	return model

}
