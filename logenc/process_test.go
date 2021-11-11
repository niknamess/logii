package logenc

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/oklog/ulid/v2"
)

func TestComparefiles2(t *testing.T) {

	tests := []struct {
		nums1  []int
		num2   []int
		result []int
	}{
		{
			nums1:  []int{},
			num2:   []int{},
			result: []int{},
		},
		{
			nums1:  []int{1, 2, 3},
			num2:   []int{},
			result: []int{1, 2, 3},
		},
	}

	r := rand.New(rand.NewSource(99))

	qtype := "0,"
	//	XML_MESSAGE := "Состояние '" + randomdata.IpV4Address() + "Cервер КС_UDP/Пинг',"
	XML_DETAILS := "Context:  -- void tmcs::AbstractMonitor::,"
	//address := randomdata.ProvinceForCountry("GB")

	now := time.Now().UnixNano()
	entropy := rand.New(rand.NewSource(now))
	timestamp := ulid.Timestamp(time.Now())
	for i, _ := range tests {
		ML_APPNAME := strconv.Itoa(r.Intn(10)) + "TMCS TEST,"
		XML_APPPATH := strconv.Itoa(r.Intn(10)) + "/TEST/TEST,"
		XML_APPPID := strconv.Itoa(r.Intn(1000)) + "," // "7481,"
		XML_THREAD := strconv.Itoa(r.Intn(10)) + ","   //"88,"
		XML_MESSAGE := "Состояние '" + randomdata.IpV4Address() + "Cервер КС_UDP/Пинг',"
		address := randomdata.ProvinceForCountry("GB") + "\n"
		time1 := randomdata.FullDate() + ","
		file, err := os.OpenFile("test"+strconv.Itoa(i), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		time_ulid := ulid.MustNew(timestamp, entropy)
		ulid1 := time_ulid.String()
		//mkulid
		LINE := ML_APPNAME + XML_APPPATH + XML_APPPID + XML_THREAD + time1 + qtype + ulid1 + XML_MESSAGE + XML_DETAILS + address
		//write to file
		file.Write([]byte(LINE))
	}
	/*
		for i, k := range tests {
			result, err := os.OpenFile("test", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}
			file1, err := os.OpenFile("test0", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer file1.Close()

			file2, err := os.OpenFile("test1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}
			defer file2.Close()
			scanner1 := bufio.NewScanner(file1)

			scanner2 := bufio.NewScanner(file2)

			info1, err := os.Stat("test0")
			info2, err := os.Stat("test1")
			if info1.Size() > info2.Size() {
				for scanner1.Scan() {

				}
			} else if info1.Size() < info2.Size() {
				for scanner2.Scan() {

				}
			} else {
				for scanner1.Scan() {

				}
			}
	*/
	//Comparefiles2(
	//"test0",
	//"test1",
	//"Test")

	//}

}
