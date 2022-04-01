package bleveSI

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/oklog/ulid/v2"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

var (
	Logger *log.Logger
)

func bleveIndex(fileN string) (bleve.Index, error) {

	dir := "./blevestorage/"
	extension := ".bleve"
	metaname := dir + fileN + extension
	index, err := bleve.Open(metaname)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.NewUsing(metaname, mapping, scorch.Name, scorch.Name, nil)
	}

	return index, err
}
func ProcBleve(fileN string, file string) {
	var count int = 0
	if !logenc.CheckFileSum(file, "", "") {
		return
	}
	var wg sync.WaitGroup
	index, err := bleveIndex(fileN)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data logenc.LogList
	ch := make(chan string, 100)

	for i := runtime.NumCPU() + 1; i > 0; i-- {
		wg.Add(1)
		go func() {
			//wg.Add(1)
			defer wg.Done()
			batch := index.NewBatch()
		brloop:
			for {

				select {
				case line, ok := <-ch:
					if !ok {
						break brloop
					}
					if count == 1000 {
						err = index.Batch(batch)
						if err != nil {
							fmt.Println("index.Batch(batch) err: ", err)
						}
						count = 0
						batch = index.NewBatch()
					}

					data = logenc.ProcLineDecodeXML(line)
					if len(data.XML_RECORD_ROOT) > 0 {
						batch.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
						count++
						//if count == 100 {
						//	fmt.Println(count)
						//	}
					}

				}
			}
			err = index.Batch(batch)
			if err != nil {
				fmt.Println("index.Batch(batch) err: ", err)
			}

		}()
	}
	err = logenc.ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		fmt.Println("ReadLines: ", err)
		close(ch)
		return
	}
	close(ch)
	wg.Wait()
	index.Close()
	logenc.WriteFileSum(file, "", "")
}

func ProcBleveSearchv2(fileN string, word string) []string {
	//var query *query.MatchQuery
	dir := "./blevestorage/"
	extension := ".bleve"
	filename := fileN
	metaname := dir + filename + extension
	index, _ := bleve.OpenUsing(metaname, nil)

	query := bleve.NewMatchQuery(word)
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
	//sort
	for i := len(docs); i > 0; i-- {
		for j := 1; j < i; j++ {
			j2, _ := ulid.Parse(docs[j-1])
			j1, _ := ulid.Parse(docs[j])
			if j2.Compare(j1) == 1 {
				intermediate := docs[j]
				docs[j] = docs[j-1]
				docs[j-1] = intermediate
			}

		}
	}

	index.Close()
	return docs

}
