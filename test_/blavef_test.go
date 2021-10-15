package test

import (
	"testing"

	"gitlab.topaz-atcs.com/tmcs/logi2/bleveSI"
)

func BenchmarkProcFileBreveTESTSPEED(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcFileBreveSPEED("12-09-20211", "/home/nik/projects/Course/tmcs-log-agent-storage/13-09-2021")

}

func BenchmarkProcFileBreve(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcFileBreve("12-09-20211", "/home/nik/projects/Course/tmcs-log-agent-storage/13-09-2021")

}

func BenchmarkProcFileBreveS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcFileBreveSLOWLY("12-09-20211", "/home/nik/projects/Course/tmcs-log-agent-storage/13-09-2021")

}
