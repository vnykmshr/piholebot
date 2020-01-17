package main

import (
	"log"

	"github.com/vnykmshr/piholebot/pihole"
)

func main() {
	err := pihole.NewPiHoleBotModule().DoTheDew()
	if err != nil {
		log.Fatalln(err)
	}
}
