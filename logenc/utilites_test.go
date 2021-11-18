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
