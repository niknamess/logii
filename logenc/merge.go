package logenc

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"

	"github.com/oklog/ulid/v2"
)

var dlog bool = !!false

func MergeLines(ch1 chan LogList, ch2 chan LogList) chan LogList {
	res := make(chan LogList)
	var nullULID string = "00000000000000000000000000"

	var savedUlid ulid.ULID

	writeRes := func(line LogList, uu ulid.ULID) {
		if uu.Compare(savedUlid) < 1 {
			if dlog {
				fmt.Println("   !write:", savedUlid, "  ", uu)
			}
			return
		}
		savedUlid = uu
		res <- line
		if dlog {
			fmt.Println("    write:", uu)
		}
	}

	go func() {
		entropy := rand.New(rand.NewSource(1))
		minUlid := ulid.MustNew(0, entropy)
		emptyUlid, _ := ulid.ParseStrict("")
		var ulid1 ulid.ULID
		var ulid2 ulid.ULID
		var line1 LogList
		var line2 LogList
		var ok1, ok2 bool
		for {

			if ulid1 == emptyUlid {
				line1, ok1 = <-ch1
				if ok1 && len(line1.XML_RECORD_ROOT) != 0 && line1.XML_RECORD_ROOT[0].XML_ULID != nullULID {
					if dlog {
						fmt.Println("ulid1 read", line1)
					}
					ulid1, _ = ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)
				}
			}
			if ulid2 == emptyUlid {
				line2, ok2 = <-ch2
				if ok2 && len(line2.XML_RECORD_ROOT) != 0 && line2.XML_RECORD_ROOT[0].XML_ULID != nullULID {
					if dlog {
						fmt.Println("ulid2 read", line2)
					}
					ulid2, _ = ulid.ParseStrict(line2.XML_RECORD_ROOT[0].XML_ULID)
				}
			}

			// если входные данные кончились, то закрываем выходной канал.
			if !ok1 && !ok2 {
				if dlog {
					fmt.Println("stop")
				}
				close(res)
				return
			}

			// отдельно обрабатываем случай когда один из входных каналов закрыт или выдает невалидные данные
			bestUlid := emptyUlid
			var bestLine LogList

			if ulid1.Compare(minUlid) < 1 {
				ulid1 = emptyUlid
				bestUlid = ulid2
				bestLine = line2
			}

			if ulid2.Compare(minUlid) < 1 {
				ulid2 = emptyUlid
				bestUlid = ulid1
				bestLine = line1
				if bestUlid.Compare(minUlid) < 1 {
					// в случае если нет ни одного ULID
					if dlog {
						fmt.Println("  check: no one")
					}
					continue
				}
			}

			if bestUlid.Compare(minUlid) > 0 {
				if dlog {
					fmt.Println("  check: only one", bestLine)
				}
				writeRes(bestLine, bestUlid)

				ulid1 = emptyUlid
				ulid2 = emptyUlid
				continue
			}

			// сравниваем гарантированно валидные ulid
			if ulid1.Compare(ulid2) == 1 {
				if dlog {
					fmt.Println("  check: ulid1>ulid2", ulid2, " < ", ulid1)
				}
				writeRes(line2, ulid2)
				ulid2 = emptyUlid
			} else if ulid1.Compare(ulid2) == -1 {
				if dlog {
					fmt.Println("  check: ulid2>ulid1", ulid1, " < ", ulid2)
				}
				writeRes(line1, ulid1)
				ulid1 = emptyUlid
			} else {
				if dlog {
					fmt.Println("  check: ulid1=ulid2", ulid1, " = ", ulid2)
				}
				writeRes(line1, ulid1)

				ulid1 = emptyUlid
				ulid2 = emptyUlid
			}
		}
	}()

	return res
}

/*
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
*/
func CreateDir(path string) {
	fileN := filepath.Base(path)
	//Create a folder/directory at a full qualified path
	err := os.MkdirAll("./repdata/"+fileN, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteOldsFiles(path string, labels string) {
	fileN := filepath.Base(path)
	err := os.Remove("./repdata/" + fileN + "/" + fileN + labels)
	if err != nil {
		log.Fatal(err)
	}

}

func RenameFile(path string, label string) {
	fileN := filepath.Base(path)
	Original_Path := "./repdata/" + fileN + "/" + fileN
	New_Path := "./repdata/" + fileN + "/" + fileN + label
	e := os.Rename(Original_Path, New_Path)
	if e != nil {
		log.Fatal(e)
	}
}

func OpenCreateFile(path string, label string, fileOs *os.File) *os.File {
	fileN := filepath.Base(path)
	file, err := os.OpenFile("./repdata/"+fileN+"/"+fileN+label, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return file
}

func CopyFile(path string, label string, fileOs *os.File) *os.File {
	fileN := filepath.Base(path)
	file, err := os.OpenFile("./repdata/"+fileN+"/"+fileN+label, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	bytesWritten, err := io.Copy(file, fileOs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bytes Written: %d\n", bytesWritten)
	return file
}

func Merge(path string) {
	fileN := filepath.Base(path)
	var wg sync.WaitGroup

	ch1 := make(chan LogList, 100)
	ch2 := make(chan LogList, 100)
	original, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	if CheckFileSum(path, "rep") == true {
		RenameFile(path, "old")
		CopyFile(path, "new", original)
		OpenCreateFile(path, "old", original)
		fileNew, err := os.OpenFile("./repdata/"+fileN+"/"+fileN, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

		if err != nil {
			log.Fatal(err)
		}
		FC, _ := os.Open("./repdata/" + fileN + "/" + fileN + "new")
		defer FC.Close()
		FN, _ := os.Open("./repdata/" + fileN + "/" + fileN + "old")
		defer FN.Close()
		scanner1 := bufio.NewScanner(FN)
		scanner2 := bufio.NewScanner(FC)

		wg.Add(1)
		go func() {
			defer wg.Done()
			for scanner1.Scan() {
				str1 := ProcLineDecodeXML(scanner1.Text())
				if len(str1.XML_RECORD_ROOT) != 0 {
					ch1 <- str1
				}
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			for scanner2.Scan() {
				str2 := ProcLineDecodeXML(scanner2.Text())
				if len(str2.XML_RECORD_ROOT) != 0 {
					ch2 <- str2
				}
			}
		}()

		wg.Wait()
		close(ch1)
		close(ch2)
		//err = os.Remove("test" + fileN)
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//	}

		f, _ := os.Create("test" + fileN)
		ch3 := MergeLines(ch1, ch2)
		for val := range ch3 {

			if len(val.XML_RECORD_ROOT) != 0 {

				xmlline := EncodeXML(val)
				line := EncodeLine(xmlline)
				f.WriteString(line + "\n")
				fileNew.WriteString(line + "\n")
			}
		}
		f.Close()
		fileNew.Close()
		DeleteOldsFiles(path, "old")
		DeleteOldsFiles(path, "new")
	}

}
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func Replication(path string) {
	CreateDir("")

	fileN := filepath.Base(path)
	original, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	files, err := ioutil.ReadDir("/home/nik/projects/Course/logi2/logenc/repdata/")
	if err != nil {
		//CreateDir(path)
		log.Fatal(err)
	}
	//defer files.Close()
	//files.Size()
	ok, err := IsDirEmpty("/home/nik/projects/Course/logi2/logenc/repdata/")
	if err != nil {
		fmt.Println(err)

	}
	if ok == true {
		CreateDir(path)
		CopyFile(path, "", original)
		WriteFileSum(path, "rep")
	}
	for _, f := range files {
		fmt.Println(f.Name())
		if f.Name() == fileN {
			Merge(path)
			WriteFileSum("./repdata/"+fileN+"/"+fileN, "rep")
			//RemoveLine("md5", fileN, "rep")
			return
		}
	}
	if ok == false {
		CreateDir(path)
		CopyFile(path, "", original)
		WriteFileSum(path, "rep")
	}

}
