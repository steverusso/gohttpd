// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"log"
	"net/http"

	"github.com/steverusso/gohttpd"
)

func main() {
	dir, cfg := gohttpd.LoadConfig()

	h, err := gohttpd.GetHandler(dir, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.TLS {
		if err := gohttpd.ServeTLS(h); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := http.ListenAndServe(cfg.Addr(), h); err != nil {
		log.Fatal(err)
	}
}
