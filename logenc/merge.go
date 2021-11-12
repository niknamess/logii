package logenc

import (
	"fmt"
	"math/rand"

	"github.com/oklog/ulid/v2"
)

func MergeLines(ch1 chan LogList, ch2 chan LogList) chan LogList {
	res := make(chan LogList)
	//var wg sync.WaitGroup

	go func() {
		//saveulid1:= emptyUlid
		entropy := rand.New(rand.NewSource(1))
		minUlid := ulid.MustNew(0, entropy)
		emptyUlid, _ := ulid.ParseStrict("")
		//saveulid1 := emptyUlid
		//saveulid2 := emptyUlid
		var ulid1 ulid.ULID
		var ulid2 ulid.ULID

		for {
			var line1 LogList
			var line2 LogList
			//var ulid1 ulid.ULID
			//var ulid2 ulid.ULID
			var ok1, ok2 bool

			//ulid1 = emptyUlid
			if ulid1 == emptyUlid {

				line1, ok1 = <-ch1
				if ok1 && len(line1.XML_RECORD_ROOT) != 0 && line1.XML_RECORD_ROOT[0].XML_ULID != "00000000000000000000000000" {

					ulid1, _ = ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)
					//fmt.Println("check lin1", line1.XML_RECORD_ROOT[0].XML_ULID)

				}
			}
			if ulid2 == emptyUlid {
				//ulid2 = emptyUlid
				line2, ok2 = <-ch2
				if ok2 && len(line2.XML_RECORD_ROOT) != 0 && line2.XML_RECORD_ROOT[0].XML_ULID != "00000000000000000000000000" {
					ulid2, _ = ulid.ParseStrict(line2.XML_RECORD_ROOT[0].XML_ULID)
					//fmt.Println("check lin2", line2.XML_RECORD_ROOT[0].XML_ULID)

				}
				//bestline
			}
			fmt.Println("ulid1", ulid1)
			fmt.Println("ulid2", ulid2)
			fmt.Println("check")
			//fmt.Println("start5")
			
			if !ok1 && !ok2 {
				fmt.Println("Stop")
				fmt.Println("ulid1", ulid1)
				fmt.Println("ulid2", ulid2)
				if ulid1.String() != "00000000000000000000000000" {
					res <- line1
					fmt.Println("1", line1)
				} else if ulid2.String() != "00000000000000000000000000" {
					res <- line2
					fmt.Println("2", line2)
				}

				fmt.Println("stop")
				close(res)
				return
			}

			//fmt.Println(ulid1.Compare(ulid2))

			// отдельно обрабатываем случай когда один из входных каналов закрыт или выдает невалидные данные
			bestUlid := emptyUlid
			var bestLine LogList

			if ulid1.Compare(minUlid) < 1 {
				bestUlid = ulid2
				bestLine = line2
				//saveulid1 = ulid2

			} else if ulid2.Compare(minUlid) < 1 {
				bestUlid = ulid1
				bestLine = line1
				//saveulid2 = ulid1

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
				//saveulid1 = ulid1
				ulid2 = emptyUlid

				//fmt.Println(line2, "line2")
				//fmt.Println("ulid2", ulid2)
				//fmt.Println("ulid2SSSS", ulid2.String())
				//fmt.Println("ulid1", ulid1)
				//fmt.Println("ulid1SSSS", ulid1.String())
				//fmt.Println("line2", line2)

				//new.Write([]byte(scanner2.Text()))

			} else if ulid1.Compare(ulid2) == -1 {
				res <- line1
				fmt.Println("1", line1)

				//saveulid2 = ulid2
				ulid1 = emptyUlid
				//fmt.Println(line1, "line1")

				//fmt.Println("ulid2", ulid2)
				//fmt.Println("ulid2SSSS", ulid2.String())
				//fmt.Println("ulid1", ulid1)
				//fmt.Println("ulid1SSSS", ulid1.String())
				//fmt.Println(ulid1)
				//new.Write([]byte(scanner1.Text()))

			} else {
				res <- line1
				fmt.Println("1", line1)

				ulid1 = emptyUlid
				ulid2 = emptyUlid
			}
			//	if !ok1 && !ok2 {
			//	close(res)
			//return
			//}
		}

	}()

	return res
}
