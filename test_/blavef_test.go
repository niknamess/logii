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

func BenchmarkProcFileBreveS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcFileBreveSLOWLY("test777", "./view/22-06-2021")

}

func BenchmarkProcFileSSS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveScorch("test123", "./view/22-06-2021")
	//generator.Remove("/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded193","test")

}

func BenchmarkProcFileBatch(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBlev("test87", "./view/22-06-2021")
	//generator.Remove("/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded193","test")

}
func TestProcBleveS(t *testing.T) {
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

func TestBleveSearch(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveSearchv2("test4", "0001GD2DJBAPMGFE7W2XFHAVYT")

}

/*
func TestProcBleveSearchv2(t *testing.T) {
	type args struct {
		fileN string
		word  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "t1",
			args: args{in0: "08092021224536920",""},
			want: [ ],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bleveSI.ProcBleveSearchv2(tt.args.fileN, tt.args.word); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcBleveSearchv2() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
