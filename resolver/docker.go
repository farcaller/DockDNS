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
		for _, n := range c.Names {
			n = strings.TrimPrefix(n, "/")
			containers[n] = ipaddr
		}
	}
	r.containers = containers
}
