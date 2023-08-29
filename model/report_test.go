package model

import (
	"reflect"
	"testing"
)

func TestFileToReport(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Report
		wantErr bool
	}{
		{
			args: args{
				filePath: "../data/report_test.json",
			},
			want: &Report{
				Data: WorkData{
					DailyReport: DayToDailyData{
						"20230104": {
							TargetDate:  "2023-01-04",
							StartTime:   "10:00",
							EndTime:     "18:00",
							RelaxTime:   "00:30",
							WorkContent: "test1",
						},
						"20230105": {
							TargetDate:  "2023-01-05",
							StartTime:   "10:30",
							EndTime:     "18:30",
							RelaxTime:   "01:00",
							WorkContent: "test2",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileToReport(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileToReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileToReport() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReport_ToQuery(t *testing.T) {
	type fields struct {
		Data WorkData
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				Data: WorkData{
					DailyReport: DayToDailyData{
						"01": DailyData{
							ID:          "12345",
							TargetDate:  "0101",
							StartTime:   "10:00",
							EndTime:     "18:00",
							RelaxTime:   "00:30",
							WorkTime:    "07:30",
							WorkContent: "a, b, c",
						},
					},
				},
			},
			want: `data%5BDailyReport%5D%5B01%5D%5Bend_time%5D=18%3A00&` +
				`data%5BDailyReport%5D%5B01%5D%5Bid%5D=12345&` +
				`data%5BDailyReport%5D%5B01%5D%5Brelax_time%5D=00%3A30&` +
				`data%5BDailyReport%5D%5B01%5D%5Bstart_time%5D=10%3A00&` +
				`data%5BDailyReport%5D%5B01%5D%5Btarget_date%5D=0101&` +
				`data%5BDailyReport%5D%5B01%5D%5Bwork_content%5D=a%2C+b%2C+c&` +
				`data%5BDailyReport%5D%5B01%5D%5Bwork_time%5D=07%3A30`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Report{
				Data: tt.fields.Data,
			}
			got, err := r.ToQuery()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDailyData_CalcWorkTime(t *testing.T) {
	type fields struct {
		StartTime string
		EndTime   string
		RelaxTime string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				StartTime: "10:00",
				EndTime:   "18:00",
				RelaxTime: "00:30",
			},
			want: "07:30",
		},
		{
			fields: fields{
				StartTime: "09:30",
				EndTime:   "18:00",
				RelaxTime: "00:30",
			},
			want: "08:00",
		},
		{
			fields: fields{
				StartTime: "09:30",
				EndTime:   "20:15",
				RelaxTime: "01:30",
			},
			want: "09:15",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DailyData{
				StartTime: tt.fields.StartTime,
				EndTime:   tt.fields.EndTime,
				RelaxTime: tt.fields.RelaxTime,
			}
			if err := d.CalcWorkTime(); (err != nil) != tt.wantErr {
				t.Errorf("CalcWorkTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if d.WorkTime != tt.want {
				t.Errorf("CalcWorkTime() workTime = %v, want %v", d.WorkTime, tt.want)
			}
		})
	}
}
