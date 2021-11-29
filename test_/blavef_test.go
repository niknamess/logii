package test

import (
	"testing"

	"gitlab.topaz-atcs.com/tmcs/logi2/bleveSI"
)

<<<<<<< HEAD
func BenchmarkProcFileBreve(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBlev("test4", "./view/22-06-2021")

}

func ABenchmarkProcFileBreveS(b *testing.B) {
=======
func BenchmarkProcFileBreveTESTSPEED(b *testing.B) {
>>>>>>> 7cc21d8bc26936e7ef731a8b2d1dc24da8cf5e15
	type args struct {
		fileN string
		file  string
	}

<<<<<<< HEAD
	bleveSI.ProcFileBreveSLOWLY("test777", "./view/22-06-2021")

}

func ABenchmarkProcFileSSS(b *testing.B) {
=======
	bleveSI.ProcFileBleveSPEED("12-09-20211", "/home/nik/projects/Course/tmcs-log-agent-storage/13-09-2021")

}

func BenchmarkProcFileBreve(b *testing.B) {
>>>>>>> 7cc21d8bc26936e7ef731a8b2d1dc24da8cf5e15
	type args struct {
		fileN string
		file  string
	}

<<<<<<< HEAD
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
=======
	bleveSI.ProcFileBreve("22-06-20211", "./view/22-06-2021")

}

func BenchmarkProcFileBreveS(b *testing.B) {
>>>>>>> 7cc21d8bc26936e7ef731a8b2d1dc24da8cf5e15
	type args struct {
		fileN string
		file  string
	}

<<<<<<< HEAD
	bleveSI.ProcBlev("test7", "./view/22-06-2021")

}

func ATestProcBleveScorch(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveScorch("test5", "./view/22-06-2021")
=======
	bleveSI.ProcFileBreveSLOWLY("22-06-20211", "./view/22-06-2021")

>>>>>>> 7cc21d8bc26936e7ef731a8b2d1dc24da8cf5e15
}
