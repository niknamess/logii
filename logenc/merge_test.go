package logenc

import (
	"testing"
)

func TestReplication(t *testing.T) {
	type args struct {
		path string
	}

	Replication("/home/nik/projects/Course/logi2/repdata/Test/19-05-2021")

}

func TestMerge(t *testing.T) {
	type args struct {
		path string
	}

	Merge("/home/nik/projects/Course/logi2/repdata/Test/19-05-2021")

}
