package benchmark

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"sync/atomic"
	"time"
)

// RequestPool manages multiple HTTP requests with different load strategies
type RequestPool struct {
	requests []*http.Request
	strategy string
	counter  uint64 // For round-robin
	rng      *rand.Rand
	sizes    []int // Request sizes in bytes
}

// NewRequestPool creates a new request pool with the specified strategy
func NewRequestPool(requests []*http.Request, strategy string) *RequestPool {
	r := &RequestPool{
		requests: requests,
		strategy: strategy,
		counter:  0,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		sizes:    make([]int, len(requests)),
	}

	for i, s := range r.requests {
		if s.Body != nil {
			body, err := httputil.DumpRequest(r.requests[i], true)
			if err != nil {
				continue
			}
			r.sizes[i] = len(body)
		}
	}

	return r
}

// GetRequest returns the next request based on the load strategy
func (rp *RequestPool) GetRequest() (*http.Request, int) {
	if len(rp.requests) == 0 {
		return nil, 0
	}

	if len(rp.requests) == 1 {
		return rp.requests[0], rp.sizes[0]
	}

	switch rp.strategy {
	case "round-robin":
		idx := atomic.AddUint64(&rp.counter, 1) - 1
		return rp.requests[idx%uint64(len(rp.requests))], rp.sizes[idx%uint64(len(rp.requests))]
	case "random":
		fallthrough
	default:
		idx := rp.rng.Intn(len(rp.requests))
		return rp.requests[idx], rp.sizes[idx]
	}
}

// Size returns the number of requests in the pool
func (rp *RequestPool) Size() int {
	return len(rp.requests)
}
