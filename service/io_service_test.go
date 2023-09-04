package service

import "testing"

func Test_ioService_Input(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				filePath: "../data/io_test.txt",
			},
			want: "Hello,\nworld!\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ioService{}
			got, err := s.Input(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Input() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Input() got = %v, want %v", got, tt.want)
			}
		})
	}
}
