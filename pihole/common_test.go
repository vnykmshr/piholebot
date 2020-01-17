package pihole

import (
	"testing"
)

func Test_read(t *testing.T) {
	type args struct {
		cfg    interface{}
		base   string
		module string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "t1",
			args: args{
				cfg:    nil,
				base:   "files/etc/piholebot",
				module: "invalid",
			},
			want: false,
		},
		{name: "t2",
			args: args{
				cfg:    &Config{},
				base:   "../files/etc/piholebot",
				module: "piholebot",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := read(tt.args.cfg, tt.args.base, tt.args.module); got != tt.want {
				t.Errorf("read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_join(t *testing.T) {
	type args struct {
		basePath string
		paths    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				basePath: "http://test",
				paths:    []string{"d1", "d2"},
			},
			want: "http://test/d1/d2",
		},
		{
			name: "t1",
			args: args{
				basePath: "http://test/",
				paths:    []string{"/d1/", "/d2/"},
			},
			want: "http://test/d1/d2",
		},
		{
			name: "t1",
			args: args{
				basePath: ":",
				paths:    []string{"/d1/", "/d2/"},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := join(tt.args.basePath, tt.args.paths...); got != tt.want {
				t.Errorf("join() = %v, want %v", got, tt.want)
			}
		})
	}
}
