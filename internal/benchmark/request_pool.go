package benchmark

import (
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

// RequestPool manages multiple HTTP requests with different load strategies
type RequestPool struct {
	requests []*http.Request
	strategy string
	counter  uint64 // For round-robin
	rng      *rand.Rand
}

// NewRequestPool creates a new request pool with the specified strategy
func NewRequestPool(requests []*http.Request, strategy string) *RequestPool {
	return &RequestPool{
		requests: requests,
		strategy: strategy,
		counter:  0,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetRequest returns the next request based on the load strategy
func (rp *RequestPool) GetRequest() *http.Request {
	if len(rp.requests) == 0 {
		return nil
	}

	if len(rp.requests) == 1 {
		return rp.requests[0]
	}

	switch rp.strategy {
	case "round-robin":
		idx := atomic.AddUint64(&rp.counter, 1) - 1
		return rp.requests[idx%uint64(len(rp.requests))]
	case "random":
		fallthrough
	default:
		idx := rp.rng.Intn(len(rp.requests))
		return rp.requests[idx]
	}
}

// Size returns the number of requests in the pool
func (rp *RequestPool) Size() int {
	return len(rp.requests)
}
