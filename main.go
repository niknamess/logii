package main

import (
	"flag"
	"fmt"

	"sort"

	generator "gitlab.topaz-atcs.com/tmcs/logi2/generate_logs"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/web"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
)

func main() {
	flagFile := flag.String("f", "", "parse log file")
	//flagFile := kingpin.Arg("f", "Directory path(s) to look for files").Default("./logtest/test/22-06-2021").String()
	flagDir := flag.String("d", "", "parse dir")
	flagSearch := flag.String("s", "", "search")
	flagServ := flag.String("z", "", "server")
	flagWrite := flag.String("w", "", "write_logs")
	flagGen := flag.String("g", "", "generate_logs")
	flagWeb := flag.String("p", "", "web_interface")
	flagTest := flag.String("c", "", "web_interface and generate log")
	flagProm := flag.String("m", "", "prometheus")
	flagVFC := flag.String("v", "", "vfc")
	flagR := flag.String("r", "", "remove")
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
		//logenc.ProcBleveSearch(*flagBleveSearch)
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

	if len(*flagTest) > 0 {
		// TODO same as flagGen
		generator.ProcGenN()

	}

	if len(*flagSearch) > 0 {
		var text string
		var limit int

		var MassStr []logenc.Data

		fmt.Print("Enter limit: ")
		fmt.Scanln(&limit)
		fmt.Print("Enter text: ")
		fmt.Scanln(&text)

		chRes := make(chan logenc.Data, 100)
		go func() {
			scan := &logenc.Scan{}
			scan.Find = *flagSearch
			scan.Text = text
			scan.ChRes = chRes
			scan.LimitResLines = limit
			scan.Search()
			close(scan.ChRes)
		}()

	ext:
		for i := 0; i < limit; i++ {

			select {

			case data, ok := <-chRes:
				if !ok {
					break ext
				}
				MassStr = append(MassStr, data)

			}
		}
		sort.Slice(MassStr, func(i, j int) (less bool) {
			return MassStr[i].ID < MassStr[j].ID
		})
		fmt.Printf("%+v\n", MassStr)
		return
	}
}
