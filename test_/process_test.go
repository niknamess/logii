package test

import (
	"testing"

	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

func BenchmarkProcMapFile(b *testing.B) {
	type args struct {
		file string
	}

	logenc.ProcMapFile("/home/nik/projects/Course/tmcs-log-agent-storage/26-05-2021")
	//t.StartTimer()
}

func BenchmarkProcMapFilePP(b *testing.B) {
	type args struct {
		file string
	}
	//b.SetBytes(1)
	//for i := 0; i < b.N; i++ {
	//for i := 0; i < 6; i++ {

	logenc.ProcMapFileREZERV("/home/nik/projects/Course/tmcs-log-agent-storage/26-05-2021")

	//}
	//t.StartTimer()
	//}
}
