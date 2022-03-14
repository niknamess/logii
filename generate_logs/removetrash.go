package generate_logs

import (
	"io/ioutil"
	"strings"

	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

//Remove file in dir
func Example() {
	Remove("./genrlogs./", "gen_logs_coded")
	Remove("./repdata/", "gen_logs_coded")
}

func Remove(dirpath string, lineS string) {
	//var count int = 0
	files, _ := ioutil.ReadDir(dirpath)

	for _, file := range files {
		//go R(count)
		//fmt.Println(count)
		fileN := file.Name()
		//fmt.Println(fileN)
		contain := strings.Contains(fileN, lineS)
		if contain {
			logenc.DeleteOldsFiles(dirpath, fileN, "")
		}

	}

}
