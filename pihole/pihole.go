package pihole

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// DoTheDew fetch stats, compose and post tweet
func (m *Module) DoTheDew() error {
	log.Printf("[%s][%s] Doing the dew!", m.Config.Server.Name, m.Version)
	stats, err := m.fetch()
	if err != nil {
		return wrap("failed to fetch", err)
	}

	if m.Version == "decompose" {
		stats.AdsBlockedToday = 0
	}

	msg := m.compose(stats)
	if msg == "" {
		return wrap("failed to compose", errors.New("empty message"))
	}

	if !m.Config.Twitter.Enabled {
		log.Printf("[%s][%s][dry] %s", m.Config.Server.Name, m.Version, msg)
		return nil
	}

	tweet, err := m.twitter.PostTweet(msg, nil)
	if err != nil {
		return wrap("failed to tweet", err)
	}

	log.Printf("[%s][%s][%s] %s", m.Config.Server.Name, m.Version, tweet.CreatedAt, tweet.Text)
	return nil
}

func (m *Module) fetch() (Stats, error) {
	var data Stats
	resp, err := m.client.Get(join(m.Config.Server.Host, "admin", "api.php"))
	if err != nil {
		return data, wrap("failed to request", err)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, wrap("failed to unmarshal", err)
	}

	return data, nil
}

func (m *Module) compose(stats Stats) string {
	if stats.AdsBlockedToday == 0 {
		return ""
	}

	var template = `Today, I have blocked %d queries processing %d DNS requests from %d clients. Ads blocked: %s%%, Blocklist: %d #pihole`
	return fmt.Sprintf(template, stats.AdsBlockedToday, stats.DNSQueriesToday, stats.UniqueClients, fmt.Sprintf("%.2f", stats.AdsPercentageToday), stats.DomainsBeingBlocked)
}
