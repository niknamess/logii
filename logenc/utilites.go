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

	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	for _, logstr := range val.XML_RECORD_ROOT {
		// fmt.Println(logstr.XML_APPNAME)
		err := writer.Write([]string{logstr.XML_APPNAME, logstr.XML_APPPATH, logstr.XML_APPPID, logstr.XML_THREAD, logstr.XML_TIME, logstr.XML_ULID, logstr.XML_TYPE, logstr.XML_MESSAGE, logstr.XML_DETAILS, logstr.DT_FORMAT})
		if err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	writer.Flush()
	//fmt.Println(buf.Len())
	return buf.String()
}
