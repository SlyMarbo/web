package web

import (
	"net"
	"net/http"
	"strings"
)

// ReverseProxy is used to serve multiple domains simultaneously
// on the same port and interface.
type ReverseProxy struct {
	proxy    map[string]http.Handler // Domain to handler.
	NotFound http.Handler
}

func NewProxy() *ReverseProxy {
	out := new(ReverseProxy)
	out.proxy = make(map[string]http.Handler)
	return out
}

// Register sets the given handler to the domain. If the domain has
// already been assigned, Register will panic.
func (p *ReverseProxy) Register(domain string, h http.Handler) {
	if _, ok := p.proxy[domain]; ok {
		panic("Domain already registered.")
	}
	p.proxy[domain] = h
}

// Register sets the given site to its domain. If the domain has
// already been assigned, Register will panic.
func (p *ReverseProxy) RegisterSite(site *Site) {
	if _, ok := p.proxy[site.Name]; ok {
		panic("Domain already registered.")
	}
	p.proxy[site.Name] = site
}

// ServeHTTP satisfies the http.Handler interface.
func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	var err error
	if strings.Index(host, ":") >= 0 {
		host, _, err = net.SplitHostPort(r.Host)
		if err != nil {
			if p.NotFound != nil {
				p.NotFound.ServeHTTP(w, r)
			}
			return
		}
	}

	for domain, handler := range p.proxy {
		if strings.HasSuffix(host, domain) {
			handler.ServeHTTP(w, r)
			return
		}
	}

	if p.NotFound != nil {
		p.NotFound.ServeHTTP(w, r)
	}
}
