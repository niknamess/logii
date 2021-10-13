package bleveSI

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/blevesearch/bleve"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

var (
	Logger *log.Logger
	mu     sync.Mutex
)

//slowly
func ProcFileBreve(fileN string, file string) {
	//func ProcFileBreve(file string) {
	var wg sync.WaitGroup
	//var counter int32 = 0
	var data logenc.LogList
	dir := "./blevestorage/"
	extension := ".bleve"
	filename := fileN
	metaname := dir + filename + extension
	//metaname := "example.bleve"
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
					data = logenc.ProcLineDX(line)
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
}

func ProcFileBreveSPEED(fileN string, file string) {
	var data logenc.LogList
	if len(file) <= 0 {
		return
	}

	dir := "./blevestorage/"
	extension := ".bleve"
	metaname := dir + fileN + extension
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
	//for i := runtime.NumCPU() + 1; i > 0; i-- {
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
				go func(line string) {
					//wg.Add(1)
					//defer wg.Done()
					data = logenc.ProcLineDX(line)
					//	atomic.AddInt32(&counter, 1)
					//	defer println(counter)
					if len(data.XML_RECORD_ROOT) > 0 {
						index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
					}
					//defer printlm
				}(line)
			}
		}
		//println("Stop") //wg.Done()
	}()
	//}
	err = logenc.ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}
	//wg.Wait()
	close(ch)
	//wg.Wait()
	index.Close()
}

//func ProcBleveSearch(dir string) []string {
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
