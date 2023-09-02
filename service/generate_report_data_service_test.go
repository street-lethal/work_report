package service

import (
	"reflect"
	"testing"
	"time"
	"work_report/model"
)

func Test_generateReportDataService_Generate(t *testing.T) {
	now := time.Now()
	monthsDiff := (now.Year()-2023)*12 + int(now.Month()) - 1
	setting := model.Setting{
		MonthsAgo: monthsDiff,
	}
	setting.DailyReport.StartsAt = "10:00"
	setting.DailyReport.RestTime = "00:30"
	type fields struct {
		Setting model.Setting
	}
	type args struct {
		works map[int]model.Work
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.DayToDailyData
	}{
		{
			fields: fields{
				Setting: setting,
			},
			args: args{
				works: map[int]model.Work{
					1: {Contents: []string{"a"}, Hours: 0},
					2: {Contents: []string{"b", "c"}, Hours: 0},
					3: {Contents: []string{"d", "e"}, Hours: 0},
					4: {Contents: []string{"f", "g"}, Hours: 7.5},
					5: {Contents: []string{"h", "i"}, Hours: 8},
					6: {Contents: []string{"j", "k", "l"}, Hours: 7.5},
					7: {Contents: []string{"m", "n"}, Hours: 0},
				},
			},
			want: model.DayToDailyData{
				"20230104": model.DailyData{
					TargetDate:  "2023-01-04",
					StartTime:   "10:00",
					EndTime:     "18:00",
					RelaxTime:   "00:30",
					WorkTime:    "07:30",
					WorkContent: "f, g",
				},
				"20230105": model.DailyData{
					TargetDate:  "2023-01-05",
					StartTime:   "10:00",
					EndTime:     "18:30",
					RelaxTime:   "00:30",
					WorkTime:    "08:00",
					WorkContent: "h, i",
				},
				"20230106": model.DailyData{
					TargetDate:  "2023-01-06",
					StartTime:   "10:00",
					EndTime:     "18:00",
					RelaxTime:   "00:30",
					WorkTime:    "07:30",
					WorkContent: "j, k, l",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := generateReportDataService{
				Setting: tt.fields.Setting,
			}
			got := s.Generate(tt.args.works)
			if len(got.Data.DailyReport) != 3 {
				t.Errorf("Generate() len() = %v, want = 20", len(got.Data.DailyReport))
			}
			if data, ok := got.Data.DailyReport["20230101"]; ok {
				t.Errorf("Generate() 20230101 = %v, want = nil", data)
			}
			if data, ok := got.Data.DailyReport["20230102"]; ok {
				t.Errorf("Generate() 20230102 = %v, want = nil", data)
			}
			if data, ok := got.Data.DailyReport["20230103"]; ok {
				t.Errorf("Generate() 20230103 = %v, want = nil", data)
			}
			if data := got.Data.DailyReport["20230104"]; !reflect.DeepEqual(data, tt.want["20230104"]) {
				t.Errorf("Generate() 20230104 = %v, want = %v", data, tt.want)
			}
			if data := got.Data.DailyReport["20230105"]; !reflect.DeepEqual(data, tt.want["20230105"]) {
				t.Errorf("Generate() 20230105 = %v, want = %v", data, tt.want)
			}
			if data := got.Data.DailyReport["20230106"]; !reflect.DeepEqual(data, tt.want["20230106"]) {
				t.Errorf("Generate() 20230106 = %v, want = %v", data, tt.want)
			}
			if data, ok := got.Data.DailyReport["20230107"]; ok {
				t.Errorf("Generate() 20230107 = %v, want = nil", data)
			}
		})
	}
}
