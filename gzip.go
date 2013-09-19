// Copyright 2013 Jamie Hall. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"compress/gzip"
	"github.com/SlyMarbo/spdy"
	"net/http"
	"strings"
)

// Gzip determines whether either the given request headers claim support
// for GZIP-encoded responses, or the connection is using SPDY, in which
// case GZIP support is guaranteed.
func Gzip(w http.ResponseWriter, r *http.Request) bool {
	return spdy.UsingSPDY(w) || strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

// GzipResponseWriter is used to provide GZIP encoding when permitted
// by the request headers. It fulfils the http.ResponseWriter interface
// so can be used completely in place of the provided ResponseWriter.
type GzipResponseWriter struct {
	w http.ResponseWriter
	g *gzip.Writer
}

// NewGzipResponseWriter takes a request and response writer and forms
// a GzipResponseWriter.
func NewGzipResponseWriter(w http.ResponseWriter, r *http.Request) *GzipResponseWriter {
	gzw := new(GzipResponseWriter)
	gzw.w = w
	if Gzip(w, r) {
		w.Header().Set("Content-Encoding", "gzip")
		gzw.g = gzip.NewWriter(w)
	}

	return gzw
}

// NewGzipResponseWriterLevel takes a request and response writer and forms
// a GzipResponseWriter using the given compression level if using GZIP.
func NewGzipResponseWriterLevel(w http.ResponseWriter, r *http.Request, level int) (*GzipResponseWriter, error) {
	gzw := new(GzipResponseWriter)
	gzw.w = w
	if Gzip(w, r) {
		w.Header().Set("Content-Encoding", "gzip")
		g, err := gzip.NewWriterLevel(w, level)
		if err != nil {
			return nil, err
		}
		gzw.g = g
	}
	return gzw, nil
}

// Close closes the underlying gzip.Writer if necessary. This must
// be called to ensure written data is flushed once all writing has
// been completed.
func (g *GzipResponseWriter) Close() error {
	if g.g != nil {
		return g.g.Close()
	}
	return nil
}

// Flush flushes the underlying gzip.Writer if necessary.
func (g *GzipResponseWriter) Flush() error {
	if g.g != nil {
		return g.g.Flush()
	}
	return nil
}

// Header returns the header map that will be sent by WriteHeader.
// Changing the header after a call to WriteHeader (or Write) has
// no effect.
func (g *GzipResponseWriter) Header() http.Header {
	return g.w.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (g *GzipResponseWriter) Write(data []byte) (int, error) {
	if g.g != nil {
		return g.g.Write(data)
	}
	return g.w.Write(data)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (g *GzipResponseWriter) WriteHeader(status int) {
	g.w.WriteHeader(status)
}
