// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/steverusso/gohttpd"
)

func main() {
	dir, cfg, err := gohttpd.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	h, err := gohttpd.GetHandler(dir, cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if cfg.TLS {
			if err := gohttpd.ServeTLS(h); err != nil {
				log.Fatal(err)
			}
			return
		}

		if err := http.ListenAndServe(cfg.Addr(), h); err != nil {
			log.Fatal(err)
		}
	}()

	// Process signals from the OS.
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM)

	for sig := range sigCh {
		if sig == syscall.SIGTERM {
			os.Exit(0)
		}
	}
}
