package service

import (
	"os"
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func Test_parseHTMLService_Parse(t *testing.T) {
	bin, err := os.ReadFile("../data/html_test.html")
	if err != nil {
		t.Errorf(err.Error())
	}
	htm := string(bin)
	type args struct {
		htm string
	}
	tests := []struct {
		name    string
		args    args
		want    *html.Node
		wantErr bool
	}{
		{
			args: args{
				htm: htm,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := parseHTMLService{}
			got, err := s.Parse(tt.args.htm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got = got.FirstChild
			if got.Data == "" {
				t.Errorf("got.Data is empty")
			}
		})
	}
}

func Test_parseHTMLService_HasClass(t *testing.T) {
	htm := `<html class="x top y"><span class="inner">dummy</span></html>`
	type args struct {
		htm  string
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				htm:  htm,
				name: "top",
			},
			want: true,
		},
		{
			args: args{
				htm:  htm,
				name: "inner",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := parseHTMLService{}
			node, err := s.Parse(tt.args.htm)
			if err != nil {
				t.Errorf(err.Error())
			}
			if got := s.HasClass(node.FirstChild, tt.args.name); got != tt.want {
				t.Errorf("HasClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHTMLService_FindFirst(t *testing.T) {
	type args struct {
		className string
	}
	tests := []struct {
		name     string
		args     args
		wantedId string
	}{
		{
			args: args{
				className: "first-class",
			},
			wantedId: "1-1",
		},
		{
			args: args{
				className: "second-class",
			},
			wantedId: "1-2",
		},
		{
			args: args{
				className: "third-class",
			},
			wantedId: "1-3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := parseHTMLService{}
			bin, _ := os.ReadFile("../data/html_test.html")
			node, _ := s.Parse(string(bin))
			got := s.FindFirst(node, func(node *html.Node) bool {
				return s.HasClass(node, tt.args.className)
			})
			for _, attr := range got.Attr {
				if attr.Key != "id" {
					continue
				}

				if attr.Val != tt.wantedId {
					t.Errorf("FindFirst() attr = %v, want %v", attr.Val, tt.wantedId)
				}
			}
		})
	}
}

func Test_parseHTMLService_FindAll_GetAttrValue(t *testing.T) {
	type args struct {
		className string
	}
	tests := []struct {
		name      string
		args      args
		wantedIds []string
	}{
		{
			args: args{
				className: "first-class",
			},
			wantedIds: []string{"1-1", "2-1"},
		},
		{
			args: args{
				className: "second-class",
			},
			wantedIds: []string{"1-2", "2-2"},
		},
		{
			args: args{
				className: "third-class",
			},
			wantedIds: []string{"1-3", "2-3"},
		},
		{
			args: args{
				className: "low",
			},
			wantedIds: []string{"1-2", "1-3", "2-2", "2-3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := parseHTMLService{}
			bin, _ := os.ReadFile("../data/html_test.html")
			node, _ := s.Parse(string(bin))
			got := s.FindAll(node, func(node *html.Node) bool {
				return s.HasClass(node, tt.args.className)
			})
			gotIds := []string{}
			for _, elm := range got {
				gotIds = append(gotIds, s.GetAttrValue(elm, "id"))
			}
			if !reflect.DeepEqual(gotIds, tt.wantedIds) {
				t.Errorf("FindAll() ids = %v, want %v", gotIds, tt.wantedIds)
			}
		})
	}
}
