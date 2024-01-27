package lb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// stats represents the load balancer stats.
type stats struct {
	startTime     time.Time
	requestCount  int64
	successCount  int64
	failuresCount int64
}

func (s *stats) IncRequestCount() {
	atomic.AddInt64(&s.requestCount, 1)
}
func (s *stats) IncSuccessCount() {
	atomic.AddInt64(&s.successCount, 1)
}
func (s *stats) IncFailureCount() {
	atomic.AddInt64(&s.failuresCount, 1)
}

// Node represents a server node.
type Node struct {
	mutex  *sync.Mutex
	client *http.Client
	url.URL
	active bool
	stats
	lastCheck time.Time
}

type HealthcheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// NewNode creates a new node.
func NewNode(url url.URL) *Node {
	return &Node{
		mutex: &sync.Mutex{},
		client: &http.Client{
			// TODO: Check if this is the best way to configure the http client
			Transport: &http.Transport{
				DisableKeepAlives:  false,
				DisableCompression: false,
				MaxIdleConns:       1,
				IdleConnTimeout:    10,
			},
			Timeout: 5 * time.Second,
		},
		URL:    url,
		active: false,
		stats: stats{
			startTime: time.Now(),
		},
		lastCheck: time.Now().AddDate(0, 0, -1),
	}
}

// Get health check url.
func (n *Node) getHealthCheckURL() string {
	return n.URL.String() + "/status"
}

// Check makes a health check request to the node and updates its status.
func (n *Node) Check() bool {
	// Start the health check
	n.lastCheck = time.Now()
	n.stats.IncRequestCount()

	if res, err := http.Get(n.getHealthCheckURL()); err == nil {
		defer res.Body.Close()

		// Network request succeeded
		if res.StatusCode == http.StatusOK {
			var data HealthcheckResponse
			err = json.NewDecoder(res.Body).Decode(&data)

			// Response payload is valid
			if err == nil && strings.ToLower(data.Status) == "ok" {
				n.stats.IncSuccessCount()
				return true
			}
		}
	}

	// Network request failed
	n.stats.IncFailureCount()
	return false
}

func (n *Node) SetActive(active bool) {
	if n.active == active {
		return
	}

	// Lock the node to update its status
	n.mutex.Lock()
	n.active = active
	n.mutex.Unlock()
}

func (n *Node) Stats() string {
	return fmt.Sprintf("Node %s req: %d, success: %d , failures: %d", n.URL.String(), n.requestCount, n.successCount, n.failuresCount)
}
