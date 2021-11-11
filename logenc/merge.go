package logenc

import (
	"github.com/oklog/ulid/v2"
)

func MergeLines(ch1 chan LogList, ch2 chan LogList) chan LogList {
	res := make(chan LogList)

	go func() {
		for {
			var line1 LogList
			var line2 LogList
			var ulid1 ulid.ULID
			var ulid2 ulid.ULID
			var ok1, ok2 bool

			line1, ok1 = <-ch1
			if len(line1.XML_RECORD_ROOT) != 0 && line1.XML_RECORD_ROOT[0].XML_ULID != "" {

				ulid1, _ = ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)
				//fmt.Println("check lin1", line1.XML_RECORD_ROOT[0].XML_ULID)
			}
			//wg.Wait()

			line2, ok2 = <-ch2
			if len(line2.XML_RECORD_ROOT) != 0 && line2.XML_RECORD_ROOT[0].XML_ULID != "" {
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
			//defer wg.Done()
			if ulid1.Compare(ulid2) == 1 {
				res <- line2
				//fmt.Println(line2, "line2")
				//fmt.Println("ulid2", ulid2)
				//fmt.Println("ulid1", ulid1)
				//fmt.Println("line2", line2)

				//new.Write([]byte(scanner2.Text()))

			} else {
				res <- line1
				//fmt.Println(line1, "line1")
				//fmt.Println("ulid2", ulid2)
				//fmt.Println("ulid1", ulid1)
				//fmt.Println(ulid1)
				//new.Write([]byte(scanner1.Text()))

			}
		}

	}()

	return res
}
