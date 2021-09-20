package logenc

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"time"
)

//XML_Structure
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

//pointer
//type TypeFile struct {
//	path string
//	dir  []string
//}

const (
	XOR_KEY = 59
	//shortForm = "2006.01.02-15.04.05"
)

//Read lines
func ReadLines(path string, fn func(line string)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fn(scanner.Text())
	}
	return scanner.Err()
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func DecodeLine(line string) string {
	data, err := base64.StdEncoding.DecodeString(line)

	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	if len(data) <= 0 {
		return ""
	}

	k := 0
	for {
		//XOR with lines
		data[k] ^= XOR_KEY
		k++
		if k >= len(data) {
			break
		}
	}
	return string(data)
}

func DecodeXML(line string) (LogList, error) {

	var v = LogList{}

	err := xml.Unmarshal([]byte(line), &v)
	return v, err
}

func EncodeCSV(val LogList) string {

	// src time format example: 08092021224536920
	//                 ddMMyyyyhhmmsszzz
	const shortForm = "02012006150405.000"

	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)
	tochka := "."
	for _, logstr := range val.XML_RECORD_ROOT {
		//TIME
		count := string(logstr.XML_TIME[0]) + string(logstr.XML_TIME[1]) + string(logstr.XML_TIME[2]) + string(logstr.XML_TIME[3]) + string(logstr.XML_TIME[4]) + string(logstr.XML_TIME[5]) + string(logstr.XML_TIME[6]) + string(logstr.XML_TIME[7]) + string(logstr.XML_TIME[8]) + string(logstr.XML_TIME[9]) + string(logstr.XML_TIME[10]) + string(logstr.XML_TIME[11]) + string(logstr.XML_TIME[12]) + string(logstr.XML_TIME[13]) + tochka + string(logstr.XML_TIME[14]) + string(logstr.XML_TIME[15]) + string(logstr.XML_TIME[16])
		t, err := time.Parse(shortForm, count)
		//fmt.Println(logstr.XML_TIME, t, err)
		//TYPE

		err = writer.Write([]string{logstr.XML_APPNAME, logstr.XML_APPPATH, logstr.XML_APPPID, logstr.XML_THREAD, t.Format(time.RFC1123), logstr.XML_ULID, logstr.XML_TYPE, logstr.XML_MESSAGE, logstr.XML_DETAILS, logstr.DT_FORMAT})
		if err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	writer.Flush()
	return buf.String()
}
