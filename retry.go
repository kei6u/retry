package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// retrier provides retry functionalities.
type retrier struct {
	ctx          context.Context
	factor       float64 // factor controls retry interval ranges.
	baseInterval time.Duration
	maxInterval  time.Duration
	maxAttempts  float64
	attempts     float64
}

// Next returns true if the next retry should be performed
// and waits for the interval before the next retry.
func (r *retrier) Next() bool {
	if r.ctx == nil {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			defaultTimeoutDuration,
		)
		r.ctx = ctx
		go func() {
			<-ctx.Done()
			cancel()
		}()
	}
	defer func() {
		r.attempts++
	}()
	if r.attempts == 0 {
		return true
	}
	if r.attempts == r.maxAttempts {
		return false
	}
	interval := float64(r.baseInterval) * math.Pow(2, r.attempts)
	factoredInterval := interval / r.factor
	waitDuration := time.Duration(randomBetween(factoredInterval, interval))
	if r.maxInterval < waitDuration {
		waitDuration = r.maxInterval
	}
	select {
	case <-r.ctx.Done():
		return false
	case <-time.After(waitDuration):
		return true
	}
}

type strategy interface {
	new() retrier
}

// New creates a new Retrier.
func New(s strategy) retrier {
	return s.new()
}

// defining this as a global variable for testing.
var defaultTimeoutDuration = time.Minute

// randomBetween returns a random float64 number between min and max.
func randomBetween(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

var _ strategy = (*Constant)(nil)

// Constant provides options for constant intervals.
type Constant struct {
	// Context is for timeout or canceling retry loop.
	Context context.Context
	// Interval is the interval between retries.
	Interval time.Duration
	// MaxAttempts is the maximum number of attempts to retry.
	MaxAttempts uint
}

func (opts Constant) new() retrier {
	var (
		maxAttempts  uint
		baseInterval = time.Second
		maxInterval  = time.Second
	)
	const (
		factor   float64 = 1 // for constant interval.
		attempts float64 = 0
	)

	if opts.Interval != 0 {
		baseInterval = opts.Interval
		maxInterval = opts.Interval
	}
	if opts.MaxAttempts != 0 {
		maxAttempts = opts.MaxAttempts
	}

	return retrier{
		ctx:          opts.Context,
		factor:       factor,
		baseInterval: baseInterval,
		maxInterval:  maxInterval,
		maxAttempts:  float64(maxAttempts),
		attempts:     attempts,
	}
}

var _ strategy = (*ExponentialBackoff)(nil)

// ExponentialBackoff provides options for the exponential backoff algorithm.
//
// An interval can be computed by this expression.
//
// interval = baseInterval * (2 ^ retryAttempts)
//
// Then, randomly choose a float64 number from `interval / 2` to `interval`.
// If a chosen float64 number is more than maxInterval, use maxInterval instead.
type ExponentialBackoff struct {
	// Context is for timeout or canceling retry loop.
	Context context.Context
	// BaseInterval controls the rate of exponential backoff interval growth.
	BaseInterval time.Duration
	// MaxInterval is the maximum wait duration to retry.
	MaxInterval time.Duration
	// MaxAttempts is the maximum number of attempts to retry.
	MaxAttempts uint
}

func (opts ExponentialBackoff) new() retrier {
	var (
		maxAttempts  uint
		baseInterval = time.Second
		maxInterval  = 64 * time.Second
	)
	const (
		factor   float64 = 2
		attempts float64 = 0
	)

	if opts.BaseInterval != 0 {
		baseInterval = opts.BaseInterval
	}
	if opts.MaxInterval != 0 {
		maxInterval = opts.MaxInterval
	}
	if opts.MaxAttempts != 0 {
		maxAttempts = opts.MaxAttempts
	}

	return retrier{
		ctx:          opts.Context,
		factor:       factor,
		baseInterval: baseInterval,
		maxInterval:  maxInterval,
		maxAttempts:  float64(maxAttempts),
		attempts:     attempts,
	}
}
