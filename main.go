package main

import (
	"log"
	"time"

	"github.com/jpillora/backoff"
	"github.com/vnykmshr/piholebot/pihole"
)

var version string

func main() {
	MaxAttempts := 5.0

	b := &backoff.Backoff{
		Min:    1 * time.Second,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	var err error
	for {
		err = pihole.NewPiHoleBotModule(version).DoTheDew()
		if err != nil {
			if b.Attempt() >= MaxAttempts {
				break
			}

			d := b.Duration()
			log.Printf("[pihole][main][%s] error: %s, retrying in %s", err, d, version)
			time.Sleep(d)
			continue
		}

		b.Reset()
		break
	}

	if err != nil {
		log.Fatalf("[pihole][main][%s] persistent error: %s, retries %f exhausted", version, err, MaxAttempts)
	}
}
