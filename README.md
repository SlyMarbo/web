web
===

<<<<<<< HEAD
Example use:
```go
package main

import (
	"github.com/SlyMarbo/web"
	"io"
	"net/http"
	"os"
)

// Simple 404 handler.
func notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

// HTML-specific handler.
func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	serveContent(w, r)
}

// Solid handler for generic content.
func serveContent(w http.ResponseWriter, r *http.Request) {
	// Add GZIP-encoding if supported by the client.
	writer := web.NewGzipResponseWriter(w)
	defer writer.Close()

	// Ensure the file exists.
	f, err := os.Open("." + r.RequestURI)
	if err != nil {
		notFound(w, r)
	}
	stat, err := f.Stat()
	if err != nil {
		notFound(w, r)
	}

	// Add caching.
	web.Cache(writer, stat.ModTime(), web.OneYear)

	// Send the file data. This will be compressed if allowed.
	_, err = io.Copy(writer, f)
	if err != nil {
		notFound(w, r)
	}
}

func main() {
	// Redirect http://example.com requests to https://example.com.
	redirector := web.NewSite("example.com", 80, notFound)
	redirector.HasPrefix(web.RedirectToHttpsHandler, "/")

	site := web.NewSite("example.com", 443, notFound)
	site.Equals(http.HandleFunc(serveHTML), "/", "/index.html")
	site.HasPrefix(http.HandleFunc(serveContent), "/")

	err := new(web.Server).Add(redirector).Add(site).Serve()
	if err != nil {
		panic(err)
	}
}
```
=======
Web is a small helper library in Go for building web servers.
>>>>>>> 801f3d5c633e01fb0e9c1a75d7388c4643e6f264
