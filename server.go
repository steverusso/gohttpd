// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

type Server interface {
	Start(h http.Handler) <-chan error
}

// newServer returns the appropriate type of Server based on the config options.
func newServer(c *config) Server {
	if c.isProduction() {
		return &productionServer{c.domains()}
	}
	return &localServer{c.addr()}
}

type productionServer struct {
	domains []string
}

func (srv *productionServer) Start(h http.Handler) <-chan error {
	ch := make(chan error)

	// Production HTTPS server on port 443 (using autocert convenience function).
	go func() {
		if err := http.Serve(autocert.NewListener(srv.domains...), h); err != nil {
			ch <- fmt.Errorf("listening and serving HTTPS: %v", err)
		}
	}()

	// Redirect HTTP requests on port 80 to HTTPS.
	go func() {
		if err := srv.redirectToHTTPS(); err != nil {
			ch <- fmt.Errorf("redirecting to HTTPS: %v", err)
		}
	}()

	return ch
}

// redirectToHTTPS redirects all requests on port 80 to HTTPS requests on port 443.
func (srv *productionServer) redirectToHTTPS() error {
	redir := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+srv.domains[0]+":443"+r.RequestURI, http.StatusMovedPermanently)
	})

	return http.ListenAndServe(":80", redir)
}

type localServer struct {
	addr string
}

func (srv *localServer) Start(h http.Handler) <-chan error {
	ch := make(chan error)

	go func() {
		if err := http.ListenAndServe(srv.addr, h); err != nil {
			ch <- err
		}
	}()

	return ch
}
