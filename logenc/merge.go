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

	"github.com/oklog/ulid/v2"
)

var dlog bool = false

func MergeLines(ch1 chan LogList, ch2 chan LogList) chan LogList {
	res := make(chan LogList)
	var nullULID string = "00000000000000000000000000"
	var count int = 0

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
		count++
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
					_, err := ulid.Parse(line1.XML_RECORD_ROOT[0].XML_ULID)
					if err == nil {

						ulid1, _ = ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)
					} //else {
					//	res <- line1
					//}

				}
			}
			if ulid2 == emptyUlid {
				line2, ok2 = <-ch2
				if ok2 && len(line2.XML_RECORD_ROOT) != 0 && line2.XML_RECORD_ROOT[0].XML_ULID != nullULID {
					if dlog {
						fmt.Println("ulid2 read", line2)
					}
					_, err := ulid.Parse(line2.XML_RECORD_ROOT[0].XML_ULID)
					if err == nil {
						ulid2, _ = ulid.ParseStrict(line2.XML_RECORD_ROOT[0].XML_ULID)
					} //else {
					//	res <- line2
					//}
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
	//fmt.Println("count", count)
	return res
}

func CreateDir(dirpath string, path string) {
	fileN := filepath.Base(path)
	//Create a folder/directory at a full qualified path
	err := os.MkdirAll(dirpath+fileN, os.ModePerm)
	if err != nil {
		log.Println("CreateDir:", err)
	}
}

func DeleteOldsFiles(dirpath string, path string, labels string) {
	fileN := filepath.Base(path)
	log.Println("RemoveOldfile:", dirpath+fileN+labels)
	//err := os.Remove(dirpath + "/" + fileN + labels)
	err := os.Remove(dirpath + fileN + labels)
	if err != nil {
		log.Println("DeleteOldsFiles err:", err)
	}

}

func RenameFile(dirpath string, path string, label string) {
	fileN := filepath.Base(path)
	Original_Path := dirpath + fileN
	New_Path := dirpath + fileN + label
	e := os.Rename(Original_Path, New_Path)
	if e != nil {
		log.Println("RenameFile:", e)
	}
}

func OpenCreateFile(dirpath string, path string, label string) {
	fileN := filepath.Base(path)
	file, err := os.OpenFile(dirpath+fileN+label, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		//log.Fatal(err)
		log.Println(err)
	}
	file.Close()
}

func CopyFile(dirpath string, path string, label string, fileOs *os.File) {
	fileN := filepath.Base(path)
	file, err := os.OpenFile(dirpath+fileN+label, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("CopyFile error", err)
		return
	}
	defer file.Close()

	bytesWritten, err := io.Copy(file, fileOs)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Bytes Written: %d\n", bytesWritten)
	}
}

func Merge(dirpath string, path string) {
	fileN := filepath.Base(path)
	//var wg sync.WaitGroup

	ch1 := make(chan LogList, 100)
	ch2 := make(chan LogList, 100)
	original, err := os.Open(path)
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
	}
	defer original.Close()
	if !CheckFileSum(path, "rep", "") {
		return
	} else {
		RenameFile(dirpath, path, "old")
		CopyFile(dirpath, path, "new", original)
		OpenCreateFile(dirpath, path, "old")
		fileNew, err := os.OpenFile(dirpath+fileN, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

		if err != nil {
			log.Fatal("Open in Merge", err)
		}
		FC, _ := os.Open(dirpath + fileN + "new")
		defer FC.Close()
		FN, _ := os.Open(dirpath + fileN + "old")
		defer FN.Close()
		scanner1 := bufio.NewScanner(FN)
		scanner2 := bufio.NewScanner(FC)

		go func() {

			for scanner1.Scan() {
				str1 := ProcLineDecodeXML(scanner1.Text())
				if len(str1.XML_RECORD_ROOT) != 0 {
					ch1 <- str1
				}
			}
			close(ch1)

		}()

		go func() {

			for scanner2.Scan() {
				str2 := ProcLineDecodeXML(scanner2.Text())
				if len(str2.XML_RECORD_ROOT) != 0 {
					ch2 <- str2
				}
			}
			close(ch2)

		}()

		//f, _ := os.Create("test" + fileN)
		ch3 := MergeLines(ch1, ch2)
		for val := range ch3 {

			if len(val.XML_RECORD_ROOT) != 0 {

				xmlline := EncodeXML(val)
				line := EncodeLine(xmlline)
				//f.WriteString(line + "\n")
				fileNew.WriteString(line + "\n")
			}
		}

		//f.Close()
		fileNew.Close()
		DeleteOldsFiles(dirpath, path, "old")
		DeleteOldsFiles(dirpath, path, "new")
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
	var dirpath string = "./repdata/"
	CreateDir(dirpath, "")

	fileN := filepath.Base(path)
	original, err := os.Open(path)
	if err != nil {
		log.Fatal("Replication OpenFile", err)
	}
	defer original.Close()

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {

		log.Fatal("ReadDir", err)
	}

	ok, err := IsDirEmpty(dirpath)
	if err != nil {
		fmt.Println("DirEmpty", err)

	}
	if ok {
		//CreateDir(path)
		CopyFile(dirpath, path, "", original)
		WriteFileSum(dirpath+fileN, "rep", "")
	} else {
		for _, f := range files {
			//fmt.Println(f.Name())
			if f.Name() == fileN {
				Merge(dirpath, path)
				WriteFileSum(dirpath+fileN, "rep", "")

				return
			}
		}
	}
	if !ok {
		//CreateDir(path)
		CopyFile(dirpath, path, "", original)
		WriteFileSum(path, "rep", "")
	}

}
