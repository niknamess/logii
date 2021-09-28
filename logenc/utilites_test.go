package logenc

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_datestr2time(t *testing.T) {
	type args struct {
		in0 string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "t1",
			args: args{in0: "08092021224536920"},
			want: time.Date(2021, 9, 8, 22, 45, 36, 920000000, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := datestr2time(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("datestr2time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeXML(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    LogList
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestDecodeXML",
			//args: '<loglist><log module_name="ProcessLivePanel" app_path="/usr/local/lemz/atcs/bin/processlivepanel" app_pid="2383" thread_id="140042514286464" time="19052021050634395" type="2" message="The &quot;text&quot; attribute is missed for the button with handle {0e5d6759-f878-42dd-9548-0db403f64b19}" ext_message="WARNING"/></loglist>'},
			//want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeXML(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeXML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeXML(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t2",
			//args: args{in0: "1,INFO,TSS Service,/usr/local/Lemz/tss/tss_service,1787,ntp_cl,Fri, 20 Aug 2021 00:34:59 UTC,0001GHXY5KGAQMHEEYM7MVVNTS,Значение поля стратум в NTP ответе полученном от NTP сервера 192.168.1.252:123 изменилось: 10 -> 1.,,"},
			//: "",
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeXML(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeXML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeLine(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "t3",
			//args: args{in0: "08092021224536920"},
			//want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeLine(tt.args.line); got != tt.want {
				t.Errorf("EncodeLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeLine(t *testing.T) {
	type args struct {
		line string
	}
	line := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZb2hoG2heSU1SWF4ZG1pLS2RLWk9TBhkUTkhJFFdUWFpXFHdeVkEUT0hIFE9ISGRIXklNUlheGRtaS0tkS1JfBhkKDA8JGRtPU0leWl9kUl8GGVVPS2RYVxkbT1JWXgYZCQ0LDgkLCQoLCwoMCwMICA0ZG05XUl8GGQsLCwp8em0PbA8CeGNoen1zDWoDfWx2en9jGRtPQkteBhkLGRtWXkhIWlxeBhkdSk5UTwDrrOuG64vqvOuO64brg"
	b64data := line[strings.IndexByte(line, ',')+1:]
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "TestDecodeLine",
			args: args{line: b64data},
			want: "<loglist><log module_name=\"TSS Service\" app_path=\"/usr/local/Lemz/tss/tss_service\" app_pid=\"1742\" thread_id=\"ntp_cl\" time=\"07072021004244057\" ulid=\"0001GE9Y441W7TF1P9N5EG3HQH\" type=\"0\" message=\"&quot;Значение поля стратум в NTP ответе полученном от NTP сервера 192.168.1.252:123 изменилось: 10 -&gt; 1.&quot;\"/></loglist>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeLine(tt.args.line); got != tt.want {
				t.Errorf("DecodeLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
