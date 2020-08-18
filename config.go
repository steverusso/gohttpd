// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"fmt"
	"strings"
)

type config struct {
	domainCSV string
	port      int
	mem       bool
}

func (c *config) isProduction() bool {
	return c.domainCSV != ""
}

func (c *config) domains() []string {
	return strings.Split(c.domainCSV, ",")
}

func (c *config) addr() string {
	return fmt.Sprintf(":%d", c.port)
}
