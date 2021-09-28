package generator

import (
	"encoding/base64"
	"encoding/xml"
	"log"
	"math/rand"
	"os"

	//"strings"
	"time"

	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	//"github.com/kataras/tablewriter"
	//"github.com/lensesio/tableprinter"
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
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

const (
	XOR_KEY = 59
)

func EncryptDecrypt(input []byte) (output string) {
	data := base64.StdEncoding.Strict().EncodeToString(input)
	result := []byte(data)

	if len(data) <= 0 {
		return ""
	}

	for i := 0; i < len(data); i++ {
		result[i] ^= XOR_KEY
	}
	//print(result)
	return string(result)
}

func init() {
	file, err := os.OpenFile("/home/nik/projects/logs/test/gen_logs1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "", 0)
	WarningLogger = log.New(file, "", 0)
	ErrorLogger = log.New(file, "", 0)

}

func ProcGenN(dir string) {
	i := 0
	XML_APPNAME := "TMCS Monitor,"
	XML_APPPATH := "/usr/local/Lemz/tmcs/monitor/tmcs_monitor,"
	XML_APPPID := "7481,"
	XML_THREAD := "88,"
	time1 := "Fri, 20 Aug 2021 00:43:44 UTC"
	qtype := "0"
	XML_ULID := "0001GHXYQ6EM4972TMPV0E0W6Q"
	XML_MESSAGE := "Состояние '[192.168.1.128] Cервер КС_UDP/Пинг'"
	XML_DETAILS := "Context:  -- void tmcs::AbstractMonitor::"
	address := "sajjsaj"
	line := XML_APPNAME + XML_APPPATH + XML_APPPID + XML_THREAD + time1 + qtype + XML_ULID + XML_MESSAGE + XML_DETAILS + address
	line1 := "<loglist><log module_name=\"TMCS Monitor\" app_path=\"/usr/local/Lemz/tmcs/monitor/tmcs_monitor\" app_pid=\"4913\" thread_id=\"\" time=\"29052021000147040\" ulid=\"0001GB313BF4HPFYCDY3QTZ6A6\" type=\"3\" message=\"Состояние '[192.168.1.120] Cервер КС_RLI/КСВ Топаз' изменилось на 'Ошибка'\" ext_message=\"Context:  -- void tmcs::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/tmcs_plugin/src/AbstractMonitor.cpp:686\"/></loglist>"

	//file, err := os.OpenFile("/home/nik/projects/logs/r/gen_logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	infof := func(info string) {
		InfoLogger.Output(2, logenc.EncodeLine(info))
	}

	warnof := func(info string) {
		WarningLogger.Output(2, logenc.EncodeLine(info))
	}

	erorof := func(info string) {
		ErrorLogger.Output(2, logenc.EncodeLine(info))
	}

	for true {
		rand.Seed(time.Now().UnixNano())
		print("it's work ea")

		timer1 := time.NewTimer(4 * time.Second)
		//InfoLogger.Println("Starting the application...")
		infof(line1)

		<-timer1.C
		i++

		timer2 := time.NewTimer(2 * time.Second)
		infof(line1)

		<-timer2.C
		i++
		timer3 := time.NewTimer(2 * time.Second)
		warnof(line1)
		<-timer3.C
		i++
		timer4 := time.NewTimer(2 * time.Second)
		erorof(line)
		<-timer4.C
		i++
	}
}
