package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.topaz-atcs.com/tmcs/logi/logenc"
)

func procLine(line string) {

	if len(line) == 0 {

		return
	}

	xmlline := logenc.DecodeLine(line)
	val, err := logenc.DecodeXML(xmlline)
	if err != nil {

		return
	}

	csvline := logenc.EncodeCSV(val)
	fmt.Print(csvline)
}

func procFile(file string) {
	ch := make(chan string, 100)

	for i := runtime.NumCPU() + 1; i > 0; i-- {
		go func() {
			for {
				select {
				case line := <-ch:

					procLine(line)
				}
			}

		}()
	}

	err := logenc.ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
}

func procDir(dir string) {

	filepath.Walk(dir,
		func(path string, file os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !file.IsDir() {

				procFile(path)
			}
			return nil
		})
}

func promrun() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	flagFile := flag.String("f", "", "parse log file")
	flagDir := flag.String("d", "", "parse dir")
	flagSearch := flag.String("s", "", "search")
	flagServ := flag.String("z", "", "server")
	flag.Parse()

	go promrun()

	if len(*flagServ) > 0 {
		fmt.Println("flagServ:", *flagServ)
		RunRPC(*flagServ)
		return
	}

	if len(*flagFile) > 0 {

		procFile(*flagFile)
		return
	}

	if len(*flagDir) > 0 {

		procDir(*flagDir)
		return
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
