package logenc

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

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
	close(ch)
}

func Promrun() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
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

func ProcFileBreve(file string) {
	var wg sync.WaitGroup
	var counter int32 = 0
	var data LogList
	metaname := "example.bleve"
	index, err := bleve.Open(metaname)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(metaname, mapping)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	// search for some text
	ch := make(chan string, 100)

	for i := runtime.NumCPU() + 1; i > 0; i-- {
		go func() {
			wg.Add(1)
			defer wg.Done()

		brloop:
			for {
				select {
				case line, ok := <-ch:
					if !ok {
						break brloop
					}
					data = ProcLineBleve(line)
					//fmt.Println((data.XML_RECORD_ROOT))
					//fmt.Println(len(data.XML_RECORD_ROOT))
					atomic.AddInt32(&counter, 1)
					println(counter)
					if len(data.XML_RECORD_ROOT) > 0 {
						index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
					}
				}
			}

		}()

	}

	err = ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
	close(ch)
	wg.Wait()
}

//func ProcBleveSearch(dir string) {
func ProcBleveSearch(dir string) (searchRes string) {

	index, _ := bleve.Open("example.bleve")
	//query := bleve.NewFuzzyQuery(dir)
	query := bleve.NewMatchQuery(dir)
	query.Fuzziness = 1
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)

	//query := bleve.NewMatchQuery("0001GD2DVH34EV686NP4W3BHX7")
	//query := bleve.NewTermQuery(dir)

	//search := bleve.NewSearchRequest(query)
	//searchResults, err := index.Search(search)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	fmt.Println(searchResult)
	//for i := 0; i < 125; i++ {
	searchResult.Hits[0].ID = searchRes
	//}
	return searchRes
}
