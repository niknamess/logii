package main

import (
	"flag"
	"fmt"

	"sort"

	generator "gitlab.topaz-atcs.com/tmcs/logi2/generate_logs"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/web"
)

//var flagFile = flag.String("f", "./logtest/test/22-06-2021", "parse log file")
//var flagDir = flag.String("d", "./logtest/test/", "parse dir")
//var flagSearch = flag.String("s", "", "search")
//var flagServ = flag.String("z", "", "server")
//var flagWrite = flag.String("w", "./logtest/test/22-06-2021", "write_logs")
//var flagGen = flag.String("g", "", "generate_logs")
//var flagWeb = flag.String("p", "15000", "web_interface")
//var flagTest = flag.String("c", "15000", "web_interface and generate log")
//var flagBleve = flag.String("b", "./logtest/test/22-06-2021", "Bleve on bleve file")
//var flagBleveSearch = flag.String("bs", "", "Bleve")

func main() {
	flagFile := flag.String("f", "", "parse log file")
	flagDir := flag.String("d", "", "parse dir")
	flagSearch := flag.String("s", "", "search")
	flagServ := flag.String("z", "", "server")
	flagWrite := flag.String("w", "", "write_logs")
	flagGen := flag.String("g", "", "generate_logs")
	flagWeb := flag.String("p", "15000", "web_interface")
	flagTest := flag.String("c", "15000", "web_interface and generate log")
	flagBleve := flag.String("b", "", "Bleve on bleve file")
	flagBleveSearch := flag.String("k", "", "Bleve search")
	flag.Parse()

	go logenc.Promrun()

	if len(*flagServ) > 0 {
		fmt.Println("flagServ:", *flagServ)
		RunRPC(*flagServ)
		//return
	}

	if len(*flagFile) > 0 {

		logenc.ProcFile(*flagFile)
		//logenc.ProcBleveSearch(*flagBleveSearch)
		//return
	}

	if len(*flagDir) > 0 {

		logenc.ProcDir(*flagDir)
		//return
	}

	if len(*flagWrite) > 0 {

		logenc.ProcWrite(*flagWrite)

		//return
	}

	if len(*flagGen) > 0 {

		generator.ProcGenN()
		//return
	}

	if len(*flagGen) > 0 {
		web.ProcWeb(*flagWeb)
	}
	//return
	if len(*flagBleve) > 0 {

		logenc.ProcFileBreve(*flagBleve)

	}
	if len(*flagBleveSearch) > 0 {

		logenc.ProcBleveSearch(*flagBleveSearch)

	}
	if len(*flagTest) > 0 {
		//go web.ProcWeb(*flagTest)
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
