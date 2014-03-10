package main

import (
	"flag"
	"github.com/farcaller/dockdns/resolver"
	"github.com/farcaller/dockdns/server"
	"github.com/fsouza/go-dockerclient"
	"log"
	"os"
	"os/signal"
)

var forwardDNS = flag.String("forward", "8.8.8.8:53", "IP address of forwarder DNS")
var dockerZone = flag.String("zone", "docker.", "Docker zone name")
var dockerEndpoint = flag.String("docker", "http://10.0.7.12:5422", "Docker API endpoint")
var listenAddress = flag.String("listen", "127.0.0.1:53", "DNS listen address")

func main() {
	flag.Parse()

	dockerClient, _ := docker.NewClient(*dockerEndpoint)
	dockerResolver := resolver.NewDocker(dockerClient, *dockerZone)
	resolver := resolver.NewForward(*forwardDNS)
	server := server.New(*dockerZone, *listenAddress, dockerResolver, resolver)

	server.Run()
	log.Printf("Server listening on TCP/UDP %s\n", *listenAddress)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-sig:
			return
		}
	}
}
