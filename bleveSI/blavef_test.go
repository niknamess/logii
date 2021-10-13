package bleveSI

import "testing"

func BenchmarkProcFileBreveTESTSPEED(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	ProcFileBreveSPEED("12-09-20211", "/home/nik/projects/Course/tmcs-log-agent-storage/13-09-2021")

}

func BenchmarkProcFileBreve(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	ProcFileBreve("12-09-20211", "/home/nik/projects/Course/tmcs-log-agent-storage/13-09-2021")

}
