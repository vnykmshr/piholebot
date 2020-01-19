package pihole

import (
	"log"
	"net/http"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

// Config PiHoleBotModule config
type Config struct {
	Server  ServerConfig
	Twitter TwitterConfig
}

// ServerConfig pihole server config
type ServerConfig struct {
	Name    string
	Host    string
	Timeout time.Duration
}

// TwitterConfig bot twitter config
type TwitterConfig struct {
	Enabled           bool
	Username          string
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// Module PiHoleBotModule
type Module struct {
	version string
	cfg     *Config
	client  *http.Client
	twitter *anaconda.TwitterApi
}

// Stats PiHole response stats
type Stats struct {
	DomainsBeingBlocked int32   `json:"domains_being_blocked,omitempty"`
	DNSQueriesToday     int32   `json:"dns_queries_today,omitempty"`
	AdsBlockedToday     int32   `json:"ads_blocked_today,omitempty"`
	AdsPercentageToday  float32 `json:"ads_percentage_today,omitempty"`
	UniqueClients       int32   `json:"unique_clients,omitempty"`
}

// NewPiHoleBotModule new piholebot module
func NewPiHoleBotModule(version string) *Module {
	var cfg Config
	ok := read(&cfg, "/etc/piholebot", "piholebot") || read(&cfg, "files/etc/piholebot", "piholebot")
	if !ok {
		log.Println("failed to read config, loading defaults")
		cfg = getDefaultConfig()
	}

	return &Module{
		version: version,
		cfg:     &cfg,
		client: &http.Client{
			Timeout: cfg.Server.Timeout * time.Second,
		},
		twitter: anaconda.NewTwitterApiWithCredentials(cfg.Twitter.AccessToken, cfg.Twitter.AccessTokenSecret, cfg.Twitter.ConsumerKey, cfg.Twitter.ConsumerSecret),
	}
}

func getDefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Name:    "piholebot",
			Host:    "http://localhost",
			Timeout: 1,
		},
		Twitter: TwitterConfig{
			Enabled: false,
		},
	}
}
