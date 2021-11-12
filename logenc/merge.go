package logenc

import (
	"fmt"
	"math/rand"

	"github.com/oklog/ulid/v2"
)

func MergeLines(ch1 chan LogList, ch2 chan LogList) chan LogList {
	res := make(chan LogList)
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
				if ok1 && len(line1.XML_RECORD_ROOT) != 0 && line1.XML_RECORD_ROOT[0].XML_ULID != "00000000000000000000000000" {
					fmt.Println("line1 st", line1)
					ulid1, _ = ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)
				}
			}
			if ulid2 == emptyUlid {
				line2, ok2 = <-ch2
				if ok2 && len(line2.XML_RECORD_ROOT) != 0 && line2.XML_RECORD_ROOT[0].XML_ULID != "00000000000000000000000000" {
					fmt.Println("line2 st", line2)
					ulid2, _ = ulid.ParseStrict(line2.XML_RECORD_ROOT[0].XML_ULID)
				}
			}
			//fmt.Println("ulid1", ulid1)
			//fmt.Println("ulid2", ulid2)
			//fmt.Println("check")
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
