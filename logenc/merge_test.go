package logenc

import (
	"fmt"
	//"sync"
	"testing"
)

func TestMergeLines(t *testing.T) {
	//var wg sync.WaitGroup

	ch1 := make(chan LogList, 10)
	ch2 := make(chan LogList, 10)

	go func() {

		gen := func() LogList {
			var listlog LogList
			//var log Log
			listlog.XML_RECORD_ROOT = make([]Log, 1)
			listlog.XML_RECORD_ROOT[0].GenTestULID()
			return listlog
		}

		ch1 <- gen()
		ch1 <- gen()
		ch2 <- gen()

		dup := gen()
		ch2 <- dup
		ch1 <- dup

		ch1 <- gen()

		close(ch1)
		close(ch2)
	}()

	want := 5

	ch3 := MergeLines(ch1, ch2)
	//wg.Wait()

	got := 0
	//wg.Add(1)
	//defer wg.Done()
	for val := range ch3 {

		if len(val.XML_RECORD_ROOT) != 0 {
			got++
			fmt.Println(val.XML_RECORD_ROOT[0].XML_ULID)
		}

		//fmt.Println(val.XML_RECORD_ROOT[0].XML_ULID)
		//got++
	}

	if got != 5 {
		t.Errorf("MergeLines() = %v, want %v", got, want)
	}

	// type args struct {
	// 	ch1 chan LogList
	// 	ch2 chan LogList
	// }
	// tests := []struct {
	// 	name string
	// 	args args
	// 	want chan LogList
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		if got := MergeLines(tt.args.ch1, tt.args.ch2); !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("MergeLines() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
