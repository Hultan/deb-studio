package engine

import (
	"path"
	"testing"
)

func Test_getFirstLine(t *testing.T) {
	type args struct {
		text   string
		prefix string
		sep    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{"", "PROGRAM=", "\n"}, "", true},
		{"test", args{"test", "PROGRAM=", "\n"}, "", true},
		{"correct", args{"PROGRAM=softtube", "PROGRAM=", "\n"}, "softtube", false},
		{"extra enter", args{"PROGRAM=softtube\n", "PROGRAM=", "\n"}, "softtube", false},
		{"extra text", args{"PROGRAM=softtube\nsdkajsdlkjasd\n", "PROGRAM=", "\n"}, "softtube", false},
		{"extra spaces", args{"PROGRAM=softtube    \nsdkajsdlkjasd\n", "PROGRAM=", "\n"}, "softtube", false},
		{"extra tabs", args{"PROGRAM=softtube 2	\t		\nsdkajsdlkjasd\n", "PROGRAM=", "\n"}, "softtube 2", false},
		{"empty", args{"", "VERSION=", " \t\n"}, "", true},
		{"test", args{"test", "VERSION=", " \t\n"}, "", true},
		{"correct", args{"VERSION=2.6.9", "VERSION=", " \t\n"}, "2.6.9", false},
		{"extra enter", args{"VERSION=2.6.9\n", "VERSION=", " \t\n"}, "2.6.9", false},
		{"extra text", args{"VERSION=2.6.9\nsdkajsdlkjasd\n", "VERSION=", " \t\n"}, "2.6.9", false},
		{"extra spaces", args{"VERSION=2.6.9    \nsdkajsdlkjasd\n", "VERSION=", " \t\n"}, "2.6.9", false},
		{"extra tabs", args{"VERSION=2.6.9			\nsdkajsdlkjasd\n", "VERSION=", " \t\n"}, "2.6.9", false},
		{"empty", args{"", "ARCHITECTURE=", " \t\n"}, "", true},
		{"test", args{"test", "ARCHITECTURE=", " \t\n"}, "", true},
		{"correct", args{"ARCHITECTURE=amd64", "ARCHITECTURE=", " \t\n"}, "amd64", false},
		{"extra enter", args{"ARCHITECTURE=amd64\n", "ARCHITECTURE=", " \t\n"}, "amd64", false},
		{"extra text", args{"ARCHITECTURE=amd64\nsdkajsdlkjasd\n", "ARCHITECTURE=", " \t\n"}, "amd64", false},
		{"extra spaces", args{"ARCHITECTURE=amd64    \nsdkajsdlkjasd\n", "ARCHITECTURE=", " \t\n"}, "amd64", false},
		{"extra tabs", args{"ARCHITECTURE=amd64			\nsdkajsdlkjasd\n", "ARCHITECTURE=", " \t\n"}, "amd64", false},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := getFirstLine(tt.args.text, tt.args.prefix, tt.args.sep)
				if (err != nil) != tt.wantErr {
					t.Errorf("getFirstLine() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("getFirstLine() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_doesDirectoryExist(t *testing.T) {
	type args struct {
		workspacePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"dir exists", args{""}, true},
		{"exists, but file", args{".bashrc"}, false},
		{"dir does not exist", args{".bashrcX"}, false},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				p := path.Join(getUserHomeDirectory(), tt.args.workspacePath)
				if got := doesDirectoryExist(p); got != tt.want {
					t.Errorf("doesDirectoryExist() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
