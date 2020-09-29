package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type SiteHandler interface {
	Domains() []string
	http.Handler
}

// getHandler returns a group of Sites if the given dir is a GoHTTPd root
// directory. Otherwise, t returns a single Site from the given directory.
func getSiteHandler(dir string, mem bool) (h SiteHandler, err error) {
	if isRoot(dir) {
		if h, err = loadSites(dir, mem); err != nil {
			log.Fatal(err)
		}
		return
	}
	return NewSite(dir, mem)
}

// isRoot returns true if there is a readable GoHTTPd configuration file
// located in the given directory.
func isRoot(dir string) bool {
	cfgFile := filepath.Join(dir, "gohttpd.yaml")
	fi, err := os.Stat(cfgFile)
	return err == nil && !fi.IsDir()
}

type SiteGroup []*Site

func loadSites(dir string, mem bool) (SiteGroup, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var sg SiteGroup
	for _, f := range files {
		if f.IsDir() {
			fpath := filepath.Join(dir, f.Name())

			s, err := NewSite(fpath, mem)
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

func NewSite(dir string, mem bool) (_ *Site, err error) {
	var h http.Handler
	if mem {
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
