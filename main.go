package main

import (
	"log"
	"time"

	"github.com/jpillora/backoff"
	"github.com/vnykmshr/piholebot/pihole"
)

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
		err = pihole.NewPiHoleBotModule().DoTheDew()
		if err != nil {
			if b.Attempt() >= MaxAttempts {
				break
			}

			d := b.Duration()
			log.Printf("[pihole][main] error: %s, retrying in %s", err, d)
			time.Sleep(d)
			continue
		}

		b.Reset()
		break
	}

	if err != nil {
		log.Fatalf("persistent error: %s, retries %f exhausted", err, MaxAttempts)
	}
}
