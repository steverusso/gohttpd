// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package main

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// HandlerTable maps strings, representing HTTP routes, to `http.Handler`s.
type HandlerTable map[string]http.Handler

// ServeHTTP uses the request URL path as the key for which `http.Handler` to
// use from the table. A NotFound handler is used if there is no matching
// entry.
//
// Requests ending with "/index.html" are redirected to the same URI without
// the "index.html" part.
func (ht HandlerTable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if redirectIndex(w, r) {
		return
	}

	key := path.Clean(r.URL.Path)
	if h, ok := ht[key]; ok {
		h.ServeHTTP(w, r)
		return
	}

	ht.notFound(w, r)
}

// notFound serves the table's 404 handler if present. Otherwise, it writes a
// `NotFound` message and status code to the response.
func (ht HandlerTable) notFound(w http.ResponseWriter, r *http.Request) {
	if h404, ok := ht["/404.html"]; ok {
		h404.ServeHTTP(w, r)
		return
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// redirectIndex checks if the request URL path ends with "index.html". If it
// does, it replies to the request with a redirect to the same URI without the
// "index.html" part and returns true. Otherwise, it just returns false.
func redirectIndex(w http.ResponseWriter, r *http.Request) bool {
	const index = "index.html"

	if strings.HasSuffix(r.URL.Path, index) {
		redirect(w, r, strings.TrimSuffix(r.URL.Path, index))
		return true
	}

	return false
}

// redirect replies to the request with a StatusMovedPermanently redirect to
// the given URL plus the request's query params if present.
func redirect(w http.ResponseWriter, r *http.Request, toURL string) {
	if q := r.URL.RawQuery; q != "" {
		toURL += "?" + q
	}
	http.Redirect(w, r, toURL, http.StatusMovedPermanently)
}

// cachedFile contains the contents and HTTP header values of a file on disk.
type cachedFile struct {
	content []byte
	headers http.Header
}

// ServeHTTP writes the cachedFile's headers and content to the given Response.
func (f *cachedFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range f.headers {
		w.Header()[k] = v
	}
	w.Write(f.content)
}

// newFileMemCache returns an HandlerTable of file paths to cachedFiles. The
// files are read recursively from the given directory and entered as
// cachedFiles in the HandlerTable using their file paths as keys.
func newFileMemCache(dir string) (HandlerTable, error) {
	cache := HandlerTable{}

	return cache, filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		content, err := ioutil.ReadFile(fpath)
		if err != nil {
			return err
		}

		key := strings.TrimPrefix(fpath, dir)
		key = strings.TrimSuffix(key, "index.html")
		key = path.Clean("/" + key)

		cache[key] = &cachedFile{
			content: content,
			headers: http.Header{
				"Content-Type": {mime.TypeByExtension(filepath.Ext(fpath))},
			},
		}

		return nil
	})
}
