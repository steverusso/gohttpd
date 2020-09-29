// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package gohttpd

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	TLS  bool
	Mem  bool
	Port int
}

func LoadConfig() (string, Config, error) {
	tls := flag.Bool("tls", false, "Use TLS.")
	mem := flag.Bool("mem", false, "Cache files in memory instead of using disk.")
	port := flag.Int("port", 8080, "The port to use for the local server.")
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) <= 0 {
		return "", Config{}, errors.New("error: Must provide a directory.")
	}

	dir := flag.Args()[0]

	return dir, Config{
		TLS:  *tls,
		Mem:  *mem,
		Port: *port,
	}, nil
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
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