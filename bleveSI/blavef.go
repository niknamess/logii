package bleveSI

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/blevesearch/bleve"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

var (
	Logger       *log.Logger
	mu           sync.Mutex
	sliceLoglist []logenc.LogList
	//lines []string
)

//classic
func ProcFileBreve(fileN string, file string) {
	//func ProcFileBreve(file string) {
	var wg sync.WaitGroup
	//var counter int32 = 0
	var data logenc.LogList
	//metaname := "example.bleve"
	if len(file) <= 0 {
		return
	}

	dir := "./blevestorage/"
	extension := ".bleve"
	metaname := dir + fileN + extension
	if logenc.CheckFileSum(file) == false {
		return
	}

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
					data = logenc.ProcLineDecodeXML(line)
					//fmt.Println((data.XML_RECORD_ROOT))
					//fmt.Println(len(data.XML_RECORD_ROOT))
					//atomic.AddInt32(&counter, 1)
					//println(counter)
					if len(data.XML_RECORD_ROOT) > 0 {
						index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
					}
				}
			}

		}()

	}

	err = logenc.ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
	close(ch)
	wg.Wait()
	index.Close()
	logenc.WriteFileSum(file)
}

func ProcFileBleveSPEED(fileN string, file string) {
	var data logenc.LogList

	if len(file) <= 0 {
		return
	}

	dir := "./blevestorage/"
	extension := ".bleve"
	metaname := dir + fileN + extension
	if logenc.CheckFileSum(file) == false {
		return
	}
	index, err := bleve.Open(metaname)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(metaname, mapping)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	ch := make(chan string, 100)
	for i := 5; i > 0; i-- {
		//brloop:
		go func() {
			//wg.Add(1)
			//println("Start") //wg.Done()

			//brloop:
			for {
				select {
				case line, ok := <-ch:
					if !ok {
						break //brloop
					}
					//for i := 5; i > 0; i-- {
					go func(line string) {
						//wg.Add(1)
						//defer wg.Done()
						//BUFFER
						data = logenc.ProcLineDecodeXML(line)
						sliceLoglist = append(sliceLoglist, data)
						if len(sliceLoglist) == 100 {
							for i := 0; i < len(sliceLoglist); i++ {
								if len(data.XML_RECORD_ROOT) > 0 {
									index.Index(data.XML_RECORD_ROOT[0].XML_ULID, sliceLoglist[i])
									//sliceLoglist = nil
								}
							}
							sliceLoglist = nil
						}
						//buf := buffer.NewSpill(1024, data)
						//	b, _ := json.Marshal(data)
						//	buf := buffer.NewSpill(1024, b)

						//BUFFER

						//	atomic.AddInt32(&counter, 1)
						//	defer println(counter)
						if len(data.XML_RECORD_ROOT) > 0 {
							index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
						}
						//defer printlm
					}(line)
					//}
				}
			}
			//println("Stop") //wg.Done()
		}()
	}
	err = logenc.ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
	close(ch)
	index.Close()
	logenc.WriteFileSum(file)
}

func ProcBleveSearch(fileN string, word string) []string {
	dir := "./blevestorage/"
	extension := ".bleve"
	filename := fileN
	metaname := dir + filename + extension
	index, _ := bleve.Open(metaname)
	//query := bleve.NewFuzzyQuery(dir)
	query := bleve.NewMatchQuery(word)
	query.Fuzziness = 1
	mq := bleve.NewMatchPhraseQuery(word)
	rq := bleve.NewRegexpQuery(word)
	q := bleve.NewDisjunctionQuery(query, mq, rq)

	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Size = 1000000000000000000

	searchResult, _ := index.Search(searchRequest)
	searchRequest.Fields = []string{"XML_RECORD_ROOT"}

	docs := make([]string, 0)

	for _, val := range searchResult.Hits {
		id := val.ID
		docs = append(docs, id)
	}

	index.Close()
	return docs

}

func ProcFileBreveSLOWLY(fileName string, file string) {
	const pieces int = 10

	//func ProcFileBreve(file string) {
	var wg sync.WaitGroup
	//var counter int32 = 0
	var lines []string
	//var data logenc.LogList

	if len(file) <= 0 {
		return
	}

	dir := "./blevestorage/"
	extension := ".bleve"
	metaname := dir + fileName + extension
	if logenc.CheckFileSum(file) == false {
		return
	}

	//metaname := "example.bleve"
	index, err := bleve.Open(metaname)

	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(metaname, mapping)
	}

	defer index.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	fileF, _ := os.Open(file)
	scanner := bufio.NewScanner(fileF)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 4*1024*1024)
	for scanner.Scan() {

		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return
	}

	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}

	var datas [pieces][]logenc.LogList

	curNum := 0

	fmt.Println("lines", len(lines))

	for _, line := range lines {

		datas[curNum] = append(datas[curNum], logenc.ProcLineDecodeXML(line))

		curNum++
		if curNum == pieces {
			curNum = 0
		}

	}

	for _, data := range datas {

		//fmt.Println("data len=", len(data))
		wg.Add(1)
		go func(dataPiece []logenc.LogList) {

			for _, data := range dataPiece {
				if len(data.XML_RECORD_ROOT) > 0 {
					index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
				}
			}
			wg.Done()
		}(data)

	}

	wg.Wait()
	logenc.WriteFileSum(file)

}
