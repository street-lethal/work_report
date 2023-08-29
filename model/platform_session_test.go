package model

import (
	"reflect"
	"testing"
)

func TestFileToPlatformSession(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *PlatformSession
		wantErr bool
	}{
		{
			args: args{
				filePath: "../data/platform_session_test.json",
			},
			want: &PlatformSession{
				SessionID: "abc123",
				AWSAuth:   "123abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileToPlatformSession(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileToPlatformSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileToPlatformSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}
