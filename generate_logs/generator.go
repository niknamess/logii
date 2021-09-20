package main

import (
	"log"
	"math/rand"
	"os"
	"time"
	//"github.com/kataras/tablewriter"
	//"github.com/lensesio/tableprinter"
)

type logs struct {
	typeLogs string `header:"first name"`
	Lastname string `header:"last name"`
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

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

func main() {

	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	for true {
		rand.Seed(time.Now().UnixNano())

		timer1 := time.NewTimer(4 * time.Second)
		InfoLogger.Println("Starting the application...")
		<-timer1.C

		timer2 := time.NewTimer(5 * time.Second)
		InfoLogger.Println("Something noteworthy happened")
		<-timer2.C

		timer3 := time.NewTimer(10 * time.Second)
		WarningLogger.Println("There is something you should know about")
		<-timer3.C

		timer4 := time.NewTimer(5 * time.Second)
		ErrorLogger.Println("Something went wrong")
		<-timer4.C
	}
}
