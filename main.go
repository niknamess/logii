package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/SCU-SJL/menuscreen"
	tea "github.com/charmbracelet/bubbletea"
	generator "gitlab.topaz-atcs.com/tmcs/logi2/generate_logs"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/terminal"
	"gitlab.topaz-atcs.com/tmcs/logi2/web"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
)

var content string

// playType indicates how to play a gauge.
type playType int

const (
	playTypePercent playType = iota
	playTypeAbsolute
)

func main() {
	//content := "nope"
	flagFile := flag.String("f", "", "parse log file")
	//flagFile := kingpin.Arg("f", "Directory path(s) to look for files").Default("./logtest/test/22-06-2021").String()
	flagDir := flag.String("d", "", "parse dir")
	//flagSearch := flag.String("s", "", "search")
	flagServ := flag.String("z", "", "server")
	flagWrite := flag.String("w", "", "write_logs")
	flagGen := flag.String("g", "", "generate_logs")
	flagWeb := flag.String("p", "", "web_interface")
	flagProm := flag.String("m", "", "prometheus")
	flagVFC := flag.String("v", "", "vfc")
	flagR := flag.String("r", "", "remove")
	flagMenu := flag.String("x", "", "menu")
	//flagInfo := flag.String("i", "", "info")
	//flagDD := flag.String("o", "", "dd")

	flag.Parse()

	//go logenc.Promrun()

	if len(*flagServ) > 0 {
		fmt.Println("flagServ:", *flagServ)
		RunRPC(*flagServ)

	}

	if len(*flagFile) > 0 {

		logenc.ProcFile(*flagFile)
	}

	if len(*flagDir) > 0 {

		logenc.ProcDir(*flagDir)
	}

	if len(*flagWrite) > 0 {

		logenc.ProcWrite(*flagWrite)
	}

	if len(*flagGen) > 0 {

		generator.ProcGenN()
	}

	if len(*flagWeb) > 0 {
		print(*flagWeb)
		web.ProcWeb(*flagWeb)

	}

	if len(*flagProm) > 0 {
		logenc.Promrun(*flagProm)

	}
	if len(*flagVFC) > 0 {
		controllers.VFC(*flagVFC)

	}
	if len(*flagR) > 0 {

		generator.Example()
	}

	if len(*flagMenu) > 0 {

		TerminalUi()

	}

}

func TerminalUi() {
	/* tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Println("Could not open TTY:", err)
		os.Exit(1)
	} */

	initialModel := terminal.Model{0, false, 100, 0, 0, false, false}
	p := tea.NewProgram(initialModel)
	model, err := p.StartReturningModel()
	if err != nil {
		fmt.Println("could not start program:", err)
	}
	fmt.Print(model)
	model.Init()
	str := model.View()
	fmt.Println(str)
	idx, err := strconv.Atoi(str)
	switch choose := idx; choose {
	case 0:
		files, err := ioutil.ReadDir("./repdata/")
		if err != nil {
			log.Fatal(err)
		}
		//menu, err := menuscreen.NewMenuScreen()
		//if err != nil {
		//	panic(err)
		//}
		//defer menu.Fini()
		//menu.SetTitle("Menu").

		for i, file := range files {
			//	menu.SetTitle("").
			//		SetLine(i, "Decode file with logs").
			fmt.Println(i, file)

		}
		//Start()
		//idx, ln, _ := menu.ChosenLine()
		fmt.Print("Enter content for ProcFile:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		logenc.ProcFile(text)
	case 1:
		//fmt.Print("Enter content for flag ProcDir:")
		//reader := bufio.NewReader(os.Stdin)
		//text, _ := reader.ReadString('\n')
		logenc.ProcDir("./repdata/")
	case 2:
		fmt.Print("Enter content for flag:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		logenc.ProcWrite(text)
	case 3:
		generator.ProcGenN()
	case 4:
		fmt.Print("Enter port for run Web:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		web.ProcWeb(text)
	case 5:
		fmt.Print("Enter content for Prometheus:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		logenc.Promrun(text)
	case 6:
		controllers.VFC("10015")
	case 7:
		generator.Example()
	case 8:
		fmt.Print("Enter content for Search:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		logenc.SearchT(text)
	}
	//model.Update()
	//tty.Close()

	/* reader := bufio.NewReader(os.Stdin)
	fmt.Print("Now type something: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	fmt.Printf("You entered: %s\n", strings.TrimSpace(text)) */
}

func Menu() {
	menu, err := menuscreen.NewMenuScreen()
	if err != nil {
		panic(err)
	}
	defer menu.Fini()
	menu.SetTitle("Menu").
		SetLine(0, "Decode file with logs").
		SetLine(1, "Decode dir with file logs").
		SetLine(2, "Write decoded logs").
		SetLine(3, "Gen logs").
		SetLine(4, "Run Web").
		SetLine(5, "Run Ptometheus").
		SetLine(6, "running VFC").
		SetLine(7, "clear genlogs").
		SetLine(8, "Search word or collocation").
		Start()
	idx, ln, _ := menu.ChosenLine()

	fmt.Printf("you've chosen %d line, content is: %s\n", idx, ln)
	switch choose := idx; choose {
	case 0:
		files, err := ioutil.ReadDir("./repdata/")
		if err != nil {
			log.Fatal(err)
		}
		//menu, err := menuscreen.NewMenuScreen()
		//if err != nil {
		//	panic(err)
		//}
		//defer menu.Fini()
		//menu.SetTitle("Menu").

		for i, file := range files {
			//	menu.SetTitle("").
			//		SetLine(i, "Decode file with logs").
			fmt.Println(i, file)

		}
		//Start()
		//idx, ln, _ := menu.ChosenLine()
		fmt.Print("Enter content for ProcFile:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		logenc.ProcFile(text)
	case 1:
		//fmt.Print("Enter content for flag ProcDir:")
		//reader := bufio.NewReader(os.Stdin)
		//text, _ := reader.ReadString('\n')
		logenc.ProcDir("./repdata/")
	case 2:
		fmt.Print("Enter content for flag:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		logenc.ProcWrite(text)
	case 3:
		generator.ProcGenN()
	case 4:
		fmt.Print("Enter port for run Web:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		web.ProcWeb(text)
	case 5:
		fmt.Print("Enter content for Prometheus:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		logenc.Promrun(text)
	case 6:
		controllers.VFC("10015")
	case 7:
		generator.Example()
	case 8:
		fmt.Print("Enter content for Search:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		logenc.SearchT(text)
	}
}
