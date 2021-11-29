package test

import (
	"testing"

	"gitlab.topaz-atcs.com/tmcs/logi2/bleveSI"
)

func BenchmarkProcFileBreve(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBlev("test4", "./view/22-06-2021")

}

func ABenchmarkProcFileBreveS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcFileBreveSLOWLY("test777", "./view/22-06-2021")

}

func ABenchmarkProcFileSSS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveScorch("test123", "./view/22-06-2021")
	//generator.Remove("/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded193","test")

}

func ABenchmarkProcFileBatch(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBlev("test87", "./view/22-06-2021")
	//generator.Remove("/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded193","test")

}
func ATestProcBleveS(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBlev("test7", "./view/22-06-2021")

}

func ATestProcBleveScorch(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveScorch("test5", "./view/22-06-2021")
}
