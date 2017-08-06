package xray

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Most code of sampling are copied from GCP's stackdirver trace implements.
// https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/trace/sampling.go

// SamplingPolicy is policy to determine which data can send to AWS API or not, with rate limit function.
type SamplingPolicy interface {
	// Sample returns a Decision.
	// If Trace is false in the returned Decision, then the Decision should be
	// the zero value.
	Sample() Decision
	CanSample() bool
}

// Decision is the value returned by a call to a SamplingPolicy's Sample method.
type Decision struct {
	Trace  bool    // Whether to trace the request.
	Policy string  // Name of the sampling policy.
	Weight float64 // Sample weight to be used in statistical calculations.
}

type sampler struct {
	fraction float64
	skipped  float64
	*rate.Limiter
	*rand.Rand
	sync.Mutex
}

func (s *sampler) Sample() Decision {
	s.Lock()
	x := s.Float64()
	d := s.sample(time.Now(), x)
	s.Unlock()
	return d
}

func (s *sampler) CanSample() bool {
	return s.Sample().Trace
}

// sample contains the a deterministic, time-independent logic of Sample.
func (s *sampler) sample(now time.Time, x float64) (d Decision) {
	d.Trace = x < s.fraction
	if !d.Trace {
		// We have no reason to trace this request.
		return Decision{}
	}
	// We test separately that the rate limit is not tiny before calling AllowN,
	// because of overflow problems in x/time/rate.
	if s.Limit() < 1e-9 || !s.AllowN(now, 1) {
		// Rejected by the rate limit.
		if d.Trace {
			s.skipped++
		}
		return Decision{}
	}
	if d.Trace {
		d.Policy, d.Weight = "default", (1.0+s.skipped)/s.fraction
		s.skipped = 0.0
	}
	return
}

// NewLimitedSampler returns a sampling policy that randomly samples a given
// fraction of requests.  It also enforces a limit on the number of traces per
// second.  It tries to trace every request with a trace header, but will not
// exceed the qps limit to do it.
func NewLimitedSampler(fraction, maxqps float64) (SamplingPolicy, error) {
	if !(fraction >= 0) {
		return nil, fmt.Errorf("invalid fraction %f", fraction)
	}
	if !(maxqps >= 0) {
		return nil, fmt.Errorf("invalid maxqps %f", maxqps)
	}
	// Set a limit on the number of accumulated "tokens", to limit bursts of
	// traced requests.  Use one more than a second's worth of tokens, or 100,
	// whichever is smaller.
	// See https://godoc.org/golang.org/x/time/rate#NewLimiter.
	maxTokens := 100
	if maxqps < 99.0 {
		maxTokens = 1 + int(maxqps)
	}
	var seed int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &seed); err != nil {
		seed = time.Now().UnixNano()
	}
	s := sampler{
		fraction: fraction,
		Limiter:  rate.NewLimiter(rate.Limit(maxqps), maxTokens),
		Rand:     rand.New(rand.NewSource(seed)),
	}
	return &s, nil
}
