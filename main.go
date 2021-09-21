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
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

var (
	Logger *log.Logger
)

func procLine(line string) (csvF string) {

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
	return csvline
}

func procLineq(line string) (csvF string) {

	if len(line) == 0 {

		return
	}

	xmlline := logenc.DecodeLine(line)
	val, err := logenc.DecodeXML(xmlline)
	if err != nil {

		return
	}

	csvline := logenc.EncodeCSV(val)
	return csvline
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

func procWrite(dir string) {

	filepath.Walk(dir,
		func(path string, file os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !file.IsDir() {

				procFileWrite(path)
			}
			return nil
		})
}

func procFileWrite(file string) {

	filew, err1 := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err1 != nil {
		log.Fatal(err1)
	}

	//Logger = log.New(filew, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger = log.New(filew, "TEST: ", log.Ldate|log.Ltime)

	ch := make(chan string, 100)

	for i := runtime.NumCPU() + 1; i > 0; i-- {
		go func() {
			for {
				select {
				case line := <-ch:

					//procLine(line)
					Logger.Println(procLineq(line))
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

func promrun() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	flagFile := flag.String("f", "", "parse log file")
	flagDir := flag.String("d", "", "parse dir")
	flagSearch := flag.String("s", "", "search")
	flagServ := flag.String("z", "", "server")
	flagWrite := flag.String("w", "", "write_logs")
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

	if len(*flagWrite) > 0 {

		procWrite(*flagWrite)
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
