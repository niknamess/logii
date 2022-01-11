package terminal

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
	generator "gitlab.topaz-atcs.com/tmcs/logi2/generate_logs"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
)

// General stuff for styling the view
var (
	end           = 1
	term          = termenv.ColorProfile()
	keyword       = makeFgStyle("211")
	subtle        = makeFgStyle("241")
	progressEmpty = subtle(progressEmptyChar)
	dot           = colorFg(" • ", "236")

	// Gradient colors we'll use for the progress bar
	ramp = makeRamp("#B14FFF", "#00FFA3", progressBarWidth)
)

type tickMsg struct{}
type frameMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type Model struct {
	Choice   int
	Chosen   bool
	Ticks    int
	Frames   int
	Progress float64
	Loaded   bool
	Quitting bool
}

func (m Model) Init() tea.Cmd {
	return tick()
}

// Main update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m Model) View() string {
	var s string
	if m.Quitting {
		//return m, tea.Quit
		str := strconv.Itoa(m.Choice)
		return str
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice += 1
			if m.Choice > 8 {
				m.Choice = 8
			}
		case "k", "up":
			m.Choice -= 1
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			return m, frame()
		}

	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks -= 1
		return m, tick()
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	//var s string
	switch msg.(type) {

	case frameMsg:
		if !m.Loaded {
			m.Frames += 1
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				m.Ticks = 5
				return m, tick()
			}
			return m, frame()
		}

	case tickMsg:
		if m.Loaded {
			if m.Ticks == 0 {
				m.Quitting = true
				//View()
				return m, tea.Quit
			}
			m.Ticks -= 1
			return m, tick()
		}
	}

	return m, nil
}

// Sub-views

// The first view, where you're choosing a task
func choicesView(m Model) string {
	c := m.Choice

	tpl := "Control panel\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		checkbox("Decode file with logs", c == 0),
		checkbox("Decode dir with file logs", c == 1),
		checkbox("Write decoded logs", c == 2),
		checkbox("Gen logs", c == 3),
		checkbox("Run Web", c == 4),
		checkbox("Run Ptometheus", c == 5),
		checkbox("running VFC", c == 6),
		checkbox("clear genlogs", c == 7),
		checkbox("Search word or collocation", c == 8),
	)

	return fmt.Sprintf(tpl, choices, colorFg(strconv.Itoa(m.Ticks), "79"))
}

// The second view, after a task has been chosen
func chosenView(m Model) string {
	var msg string

	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("Decode file with logs\n\n Run ProcDir in %s...", keyword("filename"))
		//end = 0
	case 1:
		msg = fmt.Sprintf("Decode dir with file logs %s...", keyword("dirname"))
		//end = 0
	case 2:
		msg = fmt.Sprintf("Write decoded logs\n\n Okay, cool\n Enter filename -  %s.", keyword("filename"))
		//end = 0
	case 3:
		msg = fmt.Sprintf("GenLogs\n\nCool, we generate logs %s and %s...", keyword("size generate logs"), keyword("Count generate logs"))
		//end = generator.ProcGenN()
	case 4:
		port := "15000"
		//web.ProcWeb(text)
		msg = fmt.Sprintf("Run Web\n\n Start web interface ...%s.", keyword(port))

		//Добавить таймер и выход из bubbltea
		//go web.ProcWeb(port)
	case 5:
		msg = fmt.Sprintf("Run Ptometheus\n\nOkay, cool, then we’ll need a start new service.")
	case 6:
		//controllers.VFC("10015")
		//controllers.VFC("10015")
		msg = fmt.Sprintf("running VFC\n\n We start VFC service  %s ...", keyword("OK"))
		//controllers.VFC("10015")
		//end = 0
	case 7:
		//generator.Example()
		msg = fmt.Sprintf("Clear genlogs\n\n Please wait, we clear generated...")
		generator.Example()
		//end = 0

	case 8:
		fmt.Print("Enter content for Search:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		//logenc.SearchT(text)
		msg = fmt.Sprintf("Search word or collocation\n\nPlease enter word or collocation  %s...", keyword(text))
		logenc.SearchT(text)
	default:
		msg = fmt.Sprintf("Okay.\n\nYou enter the error please restart program /n/n Report a bug in %s or %s...", keyword("Contact 1"), keyword("Contact 2"))
	}

	label := "Loading..."
	if m.Loaded {
		label = fmt.Sprintf("Loaded. Following a %s seconds...", colorFg(strconv.Itoa(m.Ticks), "79"))
	}

	return msg + "\n\n" + label + "\n" + progressbar(80, m.Progress) + "%"
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(width int, percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
