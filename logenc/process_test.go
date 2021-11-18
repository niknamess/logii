package logenc

import "testing"

func TestRemoveLine(t *testing.T) {
	type args struct {
		path  string
		fileN string
		label string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemoveLine(tt.args.path, tt.args.fileN, tt.args.label)
		})
	}
}
