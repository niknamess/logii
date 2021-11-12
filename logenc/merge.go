package logenc

import (
	"fmt"
	"io"
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

func DeleteOldsFiles(path string) {
	fileN := filepath.Base(path)
	err := os.Remove("./repdata/" + fileN + "/" + fileN + "new")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("./repdata/" + fileN + "/" + fileN + "old")
	if err != nil {
		log.Fatal(err)
	}
}

func RenameFile(path string) {
	fileN := filepath.Base(path)
	Original_Path := "./repdata/" + fileN + "/" + fileN
	New_Path := "./repdata/" + fileN + "/" + fileN + "old"
	e := os.Rename(Original_Path, New_Path)
	if e != nil {
		log.Fatal(e)
	}
}

func OpenCreateFile(path string, label string, fileOs *os.File) {
	fileN := filepath.Base(path)
	old, err := os.OpenFile("./repdata/"+fileN+"/"+fileN+label, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer old.Close()

	bytesWritten, err := io.Copy(old, fileOs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bytes Written: %d\n", bytesWritten)

}

func CopyLogs(path string) {

	//fileN := filepath.Base(path)
	CreateDir(path)
	original, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	//CheckSum file (md5)
	if CheckFileSum(path) == true {
		OpenCreateFile(path, "", original)
	} else {
		OpenCreateFile(path, "new", original)

		//Rename old file
		RenameFile(path)
		//Merge and create new file

		//Delete two old file
		DeleteOldsFiles(path)

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
