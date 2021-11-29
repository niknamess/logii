package test

import (
	"testing"

	"gitlab.topaz-atcs.com/tmcs/logi2/bleveSI"
)

func ABenchmarkProcFileBreve(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcFileBreve("test4", "/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded1937")

}

func ABenchmarkProcFileBreveS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcFileBreveSLOWLY("test777", "/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded1936")

}

func ABenchmarkProcFileSSS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveScorch("test123", "/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded1938")
	//generator.Remove("/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded193","test")

}

func ABenchmarkProcFileBatch(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBlev("test87", "/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded1933")
	//generator.Remove("/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded193","test")

}
func ATestProcBleveS(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBlev("test7", "/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded1934")

}

func ATestProcBleveScorch(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveScorch("test5", "/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded1933")
}
