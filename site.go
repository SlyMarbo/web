package web

import (
	"net/http"
	"strings"
)

// Site manages the handling of requests for a particular site.
// Site's methods are used to register handlers for different
// request paths. The handlers are tried in order, so if a given
// path would match multiple handlers, the first matching handler
// registered is used.
//
// Site must be created with NewSite or NewSecureSite.
type Site struct {
	Name     string
	Port     int
	auth     []string
	handlers []*Matcher
	notFound http.HandlerFunc
}

// NewSite builds a new HTTP Site, using the given domain name
// and port number. The provided handler is called when a request
// path does not match any handlers. If nil, http.NotFoundHandler
// is used instead.
//
//		// http://example.com
//		site := NewSite("example.com", 80, nil)
//
//		// http://example.com:8080
//		site := NewSite("example.com", 8080, nil)
func NewSite(name string, port int, notFound http.HandlerFunc) *Site {
	return &Site{
		Name:     name,
		Port:     port,
		handlers: make([]*Matcher, 0, 1),
		notFound: notFound,
	}
}

// NewSite builds a new HTTPS Site, using the given domain name,
// port number, and certificate files. The provided handler is
// called when a request path does not match any handlers. If
// nil, http.NotFoundHandler is used instead.
//
//		// https://example.com
//		site := NewSecureSite("example.com", 443, "cert.pem", "key.pem", nil)
//
//		// https://example.com:8080
//		site := NewSecureSite("example.com", 8080, "cert.pem", "key.pem", nil)
func NewSecureSite(name string, port int, certFile, keyFile string, notFound http.HandlerFunc) *Site {
	return &Site{
		Name:     name,
		Port:     port,
		auth:     []string{certFile, keyFile},
		handlers: make([]*Matcher, 0, 1),
		notFound: notFound,
	}
}

// Always uses the given handler for any request.
func (s *Site) Always(handler http.Handler) {
	matchFunc := func(_ string) bool { return true }
	s.handlers = append(s.handlers, &Matcher{matchFunc, handler})
}

// Contains uses the given handler when the request path contains
// any of the given pattern strings.
func (s *Site) Contains(handler http.Handler, patterns ...string) {
	for _, pattern := range patterns {
		matchFunc := makeMatchFunc(pattern, strings.Contains)
		s.handlers = append(s.handlers, &Matcher{matchFunc, handler})
	}
}

// Equals uses the given handler when the request path is the same
// as any of the given pattern strings.
func (s *Site) Equals(handler http.Handler, patterns ...string) {
	for _, pattern := range patterns {
		matchFunc := makeMatchFunc(pattern, stringEquals)
		s.handlers = append(s.handlers, &Matcher{matchFunc, handler})
	}
}

// EqualFold uses the given handler when the request path is the same
// as any of the given pattern strings (case insensitive).
func (s *Site) EqualFold(handler http.Handler, patterns ...string) {
	for _, pattern := range patterns {
		matchFunc := makeMatchFunc(pattern, strings.EqualFold)
		s.handlers = append(s.handlers, &Matcher{matchFunc, handler})
	}
}

// HasPrefix uses the given handler when the request path starts with
// any of the given pattern strings.
func (s *Site) HasPrefix(handler http.Handler, patterns ...string) {
	for _, pattern := range patterns {
		matchFunc := makeMatchFunc(pattern, strings.HasPrefix)
		s.handlers = append(s.handlers, &Matcher{matchFunc, handler})
	}
}

// HasSuffix uses the given handler when the request path ends with
// any of the given pattern strings.
func (s *Site) HasSuffix(handler http.Handler, patterns ...string) {
	for _, pattern := range patterns {
		matchFunc := makeMatchFunc(pattern, strings.HasSuffix)
		s.handlers = append(s.handlers, &Matcher{matchFunc, handler})
	}
}

// Match uses the given handler when the given pattern returns true
// when called with the request path.
func (s *Site) Match(handler http.Handler, matchFunc MatchFunc) {
	s.handlers = append(s.handlers, &Matcher{matchFunc, handler})
}

// ServeHTTP allows Site to fulfil the http.Handler interface.
func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	for _, handler := range s.handlers {
		if handler.Match(path) {
			handler.Handler.ServeHTTP(w, r)
			return
		}
	}
	s.notFound(w, r)
}

// Matcher is used to detect and supply an http.Handler.
type Matcher struct {
	Match   MatchFunc
	Handler http.Handler
}

// MatchFunc is used to identify desired request paths.
type MatchFunc func(string) bool

func makeMatchFunc(s1 string, m func(string, string) bool) MatchFunc {
	return func(s2 string) bool {
		return m(s1, s2)
	}
}

func stringEquals(s1, s2 string) bool {
	return s1 == s2
}
