// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package gohttpd

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Handler interface {
	Domains() []string
	http.Handler
}

// GetHandler returns a group of Sites if the given dir is a GoHTTPd root
// directory. Otherwise, t returns a single Site from the given directory.
func GetHandler(dir string, cfg *Config) (h Handler, err error) {
	if isRoot(dir) {
		if h, err = loadSites(dir, cfg); err != nil {
			return nil, err
		}
		return
	}
	return NewSite(dir, cfg)
}

// isRoot returns true if there is a readable GoHTTPd configuration file
// located in the given directory.
func isRoot(dir string) bool {
	cfgFile := filepath.Join(dir, "gohttpd.yaml")
	fi, err := os.Stat(cfgFile)
	return err == nil && !fi.IsDir()
}

type SiteGroup []*Site

func loadSites(dir string, cfg *Config) (SiteGroup, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var sg SiteGroup
	for _, f := range files {
		if f.IsDir() {
			fpath := filepath.Join(dir, f.Name())

			s, err := NewSite(fpath, cfg)
			if err != nil {
				return nil, err
			}

			sg = append(sg, s)
		}
	}

	return sg, nil
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

type Site struct {
	Domain string
	http.Handler
}

func (s Site) Domains() []string {
	return []string{s.Domain}
}

func NewSite(dir string, cfg *Config) (_ *Site, err error) {
	var h http.Handler
	if cfg.Mem {
		if h, err = newFileMemCache(dir); err != nil {
			return nil, err
		}
	} else {
		h = http.FileServer(http.Dir(dir))
	}

	return &Site{
		Domain:  path.Base(dir),
		Handler: h,
	}, nil
}
