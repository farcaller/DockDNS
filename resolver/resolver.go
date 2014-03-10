package resolver

import (
	"github.com/miekg/dns"
)

type Resolverer interface {
	Lookup(proto string, w dns.ResponseWriter, req *dns.Msg)
}
