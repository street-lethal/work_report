package service

import (
	"reflect"
	"testing"
)

func Test_parseCSVService_Parse(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			args: args{
				filePath: "../data/csv_test.csv",
			},
			want: [][]string{
				{"id", "name"},
				{"101", "A"},
				{"102", "B"},
				{"105", "X"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := parseCSVService{}
			got, err := s.Parse(tt.args.filePath)
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
