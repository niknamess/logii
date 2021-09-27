package generator

import (
	"encoding/base64"
	"log"
	"math/rand"
	"os"

	//"strings"
	"time"
	//"github.com/kataras/tablewriter"
	//"github.com/lensesio/tableprinter"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

const (
	XOR_KEY = 59
	//shortForm = "2006.01.02-15.04.05"
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

	return string(result)
}

func init() {
	//file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	//WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	//ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	//FatalLogger = log.New(file, "Fatal: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func ProcGenN(dir string) {

	//file, err := os.OpenFile("/home/nik/projects/logs/r/gen_logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	file, err := os.OpenFile("gen_logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	XML_APPNAME := "TMCS Monitor"
	XML_APPPATH := "/usr/local/Lemz/tmcs/monitor/tmcs_monitor"
	XML_APPPID := "7481"
	XML_THREAD := "88"
	time1 := "Fri, 20 Aug 2021 00:44:22 UTC"
	XML_ULID := "0001GHXYQ6EM4972TMPV0E0W6Q"
	XML_MESSAGE := "Состояние '[192.168.1.128] Cервер КС_UDP/Пинг'"
	XML_DETAILS := "Context:  -- void tmcs::AbstractMonitor::"
	DT_FORMAT := ""
	//11,ERROR,TMCS Monitor,/usr/local/Lemz/tmcs/monitor/tmcs_monitor,7481,,"Fri, 20 Aug 2021 00:44:22 UTC",0001GHXYQ6EM4972TMPV0E0W6Q,Состояние '[192.168.1.128] Cервер КС_UDP/Пинг' изменилось на 'Не доступен',Context:  -- void tmcs::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/tmcs_plugin/src/AbstractMonitor.cpp:701,
	//([]string{typeM, logstr.XML_APPNAME, logstr.XML_APPPATH, logstr.XML_APPPID, logstr.XML_THREAD, t.Format(time.RFC1123), logstr.XML_ULID, logstr.XML_MESSAGE, logstr.XML_DETAILS, logstr.DT_FORMAT})
	InfoLogger = log.New(file, EncryptDecrypt([]byte("INFO")), 0)
	WarningLogger = log.New(file, EncryptDecrypt([]byte("WARNING")), 0)
	ErrorLogger = log.New(file, EncryptDecrypt([]byte("ERROR")), 0)

	infof := func(info string) {
		InfoLogger.Output(2, info)
	}

	warnof := func(info string) {
		WarningLogger.Output(2, info)
	}

	erorof := func(info string) {
		ErrorLogger.Output(2, info)
	}

	for true {
		rand.Seed(time.Now().UnixNano())
		print("it's work ea")

		timer1 := time.NewTimer(4 * time.Second)
		//InfoLogger.Println("Starting the application...")
		infof(EncryptDecrypt([]byte(XML_APPNAME + XML_APPPATH + XML_APPPID + XML_THREAD + time1 + XML_ULID + XML_MESSAGE + XML_DETAILS + DT_FORMAT)))

		<-timer1.C
		i++

		timer2 := time.NewTimer(5 * time.Second)
		infof(EncryptDecrypt([]byte(XML_APPNAME + XML_APPPATH + XML_APPPID + XML_THREAD + time1 + XML_ULID + XML_MESSAGE + XML_DETAILS + DT_FORMAT)))

		<-timer2.C
		i++
		timer3 := time.NewTimer(10 * time.Second)
		warnof(EncryptDecrypt([]byte(XML_APPNAME + XML_APPPATH + XML_APPPID + XML_THREAD + time1 + XML_ULID + XML_MESSAGE + XML_DETAILS + DT_FORMAT)))

		<-timer3.C
		i++
		timer4 := time.NewTimer(5 * time.Second)
		erorof(EncryptDecrypt([]byte(XML_APPNAME + XML_APPPATH + XML_APPPID + XML_THREAD + time1 + XML_ULID + XML_MESSAGE + XML_DETAILS + DT_FORMAT)))

		<-timer4.C
		i++
	}
}
