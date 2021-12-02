package bleveSI

import (
	"reflect"
	"testing"
)

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcBleveSearchv2(tt.args.fileN, tt.args.word); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcBleveSearchv2() = %v, want %v", got, tt.want)
			}
		})
	}
}
