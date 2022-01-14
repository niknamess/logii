package main

import (
	"flag"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	generator "gitlab.topaz-atcs.com/tmcs/logi2/generate_logs"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/terminal"
	"gitlab.topaz-atcs.com/tmcs/logi2/web"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
)

var (
	content string
	timeout = terminal.Model{0, false, 0, 0, 0, false, true}
	status  tea.Model
)

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
	flagGen := flag.Bool("g", false, "generate_logs")
	flagWeb := flag.String("p", "", "web_interface")
	flagProm := flag.String("m", "", "prometheus")
	flagVFC := flag.String("v", "", "vfc")
	flagR := flag.String("r", "", "remove")
	flagMenu := flag.String("x", "", "menu")
	flagServer := flag.Bool("s", false, "server")
	flagClient := flag.Bool("c", false, "client")
	//flagInfo := flag.String("i", "", "info")
	//flagDD := flag.String("o", "", "dd")

	flag.Parse()

	//go logenc.Promrun()

	if len(*flagServ) > 0 {
		fmt.Println("flagServ:", *flagServ)
		RunRPC(*flagServ)
		return
	}

	if len(*flagFile) > 0 {
		logenc.ProcFile(*flagFile)
		return
	}

	if len(*flagDir) > 0 {
		logenc.ProcDir(*flagDir)
		return
	}

	if len(*flagWrite) > 0 {
		logenc.ProcWrite(*flagWrite)
		return
	}

	if *flagGen {
		generator.ProcGenN(10, 200000)
		return
	}

	if len(*flagWeb) > 0 {
		print(*flagWeb)
		web.ProcWeb(*flagWeb)
		return
	}

	if len(*flagProm) > 0 {
		logenc.Promrun(*flagProm)
		return
	}
	if len(*flagVFC) > 0 {
		controllers.VFC(*flagVFC)
		return
	}
	if len(*flagR) > 0 {
		generator.Example()
		return
	}

	if len(*flagMenu) > 0 {
		MainUi()
		return
	}

	if *flagServer {
		Server()
		return
	}

	if *flagClient {
		Client()
		return
	}

}

func MainUi() {
	var test tea.Model
	var st bool
	str, model := terminal.TerminalUi()
	idx, _ := strconv.Atoi(str)
	if model == timeout && idx == 0 {
		test = terminal.Screensaver()
	} else if model != timeout {
		st = terminal.SwitchMenu(idx)

	}
	if test != nil || status != nil || st != false {
		status = nil
		test = nil
		MainUi()
	}
}
