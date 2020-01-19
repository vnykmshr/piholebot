package pihole

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ChimeraCoder/anaconda"
)

func TestModule_DoTheDew(t *testing.T) {
	m := NewPiHoleBotModule("test")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"domains_being_blocked":114789,"dns_queries_today":38588,"ads_blocked_today":5926,"ads_percentage_today":15.357106,"unique_domains":1492,"queries_forwarded":30731,"queries_cached":1931,"clients_ever_seen":12,"unique_clients":10,"dns_queries_all_types":38588,"reply_NODATA":260,"reply_NXDOMAIN":1751,"reply_CNAME":3054,"reply_IP":7320,"privacy_level":0,"status":"enabled","gravity_last_updated":{"file_exists":true,"absolute":1578780798,"relative":{"days":"5","hours":"12","minutes":"26"}}}`))
	})

	httpClient, teardown := testingHTTPClient(h)
	m.client = httpClient
	defer teardown()

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
	m := NewPiHoleBotModule("test")
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
	m := NewPiHoleBotModule("test")
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

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
  s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}
