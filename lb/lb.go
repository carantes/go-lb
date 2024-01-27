package lb

import (
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

type LoadBalancer struct {
	mutex            *sync.Mutex
	wg               *sync.WaitGroup
	nodes            []*Node
	currentNode      int
	activeNodesCount int
	checkInterval    int
	ready            bool
}

// NewLoadBalancer creates a new load balancer.
func NewLoadBalancer(nodes []*Node, alg Algorithm, interval int) *LoadBalancer {
	return &LoadBalancer{
		mutex:            &sync.Mutex{},
		wg:               &sync.WaitGroup{},
		nodes:            nodes,
		currentNode:      0,
		activeNodesCount: 0,
		checkInterval:    interval,
		ready:            false,
	}
}

// Run starts the load balancer.
func (lb *LoadBalancer) Run(addr string) {
	// Run the health check in a goroutine
	go lb.healthcheck()

	// Proxy request to the next server
	err := http.ListenAndServe(addr, http.HandlerFunc(lb.proxyHandler))

	if err != nil {
		log.Fatal(err)
	}
}

// startHealthcheck starts the health check.
func (lb *LoadBalancer) healthcheck() {
	for {
		for _, node := range lb.nodes {
			lb.wg.Add(1)
			go lb.nodeCheck(node)
		}

		// Wait for all goroutines finish to start the next round
		lb.wg.Wait()
		log.Println("Active servers:", lb.activeNodesCount)
		time.Sleep(time.Duration(lb.checkInterval) * time.Second)
	}
}

// nodeCheck checks if a node is healthy.
func (lb *LoadBalancer) nodeCheck(node *Node) {
	defer lb.wg.Done()

	// Check if the node is healthy
	// TODO: Implement a retry mechanism
	log.Println("Checking node", node.URL.String())
	ok := node.Check()

	if node.active != ok {
		lb.mutex.Lock()

		// Update the node status
		node.SetActive(ok)

		// Update the active nodes count
		if ok {
			lb.activeNodesCount++
		} else {
			lb.activeNodesCount--
		}

		// Update the load balancer status
		if lb.activeNodesCount > 0 {
			lb.ready = true
		} else {
			log.Println("Load balancer is not ready")
			lb.ready = false
		}

		lb.mutex.Unlock()
	}

	log.Println(node.Stats())
}

// Proxy request to the next server.
func (lb *LoadBalancer) proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Load balancer is not ready
	if !lb.ready {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	// Get next active server
	node := lb.nextServer()
	target := &node.URL

	// No available servers
	if target == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("No available servers"))
		return
	}

	// Proxy request to the next server
	log.Printf("Round Robin to [%d] %s", lb.currentNode, target.String())
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		// Increment the request count
		node.stats.IncRequestCount()

		// Rewrite the request to point to the next server
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.Host = target.Host
		req.Header = r.Header
	}
	proxy.ModifyResponse = func(res *http.Response) error {
		node.stats.IncSuccessCount()
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		node.stats.IncFailureCount()
		node.SetActive(false)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
	}

	proxy.ServeHTTP(w, r)
}

// Next server returns the next healthy server.
func (lb *LoadBalancer) nextServer() *Node {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i := 0; i < len(lb.nodes); i++ {
		lb.currentNode = (lb.currentNode + 1) % len(lb.nodes)

		if lb.nodes[lb.currentNode].active {
			return lb.nodes[lb.currentNode]
		}
	}

	return nil
}
