package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Options provides retry configurations.
type Options struct {
	ctx          context.Context
	factor       float64 // factor controls retry interval ranges.
	baseInterval time.Duration
	maxInterval  time.Duration
	maxAttempts  float64
	attempts     float64
}

// Next returns true if the next retry should be performed
// and waits for the interval before the next retry.
func (o *Options) Next() bool {
	if o.ctx == nil {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			defaultTimeoutDuration,
		)
		o.ctx = ctx
		go func() {
			<-ctx.Done()
			cancel()
		}()
	}
	defer func() {
		o.attempts++
	}()
	if o.attempts == 0 {
		return true
	}
	if o.attempts == o.maxAttempts {
		return false
	}
	interval := float64(o.baseInterval) * math.Pow(2, o.attempts)
	factoredInterval := interval / o.factor
	waitDuration := time.Duration(randomBetween(factoredInterval, interval))
	if o.maxInterval < waitDuration {
		waitDuration = o.maxInterval
	}
	select {
	case <-o.ctx.Done():
		return false
	case <-time.After(waitDuration):
		return true
	}
}

// defining this as a global variable for testing.
var defaultTimeoutDuration = time.Minute

// randomBetween returns a random float64 number between min and max.
func randomBetween(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

// ConstantOptions provides options for constant intervals.
type ConstantOptions struct {
	// Context is for timeout or canceling retry loop.
	Context context.Context
	// Interval is the interval between retries.
	Interval time.Duration
	// MaxAttempts is the maximum number of attempts to retry.
	MaxAttempts uint
}

// DefaultConstant is a default configuration for constant interval retry.
// An interval duration is a second and timeout is a minute.
var DefaultConstant = Constant(ConstantOptions{})

// NewConstant returns a constant interval retry configuration.
func Constant(opts ConstantOptions) Options {
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

	return Options{
		ctx:          opts.Context,
		factor:       factor,
		baseInterval: baseInterval,
		maxInterval:  maxInterval,
		maxAttempts:  float64(maxAttempts),
		attempts:     attempts,
	}
}

// ExponentialBackoffOptions provides options for the exponential backoff algorithm.
//
// An interval can be computed by this expression.
//
// interval = baseInterval * (2 ^ retryAttempts)
//
// Then, randomly choose a float64 number from `interval / 2` to `interval`.
// If a chosen float64 number is more than maxInterval, use maxInterval instead.
type ExponentialBackoffOptions struct {
	// Context is for timeout or canceling retry loop.
	Context context.Context
	// BaseInterval controls the rate of exponential backoff interval growth.
	BaseInterval time.Duration
	// MaxInterval is the maximum wait duration to retry.
	MaxInterval time.Duration
	// MaxAttempts is the maximum number of attempts to retry.
	MaxAttempts uint
}

// DefaultExponentialBackoff is a default configuration for exponential backoff retry.
// Base interval is a second, max interval is a minute and timeout is a minute.
var DefaultExponentialBackoff = ExponentialBackoff(ExponentialBackoffOptions{})

// ExponentialBackoff creates a new exponential backoff retry configuration.
func ExponentialBackoff(opts ExponentialBackoffOptions) Options {
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

	return Options{
		ctx:          opts.Context,
		factor:       factor,
		baseInterval: baseInterval,
		maxInterval:  maxInterval,
		maxAttempts:  float64(maxAttempts),
		attempts:     attempts,
	}
}
