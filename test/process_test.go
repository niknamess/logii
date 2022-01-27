package test

import (
	"testing"

	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

func BenchmarkProcMapFile(b *testing.B) {
	type args struct {
		file string
	}

	logenc.ProcMapFile("./view/22-06-2021")
	//t.StartTimer()
}

func BenchmarkProcMapFilePP(b *testing.B) {
	type args struct {
		file string
	}

	logenc.ProcMapFileREZERV("./view/22-06-2021")

}
