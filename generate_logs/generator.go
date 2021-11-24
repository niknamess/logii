package generator

import (
	"encoding/xml"
	"log"
	"math/rand"
	"os"
	"strconv"

	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/oklog/ulid/v2"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
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
	Logger        *log.Logger
)

const (
	XOR_KEY = 59
)

func init() {
	file, err := os.OpenFile("gen_logs_coded", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	filed, err := os.OpenFile("gen_logs_decoded", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "", 0)
	//WarningLogger = log.New(file, "", 0)
	ErrorLogger = log.New(file, "", 0)

	Logger = log.New(filed, "", 0)

}

func ProcGenN() {

	i := 0
	r := rand.New(rand.NewSource(99))
	XML_DETAILS := "Context:  -- void tmcs::AbstractMonitor::,"
	//address := randomdata.ProvinceForCountry("GB")
	line1 := "<loglist><log module_name=\"TMCS Monitor\" app_path=\"/usr/local/Lemz/tmcs/monitor/tmcs_monitor\" app_pid=\"4913\" thread_id=\"\" time=\"29052021000147040\"ulid=\"0001GB313BF4HPFYCDY3QTZ6A6\" type=\"3\" message=\"Состояние '[192.168.1.120] Cервер КС_RLI/КСВ Топаз' изменилось на 'Ошибка'\" ext_message=\"Context:  -- void tmcs::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/tmcs_plugin/src/AbstractMonitor.cpp:686\"/></loglist>"
	// <loglist><log module_name=7TMCS TEST app_path=3/TEST/TEST app_pid=290 thread_id=2 time=29052021000147040 ulid=01FMSHJW4C0R9RQJ5VSWWZ0PRK type=3 message=Состояние '9.23.107.141Cервер КС_UDP/Пинг' ext_message=Context:  -- void tmcs::AbstractMonitor::,Cheshire></loglist>
	infof := func(info string) {
		InfoLogger.Output(2, logenc.EncodeLine(info))
	}
	//warnof := func(info string) {
	//	WarningLogger.Output(2, logenc.EncodeLine(info))
	//}
	erorof := func(info string) {
		ErrorLogger.Output(2, logenc.EncodeLine(info))
	}
	decode := func(info string) {
		Logger.Output(2, (info))
	}
	for true {

		now := time.Now().UnixNano()
		entropy := rand.New(rand.NewSource(now))
		timestamp := ulid.Timestamp(time.Now())
		XML_APPNAME := strconv.Itoa(r.Intn(10)) + "TMCS TEST"
		XML_APPPATH := "/" + strconv.Itoa(r.Intn(10)) + "/TEST/TEST"
		XML_APPPID := strconv.Itoa(r.Intn(1000)) + "" // "7481,"
		XML_THREAD := strconv.Itoa(r.Intn(10)) + ""   //"88,"
		XML_MESSAGE := "Состояние '" + randomdata.IpV4Address() + "Cервер КС_UDP/Пинг'"
		XML_TYPE := strconv.Itoa(rand.Intn(4-1) + 1)
		address := randomdata.ProvinceForCountry("GB") + "\n"
		//time11 := randomdata.FullDate() + ","
		time1 := "29052021000147040"
		//file, err := os.OpenFile("test"+strconv.Itoa(i), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//defer file.Close()
		time_ulid := ulid.MustNew(timestamp, entropy)
		ulid1 := time_ulid.String()
		LINE := "<loglist><log module_name=" + XML_APPNAME +
			" app_path=" + XML_APPPATH +
			" app_pid=" + XML_APPPID +
			" thread_id=" + XML_THREAD +
			" time=" + time1 +
			" ulid=" + ulid1 +
			" type=" + XML_TYPE +
			" message=" + XML_MESSAGE +
			" ext_message=" + XML_DETAILS + address + "></loglist>"

		rand.Seed(time.Now().UnixNano())
		//print("it's work ea")

		timer1 := time.NewTimer(4 * time.Second)
		//InfoLogger.Println("Starting the application...")
		infof(LINE)
		decode(LINE)

		<-timer1.C
		i++

		timer4 := time.NewTimer(2 * time.Second)
		erorof(line1)
		<-timer4.C
		i++
	}
}
