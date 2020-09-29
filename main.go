// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	tls := flag.Bool("tls", false, "Use TLS.")
	mem := flag.Bool("mem", false, "Cache files in memory instead of using disk.")
	port := flag.Int("port", 8080, "The port to use for the local server.")
	flag.Usage = usage
	flag.Parse()

	dir := "/var/www"
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}

	h, err := getSiteHandler(dir, *mem)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if !*tls {
			if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), h); err != nil {
				log.Fatal(err)
			}
			return
		}

		if err := ServeTLS(h); err != nil {
			log.Fatal(err)
		}
	}()

	// Block here to catch and process signals from the OS.
	listenForSignals(func(sig os.Signal) {
		if sig == syscall.SIGTERM {
			os.Exit(0)
		}
	})
}

// listenForSignals indicates which signals we want to catch and runs the
// given function on each signal that is read from the os.Signal channel.
func listenForSignals(sigFn func(os.Signal)) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM)

	for sig := range sigCh {
		sigFn(sig)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, usageMsg)
	flag.PrintDefaults()
	os.Exit(2)
}

var usageMsg = `GoHTTPd is a specific, stubbornly simple web server.

Usage:

  gohttpd <directory> [options]

Options:

`
