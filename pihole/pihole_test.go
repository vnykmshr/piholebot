package pihole

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/ChimeraCoder/anaconda"
)

func TestModule_DoTheDew(t *testing.T) {
	m := NewPiHoleBotModule()
	type fields struct {
		cfg     *Config
		client  *http.Client
		twitter *anaconda.TwitterApi
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				cfg:     m.cfg,
				client:  m.client,
				twitter: m.twitter,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				cfg:     tt.fields.cfg,
				client:  tt.fields.client,
				twitter: tt.fields.twitter,
			}
			if err := m.DoTheDew(); (err != nil) != tt.wantErr {
				t.Errorf("Module.DoTheDew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModule_fetch(t *testing.T) {
	m := NewPiHoleBotModule()
	m.cfg.Server.Host = "http://test"
	type fields struct {
		cfg     *Config
		client  *http.Client
		twitter *anaconda.TwitterApi
	}
	tests := []struct {
		name    string
		fields  fields
		want    Stats
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				cfg:     m.cfg,
				client:  m.client,
				twitter: m.twitter,
			},
			want:    Stats{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				cfg:     tt.fields.cfg,
				client:  tt.fields.client,
				twitter: tt.fields.twitter,
			}
			got, err := m.fetch()
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_compose(t *testing.T) {
	m := NewPiHoleBotModule()
	type fields struct {
		cfg     *Config
		client  *http.Client
		twitter *anaconda.TwitterApi
	}
	type args struct {
		stats Stats
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "t1",
			fields: fields{
				cfg:     m.cfg,
				client:  m.client,
				twitter: m.twitter,
			},
			args: args{
				stats: Stats{},
			},
			want: "",
		},
		{
			name: "t1",
			fields: fields{
				cfg:     m.cfg,
				client:  m.client,
				twitter: m.twitter,
			},
			args: args{
				stats: Stats{
					AdsBlockedToday: 10,
				},
			},
			want: "Today, I have blocked 10 queries processing 0 DNS requests from 0 clients. Ads blocked: 0.00%, Blocklist: 0 #pihole",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				cfg:     tt.fields.cfg,
				client:  tt.fields.client,
				twitter: tt.fields.twitter,
			}
			if got := m.compose(tt.args.stats); got != tt.want {
				t.Errorf("Module.compose() = %v, want %v", got, tt.want)
			}
		})
	}
}