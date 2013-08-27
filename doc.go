// Package web provides various web-serving utilities and a server multiplexer
// which integrates automatic reverse proxy and can serve multiple sites on
// different ports simultaneously.
//
// Example usage:
//
//		package main
//
//		import (
//			"github.com/SlyMarbo/web"
//			"io"
//			"net/http"
//			"os"
//		)
//
//		// Simple 404 handler.
//		func notFound(w http.ResponseWriter, r *http.Request) {
//			http.NotFound(w, r)
//		}
//
//		// HTML-specific handler.
//		func serveHTML(w http.ResponseWriter, r *http.Request) {
//			w.Header().Set("Content-Type", "text/html; charset=utf-8")
//			serveContent(w, r)
//		}
//
//		// Solid handler for generic content.
//		func serveContent(w http.ResponseWriter, r *http.Request) {
//			// Add GZIP-encoding if supported by the client.
//			writer := web.NewGzipResponseWriter(w)
//			defer writer.Close()
//
//			f, err := os.Open("." + r.RequestURI)
//			if err != nil {
//				notFound(w, r)
//			}
//			stat, err := f.Stat()
//			if err != nil {
//				notFound(w, r)
//			}
//
//			web.Cache(writer, stat.ModTime(), web.OneYear)
//
//			_, err = io.Copy(writer, f)
//			if err != nil {
//				notFound(w, r)
//			}
//		}
//
//		func main() {
//			redirector := web.NewSite("example.com", 80, notFound)
//			redirector.Always(web.RedirectToHttpsHandler)
//
//			site := web.NewSite("example.com", 443, notFound)
//			site.Equals(web.Handler(serveHTML), "/", "/index.html")
//			site.HasPrefix(web.Handler(serveContent), "/")
//
//			err := new(web.Server).Add(redirector).Add(site).Serve()
//			if err != nil {
//				panic(err)
//			}
//		}
package web
