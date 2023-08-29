package model

import (
	"reflect"
	"testing"
)

func TestFileToPlatformID(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *PlatformID
		wantErr bool
	}{
		{
			args: args{
				filePath: "../config/platform_id_test.json",
			},
			want: &PlatformID{
				Email:    "test@example.com",
				Password: "p@ssw0rd",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileToPlatformID(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileToPlatformID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileToPlatformID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
