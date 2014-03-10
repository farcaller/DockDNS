package resolver

import (
	"github.com/miekg/dns"
	"time"
)

type ForwardResolver struct {
	server string
}

func NewForward(forwardDNS string) Resolverer {
	return &ForwardResolver{forwardDNS}
}

func (r *ForwardResolver) Lookup(proto string, w dns.ResponseWriter, req *dns.Msg) {
	c := &dns.Client{
		Net:          proto,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	in, _, _ := c.Exchange(req, r.server)
	w.WriteMsg(in)
}
