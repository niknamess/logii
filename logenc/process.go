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

//slowely
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

func WriteFileSum(file string) bool {

	//checksum := MD5(file)
	checksum2 := FileMD5(file)
	fileN := filepath.Base(file)
	metaname := "hashmd5"
	ind := false
	//fmt.Printf("Checksum 1: %s\n", checksum)
	fmt.Printf("current Checksum: %s\n", checksum2)
	f, _ := os.OpenFile(metaname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	//f.Write([]byte(fileN + " " + checksum2 + "\n"))
	scanner := bufio.NewScanner(f)
	_, ok := os.Stat(metaname)
	//fmt.Println(info)
	fmt.Println(ok)
	if ok != nil {
		f.Write([]byte(checksum2 + " " + fileN + "\n"))
	}
	line := 1

	for scanner.Scan() {
		//lineStr := scanner.Text()
		if strings.Contains(scanner.Text(), (checksum2 + " " + fileN)) {
			ind = false
			fmt.Println(ind)
			return ind
		} else {
			f.Write([]byte(checksum2 + " " + fileN + "\n"))
			ind = true
			fmt.Println(ind)
		}

		line++

		//mu.Unlock()

	}
	return ind
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
