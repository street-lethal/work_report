package model

import (
	"reflect"
	"testing"
)

func TestFileToSetting(t *testing.T) {
	type args struct {
		filePath string
	}
	wantedSetting := Setting{
		MonthsAgo: 1,
	}
	wantedSetting.DailyReport.StartsAt = "10:00"
	wantedSetting.DailyReport.RestTime = "00:30"
	tests := []struct {
		name    string
		args    args
		want    Setting
		wantErr bool
	}{
		{
			args: args{
				filePath: "../config/settings_test.json",
			},
			want: wantedSetting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileToSetting(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileToSetting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("FileToSetting() got = %v, want %v", got, tt.want)
			}
		})
	}
}
