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
var (
	count = 0
)

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
	print("start")
	print(string(data))
	return string(data)
}

func EncodeLine(line []byte) string {
	data := base64.StdEncoding.Strict().EncodeToString(line)
	result := []byte(data)

	if len(data) <= 0 {
		return ""
	}

	k := 0
	for {
		//XOR with lines
		result[k] ^= XOR_KEY
		k++
		if k >= len(data) {
			break
		}
	}
	//print(line)
	return string(result)
}

func DecodeXML(line string) (LogList, error) {
	//print("start")
	//print(line)
	var v = LogList{}

	err := xml.Unmarshal([]byte(line), &v)
	//print("end")

	return v, err
}

func EncodeXML(line string) (string, error) {

	empData1, err := xml.Marshal([]byte(line))
	empData2 := string(empData1)
	return empData2, err
}
func datestr2time(str string) time.Time {
	// format example: 08092021224536920  from xml
	const shortForm = "02012006150405.000"

	str2 := string(str[0:14]) + "." + string(str[14:17])
	//fmt.Println(str2)
	t, _ := time.Parse(shortForm, str2)
	return t
}

func EncodeCSV(val LogList) string {
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)
	for _, logstr := range val.XML_RECORD_ROOT {
		//TIME
		t := datestr2time(logstr.XML_TIME)
		//fmt.Println(logstr.XML_TIME, t, err)
		//TYPE
		typeM := "INFO"
		if logstr.XML_TYPE == "1" {
			typeM = "DEBUG"
		} else if logstr.XML_TYPE == "2" {
			typeM = "WARNING"
		} else if logstr.XML_TYPE == "3" {
			typeM = "ERROR"
		} else if logstr.XML_TYPE == "4" {
			typeM = "FATAL"
		}
		//id := fmt.Sprint(count)
		err := writer.Write([]string{typeM, logstr.XML_APPNAME, logstr.XML_APPPATH, logstr.XML_APPPID, logstr.XML_THREAD, t.Format(time.RFC1123), logstr.XML_ULID, logstr.XML_MESSAGE, logstr.XML_DETAILS, logstr.DT_FORMAT})
		count++
		if err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	writer.Flush()
	return buf.String()
}
