package logenc

import (
	"reflect"
	"testing"
	"time"
)

func TestGenTestULID(t *testing.T) {
	tests := []struct {
		name    string
		notwant Log
	}{
		{
			name:    "TestGenTestLogWithULID",
			notwant: Log{XML_ULID: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var listlog Log
			listlog.GenTestULID(time.Now())
			if got := listlog.XML_ULID; reflect.DeepEqual(got, tt.notwant) {
				t.Errorf("GenTestLogWithULID() = %v, want %v", got, tt.notwant)
			}
		})
	}
}

func Test_Datestr2time(t *testing.T) {
	type args struct {
		in0 string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "t1",
			args: args{in0: "08092021224536920"},
			want: time.Date(2021, 9, 8, 22, 45, 36, 920000000, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := datestr2time(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("datestr2time() = %v, want %v", got, tt.want)
			}
		})
	}
}
