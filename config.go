// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package gohttpd

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	TLS  bool
	Port int
	Log  string
}

func LoadConfig() (dir string, cfg Config) {
	flag.BoolVar(&cfg.TLS, "tls", false, "Use TLS.")
	flag.IntVar(&cfg.Port, "port", 8080, "The port to use for the local server.")
	flag.StringVar(&cfg.Log, "log", "", "The type of logging.")
	flag.Usage = usage
	flag.Parse()

	dir = "."
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}
	return
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

func usage() {
	fmt.Fprint(os.Stderr, usageMsg)
	flag.PrintDefaults()
	os.Exit(2)
}

var usageMsg = `GoHTTPd is a specific, stubbornly simple web server.

Usage:

  gohttpd [options] <directory>

Options:

`
