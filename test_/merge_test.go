package test

import (
	"fmt"
	//"sync"
	"testing"

	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

func TestMergeLines(t *testing.T) {
	//var wg sync.WaitGroup

	ch1 := make(chan logenc.LogList, 10)
	ch2 := make(chan logenc.LogList, 10)

	go func() {
		gen := func() logenc.LogList {
			var listlog logenc.LogList
			listlog.XML_RECORD_ROOT = make([]logenc.Log, 1)
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
	ch3 := logenc.MergeLines(ch1, ch2)
	got := 0
	for val := range ch3 {

		if len(val.XML_RECORD_ROOT) != 0 {
			got++
			fmt.Println(val.XML_RECORD_ROOT[0].XML_ULID)
		}
	}
	if got != 5 {
		t.Errorf("MergeLines() = %v, want %v", got, want)
	}

}
