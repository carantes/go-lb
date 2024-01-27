package main

import (
	"flag"
	"log"
	"net/url"
	"strings"

	"github.com/carantes/golb/lb"
	"github.com/carantes/golb/server"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
	st   = flag.String("type", "node", "server type: node or lb")

	// Server specific flags
	// lbAddr = flag.String("lbAddr", "http://localhost:8080", "load balancer address")

	// Load balancer specific flags
	alg      = flag.String("alg", "roundrobin", "load balancing algorithm: roundrobin or random")
	urls     = flag.String("urls", "http://localhost:8081,http://localhost:8082", "Initial urls to distribute the load")
	interval = flag.Int("interval", 5, "Heartbeat interval in seconds")
)

// ParseNodeFromString parses a string of nodes into a Node list.
func ParseNodeFromString(nodes string) []*lb.Node {
	var parsedNodes []*lb.Node

	for _, node := range strings.Split(nodes, ",") {
		parsedNode, err := url.Parse(node)

		if err != nil {
			log.Fatal("Invalid node URL")
		}

		parsedNodes = append(parsedNodes, lb.NewNode(*parsedNode))
	}

	return parsedNodes
}

func main() {

	flag.Parse()

	// Start a server node
	if *st == "node" {
		server.NewServer().Run(*addr)
	}

	// Start the load balancer
	if *st == "lb" {
		algorithm := lb.ParseAlgorithm(*alg)
		nodes := ParseNodeFromString(*urls)
		lb.NewLoadBalancer(nodes, algorithm, *interval).Run(*addr)
	}

	panic("Invalid server type")
}
