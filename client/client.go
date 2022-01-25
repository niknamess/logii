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
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, str)
	//fmt.Print("Adress:", str)
	return result
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

/* func Clien{t() {
	//reader := bufio.NewReader(os.Stdin)
	//ipaddress, _ := reader.ReadString('\n')
	//util.CheckIPAddress(ipaddress)
	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.1:8888")
	terminal.CallClear()
	fmt.Println("Now you can use only VFC(stable) and WEB(in work)")
	for {
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print("Message from server: " + message)
		//terminal.CallClear()
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		//mt.Println("Now you can use only VFC(stable) and WEB(in work)")
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// Отправляем в socket
		fmt.Fprintf(conn, text+"\n")
		// Прослушиваем ответ
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
*/
