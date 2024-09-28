package service

import (
	"reflect"
	"testing"
	"work_report/model"
)

func Test_parseJiraCSVService_Parse(t *testing.T) {
	type args struct {
		csv [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[int]model.Work
		wantErr bool
	}{
		{
			args: args{
				csv: [][]string{
					{"イシューキー", "イシューの概要", "ログに記録された時間", "Logged Seconds", "作業日"},
					{"ticket-1", "", "0.25", "900", "2023-01-01 10:00"},
					{"ticket-1", "", "0.25", "900", "2023-01-02 10:00"},
					{"ticket-1", "", "0.75", "2700", "2023-01-02 10:00"},
					{"ticket-1", "", "0.25", "900", "2023-01-05 10:00"},
					{"ticket-2", "", "0.5", "1800", "2023-01-02 10:00"},
					{"ticket-2", "", "0.5", "1800", "2023-01-03 10:00"},
					{"ticket-3", "", "1.0", "3600", "2023-01-05 10:00"},
					{"ticket-3", "", "1.0", "3600", "2023-01-07 10:00"},
				},
			},
			want: map[int]model.Work{
				1: {
					Contents: []string{"ticket-1"},
					Hours:    0.25,
				},
				2: {
					Contents: []string{"ticket-1", "ticket-2"},
					Hours:    1.5,
				},
				3: {
					Contents: []string{"ticket-2"},
					Hours:    0.5,
				},
				5: {
					Contents: []string{"ticket-1", "ticket-3"},
					Hours:    1.25,
				},
				7: {
					Contents: []string{"ticket-3"},
					Hours:    1.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := parseJiraCSVService{}
			got, err := ss.Parse(tt.args.csv)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
