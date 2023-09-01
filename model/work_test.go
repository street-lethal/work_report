package model

import (
	"reflect"
	"testing"
)

func TestWork_AddContent(t *testing.T) {
	type fields struct {
		Contents []string
		Hours    float64
	}
	type args struct {
		content string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			fields: fields{
				Contents: []string{"one", "two"},
			},
			args: args{
				content: "three",
			},
			want: []string{"one", "two", "three"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Work{
				Contents: tt.fields.Contents,
				Hours:    tt.fields.Hours,
			}
			w.AddContent(tt.args.content)
			if !reflect.DeepEqual(w.Contents, tt.want) {
				t.Errorf("AddContent() Contents = %v, want %v", w.Contents, tt.want)
			}
		})
	}
}

func TestWork_AddHour(t *testing.T) {
	type fields struct {
		Contents []string
		Hours    float64
	}
	type args struct {
		hour float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			fields: fields{
				Hours: 4.75,
			},
			args: args{
				hour: 1.5,
			},
			want: 6.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Work{
				Contents: tt.fields.Contents,
				Hours:    tt.fields.Hours,
			}
			w.AddHour(tt.args.hour)
			if !reflect.DeepEqual(w.Hours, tt.want) {
				t.Errorf("AddContent() Hours = %v, want %v", w.Hours, tt.want)
			}
		})
	}
}

func TestWork_HourMin(t *testing.T) {
	type fields struct {
		Contents []string
		Hours    float64
	}
	tests := []struct {
		name       string
		fields     fields
		wantHour   int
		wantMinute int
	}{
		{
			fields: fields{
				Hours: 7.75,
			},
			wantHour:   7,
			wantMinute: 45,
		},
		{
			fields: fields{
				Hours: 8.0,
			},
			wantHour:   8,
			wantMinute: 0,
		},
		{
			fields: fields{
				Hours: 9.25,
			},
			wantHour:   9,
			wantMinute: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Work{
				Contents: tt.fields.Contents,
				Hours:    tt.fields.Hours,
			}
			gotHour, gotMinute := w.HourMin()
			if gotHour != tt.wantHour {
				t.Errorf("HourMin() gotHour = %v, want %v", gotHour, tt.wantHour)
			}
			if gotMinute != tt.wantMinute {
				t.Errorf("HourMin() gotMinute = %v, want %v", gotMinute, tt.wantMinute)
			}
		})
	}
}
