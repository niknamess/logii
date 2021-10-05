package logenc

import (
	"fmt"
	//"go/printer"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Logger *log.Logger
)

func ProcLine(line string) (csvF string) {

	if len(line) == 0 {

		return
	}
	lookFor := "<loglist>"
	xmlline := DecodeLine(line)
	contain := strings.Contains(xmlline, lookFor)
	if contain == false {

		return xmlline
	}

	val, err := DecodeXML(xmlline)
	if err != nil {

		return
	}

	csvline := EncodeCSV(val)
	fmt.Print(csvline)
	return csvline
}

func procLineq(line string) (csvF string) {

	if len(line) == 0 {

		return
	}

	xmlline := DecodeLine(line)
	val, err := DecodeXML(xmlline)
	if err != nil {

		return
	}

	csvline := EncodeCSV(val)
	return csvline
}

func ProcFile(file string) {
	ch := make(chan string, 100)
	log.Println("1")
	for i := runtime.NumCPU() + 1; i > 0; i-- {
		go func() {
			for {
				select {
				case line := <-ch:

					ProcLine(line)
				}
			}

		}()
	}

	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
}

func ProcDir(dir string) {

	filepath.Walk(dir,
		func(path string, file os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !file.IsDir() {

				ProcFile(path)
			}
			return nil
		})
}

func ProcWrite(dir string) {

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

	filew, err1 := os.OpenFile("/home/nik/projects/logs/r/mainlogs1.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err1 != nil {
		log.Fatal(err1)
	}

	Logger = log.New(filew, "", 0)

	ch := make(chan string, 100)

	for i := runtime.NumCPU() + 1; i > 0; i-- {
		go func() {
			for {
				select {
				case line := <-ch:

					Logger.Println(procLineq(line))
				}
			}

		}()
	}

	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
}

func Promrun() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ProcBleve(dir string) {

	//mapping := bleve.NewIndexMapping()
	index, err := bleve.Open("example.bleve")
	//index, err := bleve.New("example.bleve", mapping)
	//index, _ = bleve.Open("example.bleve")
	if err != nil {
		fmt.Println(err)
		return
	}
	data := ProcFileBreve(dir)

	// index some data
	for _, i := range data.XML_RECORD_ROOT {
		index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
		println(i)
	}
	// search for some text
	query := bleve.NewMatchQuery("Service")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}

func ProcLineBleve(line string) (val LogList) {

	if len(line) == 0 {

		return
	}
	xmlline := DecodeLine(line)
	val, err := DecodeXML(xmlline)
	if err != nil {

		return
	}

	return val
}

func ProcFileBreve(file string) (data1 LogList) {
	var data LogList
	ch := make(chan string, 100)

	for i := runtime.NumCPU() + 1; i > 0; i-- {
		go func() {
			for {
				select {
				case line := <-ch:
					data = ProcLineBleve(line)

				}
			}

		}()
	}

	return data
}

func ProcBleveSearch(dir string) {
	index, err := bleve.Open("example.bleve")
	//index, err := bleve.New("example.bleve", mapping)
	//index, _ = bleve.Open("example.bleve")

	// search for some text
	query := bleve.NewMatchQuery("0001GD41BQJ9HBFXMA2QT47H04")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}
