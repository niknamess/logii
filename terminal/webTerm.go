package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Println("Enter IP:", err)
		os.Exit(1)
	}

	if err := tea.NewProgram(model{}, tea.WithInput(tty)).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	tty.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Now type something: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	fmt.Printf("You entered: %s\n", strings.TrimSpace(text))
}

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	return "Bubble Tea running. Press any key to exit...\n"
}
