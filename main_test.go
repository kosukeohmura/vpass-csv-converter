package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Test_loadSrcRecords(t *testing.T) {
	type args struct {
		srcReader  io.Reader
		isSrcFixed bool
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name: "Can load as expected - fixed source",
			args: args{
				srcReader: transform.NewReader(
					bytes.NewBufferString(`大村　幸佑　様,0000-0000-0000-0***,ＶＩＳＡ
2022/08/05,ヨドバシカメラ　通信販売,4853,１,１,4853,
2022/08/15,ＡＭＡＺＯＮ．ＣＯ．ＪＰ,3400,１,１,3400,
大村　幸佑　様,0000-0000-0000-0***,ＡｐｐｌｅＰａｙ／ｉＤ
2022/08/18,ファミリーマート／ｉＤ,340,１,１,340,ﾌｱﾐﾘ-ﾏ-ﾄ/ID
,,,,,123456
,,,,,123456,
`),
					japanese.ShiftJIS.NewEncoder()),
				isSrcFixed: true,
			},
			want: [][]string{
				{"2022/08/05", "ヨドバシカメラ　通信販売", "4853", "１", "１", "4853", ""},
				{"2022/08/15", "ＡＭＡＺＯＮ．ＣＯ．ＪＰ", "3400", "１", "１", "3400", ""},
				{"2022/08/18", "ファミリーマート／ｉＤ", "340", "１", "１", "340", "ﾌｱﾐﾘ-ﾏ-ﾄ/ID"},
			},
			wantErr: false,
		},
		{
			name: "Can load as expected - non-fixed source",
			args: args{
				srcReader: transform.NewReader(
					bytes.NewBufferString(`2022/7/4,ＡＭＡＺＯＮ．ＣＯ．ＪＰ,ご家族,1回払い,,'22/08,3844,3844,,,,,
2022/7/13,東京都水道局,ご家族,1回払い,,'22/08,8459,8459,,,,,
2022/7/16,セブン－イレブン／ｉＤ,ご家族,1回払い,,'22/08,98,98,,,,,
2022/7/19,メルカリ,ご家族,1回払い,,'22/08,2700,2700,,,,,
`),
					japanese.ShiftJIS.NewEncoder()),
				isSrcFixed: false,
			},
			want: [][]string{
				{"2022/7/4", "ＡＭＡＺＯＮ．ＣＯ．ＪＰ", "ご家族", "1回払い", "", "'22/08", "3844", "3844", "", "", "", "", ""},
				{"2022/7/13", "東京都水道局", "ご家族", "1回払い", "", "'22/08", "8459", "8459", "", "", "", "", ""},
				{"2022/7/16", "セブン－イレブン／ｉＤ", "ご家族", "1回払い", "", "'22/08", "98", "98", "", "", "", "", ""},
				{"2022/7/19", "メルカリ", "ご家族", "1回払い", "", "'22/08", "2700", "2700", "", "", "", "", ""},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadSrcRecords(tt.args.srcReader, tt.args.isSrcFixed)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadSrcRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadSrcRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFixedSrcRecords(t *testing.T) {
	type args struct {
		records [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    []*srcRecord
		wantErr bool
	}{
		{
			name: "Can parse as expected - non-fixed source",
			args: args{
				records: [][]string{
					{"2022/08/05", "ヨドバシカメラ　通信販売", "4853", "１", "１", "4853", ""},
					{"2022/08/15", "ＡＭＡＺＯＮ．ＣＯ．ＪＰ", "3400", "１", "１", "3400", ""},
					{"2022/08/18", "ファミリーマート／ｉＤ", "340", "１", "１", "340", "ﾌｱﾐﾘ-ﾏ-ﾄ/ID"},
				},
			},
			want: []*srcRecord{
				{Date: "2022/08/05", ShopName: "ヨドバシカメラ　通信販売", Amount: 4853},
				{Date: "2022/08/15", ShopName: "ＡＭＡＺＯＮ．ＣＯ．ＪＰ", Amount: 3400},
				{Date: "2022/08/18", ShopName: "ファミリーマート／ｉＤ ﾌｱﾐﾘ-ﾏ-ﾄ/ID", Amount: 340},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFixedSrcRecords(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFixedSrcRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFixedSrcRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNonFixedSrcRecords(t *testing.T) {
	type args struct {
		records [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    []*srcRecord
		wantErr bool
	}{
		{
			name: "Can parse as expected - non-fixed source",
			args: args{
				records: [][]string{
					{"2022/7/4", "ＡＭＡＺＯＮ．ＣＯ．ＪＰ", "ご家族", "1回払い", "", "'22/08", "3844", "3844", "", "", "", "", ""},
					{"2022/7/13", "東京都水道局", "ご家族", "1回払い", "", "'22/08", "8459", "8459", "", "", "", "", ""},
					{"2022/7/16", "セブン－イレブン／ｉＤ", "ご家族", "1回払い", "", "'22/08", "98", "98", "", "", "", "", ""},
					{"2022/7/19", "メルカリ", "ご家族", "1回払い", "", "'22/08", "2700", "2700", "", "", "", "", ""},
				},
			},
			want: []*srcRecord{
				{Date: "2022/7/4", ShopName: "ＡＭＡＺＯＮ．ＣＯ．ＪＰ", Amount: 3844},
				{Date: "2022/7/13", ShopName: "東京都水道局", Amount: 8459},
				{Date: "2022/7/16", ShopName: "セブン－イレブン／ｉＤ", Amount: 98},
				{Date: "2022/7/19", ShopName: "メルカリ", Amount: 2700},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNonFixedSrcRecords(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNonFixedSrcRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNonFixedSrcRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeOutputRecords(t *testing.T) {
	type args struct {
		outputRecords []*outputRecord
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "Can write as expected",
			args: args{
				outputRecords: []*outputRecord{
					{Date: "2022/7/4", Item: "ＡＭＡＺＯＮ．ＣＯ．ＪＰ", Amount: 3844},
					{Date: "2022/7/13", Item: "東京都水道局", Amount: 8459},
					{Date: "2022/7/16", Item: "セブン－イレブン／ｉＤ", Amount: 98},
					{Date: "2022/7/19", Item: "メルカリ", Amount: 2700},
				},
			},
			wantWriter: `Date,Item,Amount,Purpose,Method
2022/7/4,ＡＭＡＺＯＮ．ＣＯ．ＪＰ,3844,,
2022/7/13,東京都水道局,8459,,
2022/7/16,セブン－イレブン／ｉＤ,98,,
2022/7/19,メルカリ,2700,,
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := writeOutputRecords(tt.args.outputRecords, writer); (err != nil) != tt.wantErr {
				t.Errorf("writeOutputRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("writeOutputRecords() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
