package logenc

import (
	"fmt"

	"github.com/oklog/ulid/v2"
)

func MergeLines(ch1 chan LogList, ch2 chan LogList) chan LogList {
	res := make(chan LogList)
	//var wg sync.WaitGroup

	go func() {
		for {
			var line1 LogList
			var line2 LogList
			var ulid1 ulid.ULID
			var ulid2 ulid.ULID
			var ok1, ok2 bool

			line1, ok1 = <-ch1
			if len(line1.XML_RECORD_ROOT) != 0 && line1.XML_RECORD_ROOT[0].XML_ULID != "00000000000000000000000000" {

				ulid1, _ = ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)
				//fmt.Println("check lin1", line1.XML_RECORD_ROOT[0].XML_ULID)
			}
			//wg.Wait()

			line2, ok2 = <-ch2
			if len(line2.XML_RECORD_ROOT) != 0 && line2.XML_RECORD_ROOT[0].XML_ULID != "00000000000000000000000000" {
				ulid2, _ = ulid.ParseStrict(line2.XML_RECORD_ROOT[0].XML_ULID)
				//fmt.Println("check lin2", line2.XML_RECORD_ROOT[0].XML_ULID)
			}
			//wg.Wait()

			//wg.Wait()
			//fmt.Println("start5")
			if !ok1 && !ok2 {
				close(res)
				return

			}

			//fmt.Println(ulid1.Compare(ulid2))
			//wg.Add(1)
			//defer wg.Done()
			if ulid1.Compare(ulid2) == 1 && ulid1.String() != "00000000000000000000000000" && ulid2.String() != "00000000000000000000000000" {
				res <- line2
				fmt.Println(res, "res")

				//fmt.Println(line2, "line2")
				//fmt.Println("ulid2", ulid2)
				//fmt.Println("ulid2SSSS", ulid2.String())
				//fmt.Println("ulid1", ulid1)
				//fmt.Println("ulid1SSSS", ulid1.String())
				//fmt.Println("line2", line2)

				//new.Write([]byte(scanner2.Text()))

			} else if ulid1.String() != "00000000000000000000000000" && ulid2.String() != "00000000000000000000000000" {
				res <- line1
				fmt.Println(res, "res,0,-1")
				//fmt.Println(line1, "line1")

				//fmt.Println("ulid2", ulid2)
				//fmt.Println("ulid2SSSS", ulid2.String())
				//fmt.Println("ulid1", ulid1)
				//fmt.Println("ulid1SSSS", ulid1.String())
				//fmt.Println(ulid1)
				//new.Write([]byte(scanner1.Text()))

			}
		}

	}()

	return res
}
