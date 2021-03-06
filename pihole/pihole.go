package pihole

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
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

	topStats, err := m.fetchTopStats()
	if err != nil {
		// ignore top Stats
		log.Printf("no top stats, skipped: %s", err)
		topStats = TopStats{}
	}

	msg := m.compose(stats, topStats)
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

func (m *Module) fetchTopStats() (TopStats, error) {
	var data TopStats
	u, err := url.Parse(join(m.Config.Server.Host, "admin", "api.php"))
	if err != nil {
		return data, err
	}

	q := u.Query()
	q.Set("topItems", "1")
	q.Set("auth", m.Config.Server.Auth)
	u.RawQuery = q.Encode()

	resp, err := m.client.Get(u.String())
	if err != nil {
		return data, wrap("failed to request", err)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, wrap("failed to unmarshal", err)
	}

	return data, nil
}

func (m *Module) compose(stats Stats, topStats TopStats) string {
	if stats.AdsBlockedToday == 0 {
		return ""
	}

	var bstr []string
	var i = 1
	for d, c := range topStats.TopAds {
		bstr = append(bstr, fmt.Sprintf("#%d Blocked: %s(%d)", i, d, c))
		i++
	}

	var topString = ""
	if len(topStats.TopAds) > 0 {
		topString = bstr[0]
	}

	var template = `Today, I have blocked %d queries processing %d DNS requests from %d clients. Ads blocked: %s%%, Blocklist: %d %s #pihole`
	return fmt.Sprintf(template, stats.AdsBlockedToday, stats.DNSQueriesToday, stats.UniqueClients, fmt.Sprintf("%.2f", stats.AdsPercentageToday), stats.DomainsBeingBlocked, topString)
}
