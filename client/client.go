package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
	"unicode"

	"github.com/SCU-SJL/menuscreen"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"gitlab.topaz-atcs.com/tmcs/logi2/terminal"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("From server:", string(buf[0:n])) //From server
	}
}

func Client() {

	terminal.CallClear()
	fmt.Println("Now you can use only VFC(stable) or WEB(stable) and stop this service")
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go reader(c)
	for {
		//reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")

		//text, _ := reader.ReadString('\n') //Send server
		//ui terminal
		idx, text := menuClientMain()
		_, err := c.Write([]byte(text)) //Send server
		if err != nil {
			log.Fatal("write error:", err)
			break
		}

		if idx == 0 {
			for {
				text = webMenu()

				_, err := c.Write([]byte(text)) //Send server
				if err != nil {
					log.Fatal("write error:", err)
					break
				}
				if text == "stop" {
					break
				}
				time.Sleep(1e9)
			}
		}
		if idx == 4 {
			return
		}

		time.Sleep(1e9)
	}
}

func menuClientMain() (int, string) {
	menu, err := menuscreen.NewMenuScreen()
	if err != nil {
		panic(err)
	}
	defer menu.Fini()
	menu.SetTitle("ControlPanel").
		SetLine(0, "WEB").
		SetLine(1, "VFC").
		SetLine(2, "STOPWEB").
		SetLine(3, "STOPVFC").
		SetLine(4, "STOP CLIENT").
		SetLine(5, "STOP SERVER").
		Start()
	idx, ln, ok := menu.ChosenLine()
	if !ok {
		fmt.Println("you did not chose any items.")
		return idx, ln

	}
	fmt.Printf("you've chosen %d line, content is: %s\n", idx, ln)
	return idx, ln
}
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
func webMenu() string {
	p := tea.NewProgram(initialModel())

	tsk, err := p.StartReturningModel()
	if err != nil {
		log.Fatal(err)
	}

	str := tsk.View()
	split := strings.Split(str, ">")
	str = strings.Join(split, " ")
	split = strings.Split(str, " ")
	str = strings.Join(split, "")
	str = normalform8(str)
	//fmt.Print("Adress:", str)
	return str
}

func normalform8(s string) string {
	if last := len(s) - 8; last >= 0 {
		s = s[:last]
	}
	return s
}

type tickMsg struct{}
type errMsg error

type model struct {
	textInput textinput.Model
	err       error
}

func initialModel() model {
	ti := textinput.NewModel()
	ti.Placeholder = "192.168.0.1"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(

		m.textInput.View())

}
