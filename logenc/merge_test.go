package logenc

import (
	"fmt"
	"testing"
)

func TestMergeLines(t *testing.T) {

	ch1 := make(chan LogList, 10)
	ch2 := make(chan LogList, 10)

	go func() {
		var listlog LogList
		//var log Log
		listlog.XML_RECORD_ROOT = make([]Log, 1)

		listlog.XML_RECORD_ROOT[0].GenTestULID()
		ch1 <- listlog
		t.Log(listlog)
		listlog.XML_RECORD_ROOT[0].GenTestULID()
		ch1 <- listlog
		t.Log(listlog)
		listlog.XML_RECORD_ROOT[0].GenTestULID()
		ch2 <- listlog
		t.Log(listlog)
		listlog.XML_RECORD_ROOT[0].GenTestULID()
		ch2 <- listlog
		ch1 <- listlog
		t.Log(listlog)
		t.Log(listlog)

		listlog.XML_RECORD_ROOT[0].GenTestULID()
		ch1 <- listlog
		t.Log(listlog)

		close(ch1)
		close(ch2)
	}()

	want := 5

	ch3 := MergeLines(ch1, ch2)

	got := 0
	for val := range ch3 {
		fmt.Println(val.XML_RECORD_ROOT[0].XML_ULID)
		got++
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
