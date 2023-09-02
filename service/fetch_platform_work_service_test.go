package service

import (
	"reflect"
	"testing"
	"time"
	"work_report/model"
)

func Test_fetchPlatformWorkService_AttachDailyIDsToReport(t *testing.T) {
	year, month, _ := time.Now().Date()
	months := (year-2023)*12 + (int(month) - 1)
	type fields struct {
		setting          model.Setting
		ParseHTMLService ParseHTMLService
	}
	type args struct {
		ids    map[int]string
		report model.Report
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.Report
	}{
		{
			fields: fields{
				setting: model.Setting{
					MonthsAgo: months,
				},
				ParseHTMLService: NewParseHTMLService(),
			},
			args: args{
				ids: map[int]string{
					1: "a",
					3: "c",
				},
				report: model.Report{
					Data: model.WorkData{
						DailyReport: model.DayToDailyData{
							"20230101": model.DailyData{
								ID:         "---",
								TargetDate: "2023-01-01",
							},
							"20230102": model.DailyData{
								ID:         "---",
								TargetDate: "2023-01-02",
							},
							"20230103": model.DailyData{
								ID:         "---",
								TargetDate: "2023-01-03",
							},
							"20230104": model.DailyData{
								ID:         "---",
								TargetDate: "2023-01-04",
							},
						},
					},
				},
			},
			want: model.Report{
				Data: model.WorkData{
					DailyReport: model.DayToDailyData{
						"20230101": model.DailyData{
							ID:         "a",
							TargetDate: "2023-01-01",
						},
						"20230102": model.DailyData{
							ID:         "---",
							TargetDate: "2023-01-02",
						},
						"20230103": model.DailyData{
							ID:         "c",
							TargetDate: "2023-01-03",
						},
						"20230104": model.DailyData{
							ID:         "---",
							TargetDate: "2023-01-04",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fetchPlatformWorkService{
				setting:          tt.fields.setting,
				ParseHTMLService: tt.fields.ParseHTMLService,
			}
			got := s.AttachDailyIDsToReport(tt.args.ids, tt.args.report)
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("AttachDailyIDsToReport() = %v, want %v", *got, tt.want)
			}
		})
	}
}
