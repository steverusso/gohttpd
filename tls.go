// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package gohttpd

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func ServeTLS(h http.Handler) error {
	sg, ok := h.(SiteGroup)
	if !ok {
		return errors.New("can only serve SiteGroup over TLS")
	}

	errCh := make(chan error)

	// Serve the `sites` with TLS.
	go func() {
		if err := http.Serve(autocert.NewListener(sg.Domains()...), sg); err != nil {
			errCh <- err
		}
	}()

	// Redirect all requests on port 80 to HTTPS on port 443.
	go func() {
		redir := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+":443"+r.RequestURI, http.StatusMovedPermanently)
		})

		if err := http.ListenAndServe(":80", redir); err != nil {
			errCh <- err
		}
	}()

	return <-errCh
}
