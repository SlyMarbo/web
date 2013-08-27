package web

import (
	"net/http"
	"time"
)

// DoNotCache uses the Cache-Control, Pragma, and Expires HTTP headers
// to advise the client not to cache the response.
func DoNotCache(w http.ResponseWriter) {
	header := w.Header()
	header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	header.Set("Pragma", "no-cache")
	header.Set("Expires", "0")
}

// Cache uses the Last-Modified, Expires, and Vary HTTP headers to
// advise the client to cache the response for the given duration.
func Cache(w http.ResponseWriter, modTime time.Time, duration time.Duration) {
	header := w.Header()
	if !modTime.IsZero() {
		header.Set("Last-Modified", modTime.UTC().Format(http.TimeFormat))
	}
	header.Set("Expires", time.Now().Add(duration).UTC().Format(http.TimeFormat))
	header.Set("Vary", "Accept-Encoding")
}

// May be useful in cache durations.
var OneYear time.Duration = time.Hour * 24 * 366

// RedirectToHTTPS takes an HTTP request and redirects it to the same
// page, but using HTTPS. Be careful not to use in serving HTTPS, or
// an infinite redirection loop will occur.
func RedirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	url.Scheme = "https"
	url.Host = r.Host
	http.Redirect(w, r, url.String(), 301)
}

// RedirectToHttpsHandler can be used as an http.Handler which uses
// RedirectToHTTPS above.
var RedirectToHttpsHandler = Handler(RedirectToHTTPS)

// RedirectToHTTP takes an HTTPS request and redirects it to the same
// page, but using HTTP. Be careful not to use in serving HTTP, or
// an infinite redirection loop will occur.
func RedirectToHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	url.Scheme = "http"
	url.Host = r.Host
	http.Redirect(w, r, url.String(), 301)
}

// RedirectToHttpHandler can be used as an http.Handler which uses
// RedirectToHTTP above.
var RedirectToHttpHandler = Handler(RedirectToHTTP)

// Handler can be used as a shorter http.HandlerFunc.
type Handler func(http.ResponseWriter, *http.Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}
