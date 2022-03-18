package generate_logs

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"time"

	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

func ProcGenWiteF() {
	//Example()

	filesFrom := string(util.GetOutboundIP()[len(util.GetOutboundIP())-3:])
	//	last3  := string(s[len(s)-3:])
	logenc.CreateDir("./genlogs", "")
	file, err := os.OpenFile("./repdata/gen_files_write"+filesFrom, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	for {

		LINE := StructFile()

		rand.Seed(time.Now().UnixNano())

		InfoLogger := log.New(file, "", 0)

		infof := func(info string) {
			InfoLogger.Output(2, logenc.EncodeLine(info))
		}

		infof(LINE)

		//time.Sleep(time.Nanosecond * 1000000)
		time.Sleep(5000 * time.Millisecond)
		fmt.Println("Message add :D")

	}
	//fmt.Println("done")

}
