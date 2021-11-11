package logenc

import "github.com/oklog/ulid/v2"

func MergeLines(ch1 chan LogList, ch2 chan LogList) chan LogList {
	res := make(chan LogList)

	go func() {
		for {
			var line1 LogList
			var line2 LogList
			var ulid1 ulid.ULID
			var ulid2 ulid.ULID
			var ok1, ok2 bool

			select {
			case line1, ok1 = <-ch1:
				if len(line1.XML_RECORD_ROOT) != 0 {
					ulid1, _ = ulid.ParseStrict(line1.XML_RECORD_ROOT[0].XML_ULID)
				}

			case line2, ok2 = <-ch2:
				if len(line2.XML_RECORD_ROOT) != 0 {
					ulid2, _ = ulid.ParseStrict(line2.XML_RECORD_ROOT[0].XML_ULID)
				}
			}

			if !ok1 && !ok2 {
				close(res)
				return
			}

			if ulid1.Compare(ulid2) == 1 {
				//new.Write([]byte(scanner2.Text()))

			} else {
				//new.Write([]byte(scanner1.Text()))

			}
		}

	}()

	return res
}
