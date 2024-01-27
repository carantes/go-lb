package lb

import (
	"log"
	"strings"
)

// Algorithm enum.
type Algorithm int

const (
	AlgorithmRoundRobin Algorithm = iota
	AlgorithmLeastConnections
)

var (
	algorithmsMap = map[string]Algorithm{
		"roundrobin":       AlgorithmRoundRobin,
		"leastconnections": AlgorithmLeastConnections,
	}
)

func ParseAlgorithm(str string) Algorithm {
	c, ok := algorithmsMap[strings.ToLower(str)]

	if !ok {
		log.Fatal("Invalid load balancing algorithm")
	}

	return c
}
