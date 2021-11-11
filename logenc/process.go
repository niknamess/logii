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

func CheckFileSum(file string) bool {

	checksum2 := FileMD5(file)
	fileN := filepath.Base(file)
	hashFileName := "md5"
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
		}
		line++
	}
	scanner = nil
	return ind
}

func WriteFileSum(file string) {

	checksum2 := FileMD5(file)
	fileN := filepath.Base(file)
	hashFileName := "md5"
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

func Createdir(path string) {
	fileN := filepath.Base(path)
	//Create a folder/directory at a full qualified path
	err := os.Mkdir("./repdata/"+fileN, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
func CopyLogs(path string) {
	fileN := filepath.Base(path)
	Createdir(path)
	original, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	if CheckFileSum(path) == true {
		old, err := os.OpenFile("./repdata/"+fileN+"/"+fileN+"old", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

		if err != nil {
			log.Fatal(err)
		}
		defer old.Close()

		bytesWritten, err := io.Copy(old, original)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Bytes Written: %d\n", bytesWritten)
	} else {
		new, err := os.OpenFile("./repdata/"+fileN+"/"+fileN+"new", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer new.Close()

		bytesWritten, err := io.Copy(new, original)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Bytes Written in new file: %d\n", bytesWritten)

		err = os.Remove("./repdata/" + fileN + "/" + fileN + "new")
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove("./repdata/" + fileN + "/" + fileN + "old")
		if err != nil {
			log.Fatal(err)
		}

	}

}

/*
func Comparefiles(path1 string, path2 string, path string) {

	var data1 LogList
	var data2 LogList
	//var wg sync.WaitGroup
	var str1 string
	var str2 string

	fileN := path
	var wait int32 = 0
	new, err := os.OpenFile("/home/nik/projects/Course/logi2/repdata/Test/"+fileN, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//	ch := make(chan string, 100)

		for i := runtime.NumCPU() + 1; i > 0; i-- {
			go func() {
				wg.Add(1)
				defer wg.Done()
				for {
					if wait == 0 || wait == -1 {
						line1 = <-ch
					}

					if wait == 0 || wait == 1 {
						line2 = <-ch
					}

	for {
		err1 := ReadLines(path1, func(line1 string) {
			if wait == 0 || wait == -1 {
				str1 = line1
			}
		})
		if err1 != nil {
			log.Fatalf("ReadLines: %s", err1)
		}

		err2 := ReadLines(path2, func(line2 string) {
			if wait == 0 || wait == 1 {
				str2 = line2
			}
		})
		if err2 != nil {
			log.Fatalf("ReadLines: %s", err2)
		}

		data1 = ProcLineDecodeXML(str1)
		ulid1, _ := ulid.ParseStrict(data1.XML_RECORD_ROOT[0].XML_ULID)
		data2 = ProcLineDecodeXML(str2)
		ulid2, _ := ulid.ParseStrict(data2.XML_RECORD_ROOT[0].XML_ULID)

		if ulid1.Compare(ulid2) == 1 {
			new.Write([]byte(str2))
			//atomic.AddInt32(&counter, 1)
			atomic.StoreInt32(&wait, 1)

		} else if ulid1.Compare(ulid2) == -1 {
			new.Write([]byte(str1))
			atomic.StoreInt32(&wait, -1)

		} else {
			new.Write([]byte(str1))
		}

	}

	//}()
	//	}
	/*
		err1 := ReadLines(path1, func(line1 string) {
			ch <- line1
		})
		err2 := ReadLines(path2, func(line2 string) {
			ch <- line2
		})
		if err1 != nil {
			log.Fatalf("ReadLines: %s", err1)
		}
		if err2 != nil {
			log.Fatalf("ReadLines: %s", err2)
		}
		close(ch)
		wg.Wait()

*/
/*
func Comparefiles2(path1 string, path2 string, path string) {

	//var str1 LogList
	//	var str2 LogList
	ch1 := make(chan string, 100)

	ch2 := make(chan string, 100)

	//fileN := path
	//var wait int32 = 0
	new, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer new.Close()
	file1, err := os.OpenFile(path1, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()

	file2, err := os.OpenFile(path2, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)

	//info1, err := os.Stat(path1)
	//info2, err := os.Stat(path2)

	go func() {
		for {
			select {
			case line1, _ := <-ch1:
				ulid1, _ := ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)

			case line2, _ := <-ch2:
				ulid2, _ := ulid.ParseStrict(line2.XML_RECORD_ROOT[0].XML_ULID)

			}
		}

	}()

	for scanner1.Scan() {

		ch1 <- ProcLine(scanner1.Text())

	}
	for scanner1.Scan() {
		ch2 <- ProcLine(scanner2.Text())

	}
	/*
		for scanner1.Scan() || scanner2.Scan() {
			str1 = ProcLineDecodeXML(scanner1.Text())

			str2 = ProcLineDecodeXML(scanner2.Text())

			ulid1, _ := ulid.ParseStrict(str1.XML_RECORD_ROOT[0].XML_ULID)
			ulid2, _ := ulid.ParseStrict(str2.XML_RECORD_ROOT[0].XML_ULID)

			if ulid1.Compare(ulid2) == 1 {
				new.Write([]byte(scanner2.Text()))

			} else {
				new.Write([]byte(scanner1.Text()))

			}
		}


}
*/
