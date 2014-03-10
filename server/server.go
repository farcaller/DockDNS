package server

import (
	"github.com/farcaller/dockdns/resolver"
	"github.com/miekg/dns"
	"time"
)

type Server struct {
	zone            string
	listenAddress   string
	dockerResolver  resolver.Resolverer
	forwardResolver resolver.Resolverer
}

func New(zone, listenAddress string, dockerResolver, forwardResolver resolver.Resolverer) *Server {
	return &Server{zone, listenAddress, dockerResolver, forwardResolver}
}

func (s *Server) Run() {
	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(s.zone, func(w dns.ResponseWriter, req *dns.Msg) { s.HandleDocker("tcp", w, req) })
	tcpHandler.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) { s.HandleForward("tcp", w, req) })

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(s.zone, func(w dns.ResponseWriter, req *dns.Msg) { s.HandleDocker("udp", w, req) })
	udpHandler.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) { s.HandleForward("udp", w, req) })

	tcpServer := &dns.Server{Addr: s.listenAddress,
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second}

	udpServer := &dns.Server{Addr: s.listenAddress,
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second}

	go func() {
		if err := tcpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	go func() {
		if err := udpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}

func (s *Server) HandleDocker(proto string, w dns.ResponseWriter, req *dns.Msg) {
	s.dockerResolver.Lookup(proto, w, req)
}

func (s *Server) HandleForward(proto string, w dns.ResponseWriter, req *dns.Msg) {
	s.forwardResolver.Lookup(proto, w, req)
}
