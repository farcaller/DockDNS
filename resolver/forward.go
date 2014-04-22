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
