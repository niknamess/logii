package generate_logs

import (
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/oklog/ulid/v2"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

type LogList struct {
	XMLName         xml.Name `xml:"loglist"` //dont touch XMLName
	XML_RECORD_ROOT []Log    `xml:"log"`
}
type Log struct {
	XML_APPNAME string `xml:"module_name,attr"`
	XML_APPPATH string `xml:"app_path,attr"`
	XML_APPPID  string `xml:"app_pid,attr"`
	XML_THREAD  string `xml:"thread_id,attr"`
	XML_TIME    string `xml:"time,attr"`
	XML_ULID    string `xml:"ulid,attr"`
	XML_TYPE    string `xml:"type,attr"`
	XML_MESSAGE string `xml:"message,attr"`
	XML_DETAILS string `xml:"ext_message,attr"`
	DT_FORMAT   string `xml:"ddMMyyyyhhmmsszzz,omitempty"`
}

var (
	Logger     *log.Logger
	label      string
	labeld     string
	countFile  int = 0
	countFiled int = 0
)

const (
	XOR_KEY = 59
)

func StructFile() string {
	elem := "\""
	r := rand.New(rand.NewSource(99))
	XML_DETAILS := "Context:  -- void tmcs::AbstractMonitor::,"
	now := time.Now().UnixNano()
	entropy := rand.New(rand.NewSource(now))
	timestamp := ulid.Timestamp(time.Now())
	XML_APPNAME := strconv.Itoa(r.Intn(10)) + "TMCS TEST"
	XML_APPPATH := "/" + strconv.Itoa(r.Intn(10)) + "/TEST/TEST"
	XML_APPPID := strconv.Itoa(r.Intn(1000)) + "" // "7481,"
	XML_THREAD := strconv.Itoa(r.Intn(10)) + ""   //"88,"
	XML_MESSAGE := "Состояние '" + randomdata.IpV4Address() + "Cервер КС_UDP/Пинг'"
	XML_TYPE := strconv.Itoa(rand.Intn(4-1) + 1)
	address := randomdata.ProvinceForCountry("GB")
	time1 := "29052021000147040"
	time_ulid := ulid.MustNew(timestamp, entropy)
	ulid1 := time_ulid.String()
	LINE := "<loglist><log module_name=" + elem + XML_APPNAME + elem +
		" app_path=" + elem + XML_APPPATH + elem +
		" app_pid=" + elem + XML_APPPID + elem +
		" thread_id=" + elem + XML_THREAD + elem +
		" time=" + elem + time1 + elem +
		" ulid=" + elem + ulid1 + elem +
		" type=" + elem + XML_TYPE + elem +
		" message=" + elem + XML_MESSAGE + elem +
		" ext_message=" + elem + XML_DETAILS + address + elem + "/></loglist>"

	rand.Seed(time.Now().UnixNano())

	return LINE
}

func ProcGenN() {
	//Example()

	filesFrom := string(util.GetOutboundIP()[len(util.GetOutboundIP())-3:])
	//	last3  := string(s[len(s)-3:])
	logenc.CreateDir("./genrlogs", "")

	for true {

		LINE := StructFile()

		rand.Seed(time.Now().UnixNano())

		file, err := os.OpenFile("./genrlogs./gen_logs_coded"+filesFrom+label, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		InfoLogger := log.New(file, "", 0)

		fi, err := file.Stat()
		if err != nil {

		}
		if fi.Size() >= 200000 {
			countFile++
			fmt.Println(fi.Size())
			label = strconv.Itoa(countFile)

		}

		infof := func(info string) {
			InfoLogger.Output(2, logenc.EncodeLine(info))
		}

		infof(LINE)

		time.Sleep(time.Nanosecond * 1000000)

		if countFile >= 10 {
			return
		}

	}
}
