package logenc

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Logger   *log.Logger
	mu       sync.Mutex
	ind      bool
	fileSize int64

	//true untyped bool = true
)
var (
	sliceLoglist []LogList
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

func ProcLineDecodeXML(line string) (val LogList) {

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

func ProcMapFile(file string) map[string]string {
	//func ProcMapFile(file string) {
	if len(file) <= 0 {
		return nil
	}
	ch := make(chan string, 100)
	SearchMap := make(map[string]string)
	var wg sync.WaitGroup
	var data LogList
	var datas string
	go func() {
		for {
			select {
			case line, ok := <-ch:
				if !ok {
					break
				}
				go func(line string) {
					wg.Add(1)
					defer wg.Done()
					data = ProcLineDecodeXML(line)
					datas = ProcLine(line)
					if len(data.XML_RECORD_ROOT) > 0 {
						mu.Lock()
						SearchMap[data.XML_RECORD_ROOT[0].XML_ULID] = datas
						mu.Unlock()
					}
				}(line)
			}
		}
	}()
	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
	wg.Wait()
	close(ch)
	return SearchMap
}

//slowely not used
func ProcMapFileREZERV(file string) {
	if len(file) <= 0 {
		return
	}
	ch := make(chan string, 1000000)
	SearchMap := make(map[string]string)
	var data LogList
	var datas string
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
			data = ProcLineDecodeXML(line)
			datas = ProcLine(line)
			if len(data.XML_RECORD_ROOT) > 0 {
				SearchMap[data.XML_RECORD_ROOT[0].XML_ULID] = datas
			}
		}
		if len(ch) == 0 {
			break
		}

	}
}

func CheckFileSum(file string, typeS string) bool {

	checksum2 := FileMD5(file)
	fileN := filepath.Base(file)
	hashFileName := "md5" + typeS
	f, err := os.OpenFile(hashFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checke(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), (checksum2 + " " + fileN)) {
			ind = false
			return ind
		} else {
			ind = true
			//WriteFileSum(file)
		}
		line++
	}
	scanner = nil
	return ind
}

func WriteFileSum(file string, typeS string) {

	checksum2 := FileMD5(file)
	fileN := filepath.Base(file)
	hashFileName := "md5" + typeS
	fmt.Println(os.Getwd())
	f, err := os.OpenFile(hashFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checke(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {

		if strings.Contains(scanner.Text(), (checksum2 + " " + fileN)) {
			ind = false

			return
		} else {
			ind = true
		}

		line++
	}
	scanner = nil
	if ind == true {

		f.Write([]byte(checksum2 + " " + fileN + "\n"))
	}
	fi, _ := f.Stat()

	if fi.Size() == 0 {
		f.Write([]byte(checksum2 + " " + fileN + "\n"))
	}

}

func checke(e error) {
	if e != nil {
		panic(e)
	}
}

// FileMD5 создает md5-хеш из содержимого нашего файла.
func FileMD5(path string) string {
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(h, f)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
