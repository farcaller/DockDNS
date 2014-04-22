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
	"github.com/fsouza/go-dockerclient"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

type DockerResolver struct {
	client     *docker.Client
	zone       string
	containers map[string]string
}

func NewDocker(client *docker.Client, zone string) Resolverer {
	return &DockerResolver{client, zone, map[string]string{}}
}

func (r *DockerResolver) Lookup(proto string, w dns.ResponseWriter, req *dns.Msg) {
	r.updateContainers()
	question := req.Question[0]
	if question.Qtype != dns.TypeA || question.Qclass != dns.ClassINET {
		dns.HandleFailed(w, req)
		return
	}

	name := strings.TrimSuffix(question.Name, "."+r.zone)
	ipaddr, exists := r.containers[name]
	if !exists {
		dns.HandleFailed(w, req)
		return
	}

	answer := dns.A{
		Hdr: dns.RR_Header{
			Name:     question.Name,
			Rrtype:   dns.TypeA,
			Class:    dns.ClassINET,
			Ttl:      60,
			Rdlength: 4,
		},
		A: net.ParseIP(ipaddr),
	}
	req.Answer = []dns.RR{&answer}
	req.Response = true
	w.WriteMsg(req)
}

func (r *DockerResolver) updateContainers() {
	cont, err := r.client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Printf("Failed to update containers: %v\n", err)
		return
	}

	containers := map[string]string{}

	for _, c := range cont {
		containerJS, err := r.client.InspectContainer(c.ID)
		if err != nil {
			log.Printf("Failed to inspect container %s: %v\n", c.ID, err)
			continue
		}
		ipaddr := containerJS.NetworkSettings.IPAddress
		containers[strings.TrimPrefix(containerJS.Name, "/")] = ipaddr
	}
	r.containers = containers
}
