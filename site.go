// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package gohttpd

import (
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
)

type SiteLoader struct {
	config Config
	logger Logger
}

func NewSiteLoader(config Config) (SiteLoader, error) {
	logger, err := NewLogger(config.Log)
	if err != nil {
		return SiteLoader{}, err
	}

	return SiteLoader{config, logger}, nil
}

// GetHandler returns a group of Sites if the given dir is a GoHTTPd root
// directory. Otherwise, t returns a single Site from the given directory.
func (sl *SiteLoader) Load(dir string) (h http.Handler, err error) {
	if sl.config.TLS {
		if h, err = sl.loadSites(dir); err != nil {
			return nil, err
		}
		return
	}
	return sl.loadSingleSite(dir), nil
}

func (sl *SiteLoader) loadSites(dir string) (SiteGroup, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var sg SiteGroup
	for _, f := range files {
		if f.IsDir() {
			fpath := filepath.Join(dir, f.Name())
			sg = append(sg, sl.loadSingleSite(fpath))
		}
	}

	return sg, nil
}

type Site struct {
	Domain string
	http.Handler
}

type SiteGroup []*Site

var ignoreRegex = regexp.MustCompile(`\.(png|jpg|css|ico)+$`)

func (sl *SiteLoader) loadSingleSite(dir string) *Site {
	h := http.FileServer(http.Dir(dir))
	if sl.logger != nil {
		h = chain(h, func(_ http.ResponseWriter, r *http.Request) {
			if !ignoreRegex.MatchString(r.URL.String()) {
				sl.logger.Log(r)
			}
		})
	}

	return &Site{
		Domain:  path.Base(dir),
		Handler: h,
	}
}

func (sg SiteGroup) Domains() (d []string) {
	for _, s := range sg {
		d = append(d, s.Domain)
	}
	return
}

func (sg SiteGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, s := range sg {
		if s.Domain == r.Host {
			s.ServeHTTP(w, r)
			return
		}
	}
	// uh oh
}

func chain(h http.Handler, hf http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hf(w, r)
		h.ServeHTTP(w, r)
	})
}
