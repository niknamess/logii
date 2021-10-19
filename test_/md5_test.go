package test

import (
	"testing"

	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

func BenchmarkCheckFileSum(b *testing.B) {
	type args struct {
		file string
	}

	logenc.CheckFileSum("/home/nik/projects/Course/logi2/logtest/gen_logs")
	logenc.CheckFileSum("/home/nik/projects/Course/logi2/logtest/gen_logs1")

}
