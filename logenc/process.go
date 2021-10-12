package logenc

import (
	//	"fmt"

	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	//	"sync"
	//	"sync/atomic"

	//"github.com/blevesearch/bleve"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Logger *log.Logger
	mu     sync.Mutex
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
	//fmt.Print(csvline)
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
	//log.Println("1")
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

func ProcLineDX(line string) (val LogList) {

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

//func ProcMapFile(file string) map[string]string {
func ProcMapFile(file string) {
	ch := make(chan string, 100)
	//log.Println("1")
	SearchMap := make(map[string]string)
	var wg sync.WaitGroup
	//var counter int32 = 0
	var data LogList
	var datas string

	//fmt.Println("run 1")

	go func() {
		//wg.Add(1)
		//defer wg.Done()
		for {
			select {
			case line, ok := <-ch:
				if !ok {
					break
				}
				go func(line string) {
					wg.Add(1)
					defer wg.Done()

					//fmt.Println("run3")
					data = ProcLineDX(line)
					datas = ProcLine(line)
					//fmt.Println("stop")
					//atomic.AddInt32(&counter, 1)

					if len(data.XML_RECORD_ROOT) > 0 {
						mu.Lock()
						SearchMap[data.XML_RECORD_ROOT[0].XML_ULID] = datas
						mu.Unlock()
					}
				}(line)
				//close(ch)
			}
		}
	}()
	//close(ch)

	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}

	//close(ch)
	wg.Wait()
	close(ch)

}

func ProcMapFilePP(file string) {
	ch := make(chan string, 1000000)
	//log.Println("1")
	SearchMap := make(map[string]string)
	//var wg sync.WaitGroup
	//var counter int32 = 0
	var data LogList
	var datas string

	//fmt.Println("run 1")

	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
	fmt.Println("run")
	for {
		select {
		case line, ok := <-ch:
			if !ok {
				break
			}

			//fmt.Println("run3")
			data = ProcLineDX(line)
			datas = ProcLine(line)
			//fmt.Println("stop")
			//atomic.AddInt32(&counter, 1)

			if len(data.XML_RECORD_ROOT) > 0 {
				//mu.Lock()
				SearchMap[data.XML_RECORD_ROOT[0].XML_ULID] = datas
				//mu.Unlock()
			}

			//close(ch)
		}
		if len(ch) == 0 {
			break
		}

	}
	//b, _ := json.MarshalIndent(SearchMap, "", "  ")
	//fmt.Print(string(b))
	//close(ch)

	//close(ch)

}
