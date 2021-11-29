package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

//Remove file in dir
func Example() {
	Remove("./genrlogs./", "gen_logs_coded")
	Remove("./repdata/", "gen_logs_coded")
}

func Remove(dirpath string, lineS string) {
	var count int = 0
	files, _ := ioutil.ReadDir(dirpath)

	for _, file := range files {
		//go R(count)
		fmt.Println(count)
		fileN := file.Name()
		fmt.Println(fileN)
		contain := strings.Contains(fileN, lineS)
		if contain == true {
			logenc.DeleteOldsFiles(dirpath, fileN, "")
		}

	}

}

func DeleteFile90(dir string) {
	var cutoff = 24 * time.Hour * 70
	fileInfo, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err.Error())
	}
	now := time.Now()
	fmt.Println(now)
	for _, info := range fileInfo {
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			fmt.Printf("Deleting %s which is %s old\n", info.Name(), diff)
			logenc.DeleteOldsFiles(dir, info.Name(), "")

		}
	}
}
