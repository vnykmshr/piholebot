package pihole

import (
	"reflect"
	"testing"
)

func TestNewPiHoleBotModule(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "t1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if m := NewPiHoleBotModule("test"); (m == nil) != tt.wantErr {
				t.Errorf("NewPiHoleBotModule() m = %v, wantErr %v", m, tt.wantErr)
			}
		})
	}
}

func Test_getDefaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{
			name: "t1",
			want: Config{
				Server: ServerConfig{
					Name:    "piholebot",
					Host:    "http://localhost",
					Timeout: 1,
				},
				Twitter: TwitterConfig{
					Enabled: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDefaultConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefaultConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
