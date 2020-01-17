package pihole

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	gcfg "gopkg.in/gcfg.v1"
)

// read gcfg config
func read(cfg interface{}, base string, module string) bool {
	env := os.Getenv("PIENV")
	if env == "" {
		env = "development"
	}

	fname := path.Join(base, fmt.Sprintf("%s.%s.ini", module, env))
	err := gcfg.ReadFileInto(cfg, fname)
	if err == nil {
		log.Println("read config from ", fname)
		return true
	}
	log.Println(err)
	return false
}

// join url
func join(basePath string, paths ...string) string {
	u, err := url.Parse(basePath)
	if err != nil {
		log.Println(err)
		return ""
	}

	p2 := append([]string{u.Path}, paths...)
	result := path.Join(p2...)
	u.Path = result
	return u.String()
}
