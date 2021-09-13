package logenc

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"

	//"strings"
	"testing"
)

//Function for test structure xml file
func Test_xml(t *testing.T) {
	//Marsh
	var v = &LogList{
		XML_RECORD_ROOT: []Log{
			{
				XML_APPNAME: "GenSafetyThreats",
				XML_APPPATH: "/usr/local/lemz/atcs/bin/gensafetythreats",
				XML_APPPID:  "app_pid",
				XML_THREAD:  "asd",
				XML_TIME:    "saDas",
				XML_ULID:    "asd",
				XML_TYPE:    "dsfhjshdfgioujs",
				XML_MESSAGE: "sdhjfgoshfgoihjws",
				XML_DETAILS: "sdhjfgoshfgoasdasdasdihjws",
				DT_FORMAT:   "sdhjfgoshfgoasdasfdlskajdflsjhfihjws",
			},
		},
	}

	output, err := xml.Marshal(v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	t.Logf("%s", output)

	//Unmarsh
	var q LogList
	err = xml.Unmarshal([]byte(output), &v)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	t.Logf("%#+v", v)
	//

	file, err := os.Create("logstr.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, logstr := range q.XML_RECORD_ROOT {
		err := writer.Write([]string{logstr.XML_APPNAME, logstr.XML_APPPATH, logstr.XML_APPPID, logstr.XML_THREAD, logstr.XML_TIME, logstr.XML_ULID, logstr.XML_TYPE, logstr.XML_MESSAGE, logstr.XML_DETAILS, logstr.DT_FORMAT})
		checkError("Cannot write to file", err)
	}

}

func TestDecodeLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				line: `B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZfF5VelJJWElaXU9IGRtaS0tkS1pPUwYZFE5ISRRXVFhaVxRXXlZBFFpPWEgUWVJVFFxeVVpSSVhJWl1PSBkbWktLZEtSXwYZCgkNAwIZG09TSV5aX2RSXwYZCAsMDAMLCgMPCxkbT1JWXgYZCwILDAkLCQoJCA4CCwMCDgIZG05XUl8GGQgLYwhzeGkLCwt6CnV1DHlofGhsC39zCwtzGRtPQkteBhkLGRtWXkhIWlxeBhnroeuF6rvqu+uO64HqveuD6rQb6rnqu+uO64HrixvrhOuFG+uE64Dri+uG6rgVGRteQ09kVl5ISFpcXgYZcnV9dBkUBQcUV1RcV1JITwU=`,
			},
			want: `<loglist><log module_name="GenAircrafts" app_path="/usr/local/lemz/atcs/bin/genaircrafts" app_pid="12689" thread_id="3077801840" time="09072021235908959" ulid="30X3HCR000A1NN7BSGSW0DH00H" type="0" message="Коррекция трека по плану." ext_message="INFO"/></loglist>
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if got := DecodeLine(tt.args.line); strings.Compare(got, tt.want) != 0 {
			// t.Log(len(got))
			// t.Log(len(tt.want))
			// t.Errorf("DecodeLine() = %v, want %v", got, tt.want)
			//}
		})
	}
}

//<loglist><log module_name="GenAircrafts" app_path="/usr/local/lemz/atcs/bin/genaircrafts" app_pid="12689" thread_id="3077801840" time="09072021235908959" ulid="30X3HCR000A1NN7BSGSW0DH00H" type="0" message="Коррекция трека по плану." ext_message="INFO"/></loglist>
//<loglist><log module_name="GenAircrafts" app_path="/usr/local/lemz/atcs/bin/genaircrafts" app_pid="12689" thread_id="3077801840" time="09072021235908959" ulid="30X3HCR000A1NN7BSGSW0DH00H" type="0" message="Коррекция трека по плану." ext_message="INFO"/></loglist>
