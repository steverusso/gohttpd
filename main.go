// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize config from cli args.
	var cfg config
	flag.StringVar(&cfg.root, "root", "", "The root directory of static resources to serve.")
	flag.StringVar(&cfg.domainCSV, "domains", "", "CSV of domains which indicates production.")
	flag.IntVar(&cfg.port, "port", 8080, "The port to use for the local server.")
	flag.BoolVar(&cfg.mem, "mem", false, "Cache files in memory instead of using disk.")
	flag.Parse()

	// Start a new Server and process any errors that are sent back by the Server
	// over the error channel.
	go func() {
		srv := newServer(&cfg)
		for err := range srv.Start(newFileHandler(&cfg)) {
			if err != nil {
				log.Fatal(err)
			}
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
