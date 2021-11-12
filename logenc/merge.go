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

func MergeLines(ch1 chan LogList, ch2 chan LogList) chan LogList {
	res := make(chan LogList)
	var nullULID string = "00000000000000000000000000"
	go func() {
		entropy := rand.New(rand.NewSource(1))
		minUlid := ulid.MustNew(0, entropy)
		emptyUlid, _ := ulid.ParseStrict("")
		var ulid1 ulid.ULID
		var ulid2 ulid.ULID
		var line1 LogList
		var line2 LogList
		for {
			var ok1, ok2 bool
			if ulid1 == emptyUlid {
				line1, ok1 = <-ch1
				if ok1 && len(line1.XML_RECORD_ROOT) != 0 && line1.XML_RECORD_ROOT[0].XML_ULID != nullULID {
					fmt.Println("line1 st", line1)
					ulid1, _ = ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)
				}
			}
			if ulid2 == emptyUlid {
				line2, ok2 = <-ch2
				if ok2 && len(line2.XML_RECORD_ROOT) != 0 && line2.XML_RECORD_ROOT[0].XML_ULID != nullULID {
					fmt.Println("line2 st", line2)
					ulid2, _ = ulid.ParseStrict(line2.XML_RECORD_ROOT[0].XML_ULID)
				}
			}
			if !ok1 && !ok2 {
				fmt.Println("stop")
				close(res)
				return
			}
			// отдельно обрабатываем случай когда один из входных каналов закрыт или выдает невалидные данные
			bestUlid := emptyUlid
			var bestLine LogList

			if ulid1.Compare(minUlid) < 1 {
				bestUlid = ulid2
				bestLine = line2
			} else if ulid2.Compare(minUlid) < 1 {
				bestUlid = ulid1
				bestLine = line1

			}
			if bestUlid.Compare(minUlid) > 0 {
				res <- bestLine
				fmt.Println("best", bestLine)
				continue
			}
			// сравниваем гарантированно валидные ulid
			if ulid1.Compare(ulid2) == 1 {
				res <- line2
				fmt.Println("2", line2)
				ulid2 = emptyUlid
			} else if ulid1.Compare(ulid2) == -1 {
				res <- line1
				fmt.Println("1", line1)
				ulid1 = emptyUlid
			} else {
				res <- line1
				fmt.Println("1", line1)

				ulid1 = emptyUlid
				ulid2 = emptyUlid
			}
		}
	}()

	return res
}

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

func Checkmd5rep(path string) {
	fileN := filepath.Base(path)
	original, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	files, err := ioutil.ReadDir("./repdata")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.Name() == fileN {
			return
		} else {
			CreateDir(path)
			CopyFile(path, "", original)
			WriteFileSum(path, "rep")
		}
	}

	//OpenCreateFile(path, "", original)
	//WriteFileSum(path, "rep")

}

func Replication(path string) {

	//fileN := filepath.Base(path)
	//CreateDir(path)
	var ch1 chan LogList
	var ch2 chan LogList
	var ch3 chan LogList

	original, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	//CheckSum file (md5rep)
	if CheckFileSum(path, "rep") == true {
		//Rename old file
		RenameFile(path, "old")

		fileChanged := CopyFile(path, "new", original)
		defer fileChanged.Close()
		fileOld := OpenCreateFile(path, "old", original)
		defer fileOld.Close()
		fileNew := OpenCreateFile(path, "", original)
		defer fileNew.Close()
		//
		scanner1 := bufio.NewScanner(fileOld)
		scanner2 := bufio.NewScanner(fileChanged)
		//Merge and create new file
		for scanner1.Scan() || scanner2.Scan() {
			str1 := ProcLineDecodeXML(scanner1.Text())
			str2 := ProcLineDecodeXML(scanner2.Text())
			ch1 <- str1
			ch2 <- str2
			ch3 = MergeLines(ch1, ch2)
		}
		for val := range ch3 {

			if len(val.XML_RECORD_ROOT) != 0 {
				//got++
				//fmt.Println(val.XML_RECORD_ROOT[0].XML_ULID)
				fileNew.WriteString(val.XML_RECORD_ROOT[0].XML_ULID)
			}
		}
		//Delete two old file
		DeleteOldsFiles(path, "old")
		DeleteOldsFiles(path, "new")

	}

}

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
