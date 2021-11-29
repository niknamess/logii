package test

import (
	//"fmt"
	//"math/rand"

	//"sync"
	"testing"

	//"github.com/oklog/ulid"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
)

/*
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
			if got := logenc.datestr2time(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("datestr2time() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
/*func TestDecodeXML(t *testing.T) {
	lines := "<loglist><log module_name=\"TMCS Monitor\" app_path=\"/usr/local/Lemz/tmcs/monitor/tmcs_monitor\" app_pid=\"4913\" thread_id=\"\" time=\"29052021000147040\" ulid=\"0001GB313BF4HPFYCDY3QTZ6A6\" type=\"3\" message=\"Состояние '[192.168.1.120] Cервер КС_RLI/КСВ Топаз' изменилось на 'Ошибка'\" ext_message=\"Context:  -- void tmcs::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/tmcs_plugin/src/AbstractMonitor.cpp:686\"/></loglist>"

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
			name: "TestDecodeLine",
			args: args{line: lines},
			//want: {{ loglist} [{TMCS Monitor /usr/local/Lemz/tmcs/monitor/tmcs_monitor 4913  29052021000147040 0001GB313BF4HPFYCDY3QTZ6A6 3 Состояние '[192.168.1.120] Cервер КС_RLI/КСВ Топаз' изменилось на 'Ошибка' Context:  -- void tmcs::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/tmcs_plugin/src/AbstractMonitor.cpp:686 }]},
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
*/
func TestEncodeLine(t *testing.T) {
	line := "<loglist><log module_name=\"TMCS Monitor\" app_path=\"/usr/local/Lemz/tmcs/monitor/tmcs_monitor\" app_pid=\"4913\" thread_id=\"\" time=\"29052021000147040\" ulid=\"0001GB313BF4HPFYCDY3QTZ6A6\" type=\"3\" message=\"Состояние '[192.168.1.120] Cервер КС_RLI/КСВ Топаз' изменилось на 'Ошибка'\" ext_message=\"Context:  -- void tmcs::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/tmcs_plugin/src/AbstractMonitor.cpp:686\"/></loglist>"

	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "TestEncodeLine",
			args: args{line: line},
			want: "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZb3Z4aBt2VFVST1RJGRtaS0tkS1pPUwYZFE5ISRRXVFhaVxR3XlZBFE9WWEgUVlRVUk9USRRPVlhIZFZUVVJPVEkZG1pLS2RLUl8GGQ8CCggZG09TSV5aX2RSXwYZGRtPUlZeBhkJAgsOCQsJCgsLCwoPDAsPCxkbTldSXwYZCwsLCnx5CAoIeX0Pc2t9Ynh/Yghqb2ENeg0ZG09CS14GGQgZG1ZeSEhaXF4GGeua64Xquuq564XqtOuG64PrjhscYAoCCRUKDQMVChUKCQtmG3jrjuq764nrjuq7G+uh65pkaXdyFOuh65rrqRvrmeuF64Tri+uMHBvrg+uM64frjuuG64PrgOuF6rrqtxvrhuuLGxzrpeqz64PriuuB64scGRteQ09kVl5ISFpcXgYZeFRVT15DTwEbGxYWG01UUl8bT1ZYSAEBellIT0laWE92VFVST1RJAQFUVXhUVktUVV5VT2hPWk9eeFNaVVxeXxNqbk5SXxIAGxUVFBUVFBUVFBUVFEhJWBRXUllIFE9WWEhkS1dOXFJVFEhJWBR6WUhPSVpYT3ZUVVJPVEkVWEtLAQ0DDRkUBQcUV1RcV1JITwU=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logenc.EncodeLine(tt.args.line); got != tt.want {
				t.Errorf("EncodeLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

<<<<<<< HEAD
func ATestDecodeLine(t *testing.T) {
	type args struct {
		line string
	}
	//line1 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZb3Z4aBt2VFVST1RJGRtaS0tkS1pPUwYZFE5ISRRXVFhaVxR3XlZBFE9WWEgUVlRVUk9USRRPVlhIZFZUVVJPVEkZG1pLS2RLUl8GGQ8CCggZG09TSV5aX2RSXwYZGRtPUlZeBhkJAgsOCQsJCgsLCwoPDAsPCxkbTldSXwYZCwsLCnx5CAoIeX0Pc2t9Ynh/Yghqb2ENeg0ZG09CS14GGQgZG1ZeSEhaXF4GGeua64Xquuq564XqtOuG64PrjhscYAoCCRUKDQMVChUKCQtmG3jrjuq764nrjuq7G+uh65pkaXdyFOuh65rrqRvrmeuF64Tri+uMHBvrg+uM64frjuuG64PrgOuF6rrqtxvrhuuLGxzrpeqz64PriuuB64scGRteQ09kVl5ISFpcXgYZeFRVT15DTwEbGxYWG01UUl8bT1ZYSAEBellIT0laWE92VFVST1RJAQFUVXhUVktUVV5VT2hPWk9eeFNaVVxeXxNqbk5SXxIAGxUVFBUVFBUVFBUVFEhJWBRXUllIFE9WWEhkS1dOXFJVFEhJWBR6WUhPSVpYT3ZUVVJPVEkVWEtLAQ0DDRkUBQcUV1RcV1JITwU="
	line2 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZDG92eGgbb35obxkbWktLZEtaT1MGGRQIFG9+aG8Ub35obxkbWktLZEtSXwYZCQILGRtPU0leWl9kUl8GGQkZG09SVl4GGQkCCw4JCwkKCwsLCg8MCw8LGRtOV1JfBhkLCn11eWlxDHkDcX5xYnZ9fwMJfQtvfX94CxkbT0JLXgYZCRkbVl5ISFpcXgYZ65rrheq66rnrheq064brg+uOGxwCAxUKCQoVCgMKFQoJA3jrjuq764nrjuq7G+uh65pkbn9rFOuk64PrhuuIHBkbXkNPZFZeSEhaXF4GGXhUVU9eQ08BGxsWFhtNVFJfG09WWEgBAXpZSE9JWlhPdlRVUk9USQEBF3xJXlpPXkkbd1RVX1RVGQUHFFdUXFdSSE8F"
=======
func TestDecodeLine(t *testing.T) {
	type args struct {
		line string
	}
	line1 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZb3Z4aBt2VFVST1RJGRtaS0tkS1pPUwYZFE5ISRRXVFhaVxR3XlZBFE9WWEgUVlRVUk9USRRPVlhIZFZUVVJPVEkZG1pLS2RLUl8GGQ8CCggZG09TSV5aX2RSXwYZGRtPUlZeBhkJAgsOCQsJCgsLCwoPDAsPCxkbTldSXwYZCwsLCnx5CAoIeX0Pc2t9Ynh/Yghqb2ENeg0ZG09CS14GGQgZG1ZeSEhaXF4GGeua64Xquuq564XqtOuG64PrjhscYAoCCRUKDQMVChUKCQtmG3jrjuq764nrjuq7G+uh65pkaXdyFOuh65rrqRvrmeuF64Tri+uMHBvrg+uM64frjuuG64PrgOuF6rrqtxvrhuuLGxzrpeqz64PriuuB64scGRteQ09kVl5ISFpcXgYZeFRVT15DTwEbGxYWG01UUl8bT1ZYSAEBellIT0laWE92VFVST1RJAQFUVXhUVktUVV5VT2hPWk9eeFNaVVxeXxNqbk5SXxIAGxUVFBUVFBUVFBUVFEhJWBRXUllIFE9WWEhkS1dOXFJVFEhJWBR6WUhPSVpYT3ZUVVJPVEkVWEtLAQ0DDRkUBQcUV1RcV1JITwU="
	//line2 := "B1dUXFdSSE8FB1dUXBtWVF9OV15kVVpWXgYZb3Z4aBt2VFVST1RJGRtaS0tkS1pPUwYZFE5ISRRXVFhaVxR3XlZBFE9WWEgUVlRVUk9USRRPVlhIZFZUVVJPVEkZG1pLS2RLUl8GGQ8CCggZG09TSV5aX2RSXwYZGRtPUlZeBhkJAgsOCQsJCgsLCwoPDAsPCxkbTldSXwYZCwsLCnx5CAoIeX0Pc2t9Ynh/Yghqb2ENeg0ZG09CS14GGQgZG1ZeSEhaXF4GGeua64Xquuq564XqtOuG64PrjhscYAoCCRUKDQMVChUKCQtmG3jrjuq764nrjuq7G+uh65pkaXdyFOuh65rrqRvrmeuF64Tri+uMHBvrg+uM64frjuuG64PrgOuF6rrqtxvrhuuLGxzrpeqz64PriuuB64scGRteQ09kVl5ISFpcXgYZeFRVT15DTwEbGxYWG01UUl8bT1ZYSAEBellIT0laWE92VFVST1RJAQFUVXhUVktUVV5VT2hPWk9eeFNaVVxeXxNqbk5SXxIAGxUVFBUVFBUVFBUVFEhJWBRXUllIFE9WWEhkS1dOXFJVFEhJWBR6WUhPSVpYT3ZUVVJPVEkVWEtLAQ0DDRkUBQcUV1RcV1JITwU="
>>>>>>> 7cc21d8bc26936e7ef731a8b2d1dc24da8cf5e15
	//b64data := line[strings.IndexByte(line, ',')+1:]
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "TestDecodeLine",
<<<<<<< HEAD
			args: args{line: line2}, //"01FNBRJ7B8JEJYMFD82F0TFDC0"
=======
			args: args{line: line1},
>>>>>>> 7cc21d8bc26936e7ef731a8b2d1dc24da8cf5e15
			want: "<loglist><log module_name=\"TMCS Monitor\" app_path=\"/usr/local/Lemz/tmcs/monitor/tmcs_monitor\" app_pid=\"4913\" thread_id=\"\" time=\"29052021000147040\" ulid=\"0001GB313BF4HPFYCDY3QTZ6A6\" type=\"3\" message=\"Состояние '[192.168.1.120] Cервер КС_RLI/КСВ Топаз' изменилось на 'Ошибка'\" ext_message=\"Context:  -- void tmcs::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/tmcs_plugin/src/AbstractMonitor.cpp:686\"/></loglist>",
		},
		//{
		//	name: "TestDecodeLine",
		//	args: args{line: line2},
		//	want: "1<loglist><log module_name=\"TMCS Monitor\" app_path=\"/usr/local/Lemz/tmcs/monitor/tmcs_monitor\" app_pid=\"4913\" thread_id=\"\" time=\"29052021000147040\" ulid=\"0001GB313BF4HPFYCDY3QTZ6A6\" type=\"3\" message=\"Состояние '[192.168.1.120] Cервер КС_RLI/КСВ Топаз' изменилось на 'Ошибка'\" ext_message=\"Context:  -- void tmcs::AbstractMonitor::onComponentStateChanged(QUuid); ../../../../src/libs/tmcs_plugin/src/AbstractMonitor.cpp:686\"/></loglist>",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logenc.DecodeLine(tt.args.line); got != tt.want {
				t.Errorf("DecodeLine() = %v, want %v", got, tt.want)

			}
		})
	}
}

/*
func TestParseUlid(t *testing.T) {
	line1:="0001GHXYQ6EM4972TMPV0E0W6Q"
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "t1",
			args: 	 args{line: line1},
			want: 	" ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := datestr2time(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("datestr2time() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
