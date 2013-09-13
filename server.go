package web

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/SlyMarbo/spdy"
	"net"
	"net/http"
)

// Server serves requests for all attached sites.
// If multiple sites use the same port, it will provide
// a reverse proxy service.
type Server struct {
	sites []*Site
}

// NewServer creates and initialises a Server.
func NewServer() *Server {
	s := new(Server)
	s.sites = make([]*Site, 0, 1)
	return s
}

// NewServerFromSites creates a Server with the
// provided sites.
func NewServerFromSites(sites ...*Site) *Server {
	s := new(Server)
	s.sites = sites
	return s
}

// Add takes a Site and adds it to the Server's collection.
// Add returns the server, so multiple Add calls can be
// chained together.
func (s *Server) Add(site *Site) *Server {
	s.sites = append(s.sites, site)
	return s
}

// Serve starts listening and serving requests to the
// provided sites. Serve uses an extra goroutine for
// each port used by its sites.
func (s *Server) Serve() error {
	errChan := make(chan error)
	portMap := make(map[int][]*Site)

	// Collect sites by port.
	for _, site := range s.sites {
		if sites, ok := portMap[site.Port]; ok {
			portMap[site.Port] = append(sites, site)
		} else {
			portMap[site.Port] = []*Site{site}
		}
	}

	// Iterate through site groups by port.
	for port, sites := range portMap {

		// Single sites on a port are simple.
		if len(sites) == 1 {
			site := sites[0]
			if site.auth != nil {
				if site.SPDY {
					go serveSPDY(site.Port, site, site.auth[0], site.auth[1], errChan)
				} else {
					go serveHTTPS(site.Port, site, site.auth[0], site.auth[1], errChan)
				}
			} else {
				go serveHTTP(site.Port, site, errChan)
			}
		} else {

			// Make sure all sites on one port use TLS or none do.
			auth := sites[0].auth != nil
			for i := 1; i < len(sites); i++ {
				if auth != (sites[i].auth != nil) {
					return errors.New("Multiple sites on the same port with mixed HTTPS usage.")
				}
			}

			// Build reverse proxy and Server Name Identification (SNI) configuration.
			addr := fmt.Sprintf(":%d", port)
			proxy := NewProxy()
			server := &http.Server{Addr: addr, Handler: proxy}
			tlsConf := &tls.Config{NextProtos: []string{"http/1.1"}}
			tlsConf.Certificates = make([]tls.Certificate, len(sites))
			var err error
			for i, site := range sites {
				proxy.RegisterSite(site)
				if auth {
					// Add certificate pair if using TLS.
					tlsConf.Certificates[i], err = tls.LoadX509KeyPair(site.auth[0], site.auth[1])
					if err != nil {
						return err
					}
				}
			}

			// Build SNI.
			tlsConf.BuildNameToCertificate()

			// Create the TCP listener.
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			// Add TLS if necessary.
			if auth {
				tlsListener := tls.NewListener(listener, tlsConf)
				server.TLSConfig = tlsConf
				listener = tlsListener
			}

			// Start serving.
			go serveMany(server, listener, errChan)
		}
	}

	// Keep running until an error occurs.
	return <-errChan
}

func serveSPDY(port int, handler http.Handler, certFile, keyFile string, errChan chan<- error) {
	addr := fmt.Sprintf(":%d", port)
	err := spdy.ListenAndServeTLS(addr, certFile, keyFile, handler)
	if err != nil {
		errChan <- err
	}
}

func serveHTTPS(port int, handler http.Handler, certFile, keyFile string, errChan chan<- error) {
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServeTLS(addr, certFile, keyFile, handler)
	if err != nil {
		errChan <- err
	}
}

func serveHTTP(port int, handler http.Handler, errChan chan<- error) {
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		errChan <- err
	}
}

func serveMany(server *http.Server, listener net.Listener, errChan chan<- error) {
	err := server.Serve(listener)
	if err != nil {
		errChan <- err
	}
}
