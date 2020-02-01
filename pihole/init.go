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
	Auth    string
	Timeout time.Duration

	MaxAttempts int32
	MinDelay    time.Duration
	MaxDelay    time.Duration
	Factor      int32
	Jitter      bool
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
	Version string
	Config  *Config
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

// TopStats top stats
type TopStats struct {
	TopQueries map[string]int32 `json:"top_queries,omitempty"`
	TopAds     map[string]int32 `json:"top_ads,omitempty"`
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
		Version: version,
		Config:  &cfg,
		client: &http.Client{
			Timeout: cfg.Server.Timeout * time.Second,
		},
		twitter: anaconda.NewTwitterApiWithCredentials(cfg.Twitter.AccessToken, cfg.Twitter.AccessTokenSecret, cfg.Twitter.ConsumerKey, cfg.Twitter.ConsumerSecret),
	}
}

func getDefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Name:        "piholebot",
			Host:        "http://localhost",
			Auth:        "testing",
			Timeout:     1,
			MaxAttempts: 5,
			MinDelay:    1,
			MaxDelay:    10,
			Factor:      2,
			Jitter:      true,
		},
		Twitter: TwitterConfig{
			Enabled: false,
		},
	}
}
