// DockDNS, the simple docker-aware DNS forwarder.
// Copyright 2014 Vladimir "farcaller" Pouzanov <farcaller@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
