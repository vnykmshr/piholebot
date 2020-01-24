package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jpillora/backoff"
	"github.com/vnykmshr/piholebot/pihole"
)

var version string
var buildTime string
var versionFlag bool

func init() {
	flag.BoolVar(&versionFlag, "version", false, "binary version")
}

func main() {
	flag.Parse()

	if versionFlag {
		fmt.Println(fmt.Sprintf("PiHoleBot\nPi-Hole Ad-Blocker Tweet Bot\nVersion: %s\nBuild Time: %s", version, buildTime))
		os.Exit(0)
	}

	m := pihole.NewPiHoleBotModule(version)
	b := &backoff.Backoff{
		Min:    m.Config.Server.MinDelay * time.Second,
		Max:    m.Config.Server.MaxDelay * time.Second,
		Factor: float64(m.Config.Server.Factor),
		Jitter: m.Config.Server.Jitter,
	}

	var err error
	for {
		err = m.DoTheDew()
		if err != nil {
			if b.Attempt() >= float64(m.Config.Server.MaxAttempts) {
				log.Fatalf("[pihole][%s] exhausted attempts: %d, last err: %s", version, m.Config.Server.MaxAttempts, err)
			}

			d := b.Duration()
			log.Printf("[pihole][%s] error: %s, retrying in %s", version, err, d)
			time.Sleep(d)
			continue
		}

		// DoTheDew: success
		b.Reset()
		break
	}
}
