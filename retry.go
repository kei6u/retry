package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Options provides retry configuration.
type Options struct {
	ctx          context.Context
	factor       float64 // factor controls retry interval ranges.
	baseInterval time.Duration
	maxInterval  time.Duration
	maxAttempts  float64
	attempts     float64
}

// Next returns true if the next retry should be performed.
// It waits for the configured interval before the next retry.
func (o *Options) Next() bool {
	if o.ctx == nil {
		ctx, cancel := defaultContextWithTimeout()
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

// defaultContextWithTimeout retruns context with default timeout.
// The reason why it is implemented as a variable is for testing.
var defaultContextWithTimeout = func() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), time.Minute)
}

// randomBetween returns a random float64 number between min and max.
func randomBetween(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

// ConstantOptions provides constant interval options.
type ConstantOptions struct {
	Context     context.Context
	Interval    time.Duration
	MaxAttempts uint
}

// Constant is a default configuration for constant interval retry.
// An interval duration is a second and time out is a minute.
var DefaultConstant = Constant(ConstantOptions{})

// NewConstant returns a constant interval retry configuration.
func Constant(opts ConstantOptions) Options {
	var (
		ctx          context.Context
		maxAttempts  uint
		baseInterval               = time.Second
		maxInterval  time.Duration = time.Second
	)
	const (
		factor   float64 = 1 // for constant interval.
		attempts float64 = 0
	)

	if opts.Context != nil {
		ctx = opts.Context
	}
	if opts.Interval != 0 {
		baseInterval = opts.Interval
		maxInterval = opts.Interval
	}
	if opts.MaxAttempts != 0 {
		maxAttempts = opts.MaxAttempts
	}

	return Options{
		ctx:          ctx,
		factor:       factor,
		baseInterval: baseInterval,
		maxInterval:  maxInterval,
		maxAttempts:  float64(maxAttempts),
		attempts:     attempts,
	}
}

// ExponentialBackoffOptions provides exponential backoff options.
type ExponentialBackoffOptions struct {
	// Context is for timeout or canceling retry loop.
	Context context.Context
	// BaseInterval controls the rate of exponential backoff interval growth.
	// An interval can be computed by the expression below.
	// interval = baseInterval * (2 ^ retryAttempts)
	// Randomly choose between interval / 2 and interval if the result is less than maxInterval.
	BaseInterval time.Duration
	// MaxInterval is the maximum duration to wait for a retry.
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
		ctx          context.Context
		maxAttempts  uint
		baseInterval = time.Second
		maxInterval  = 64 * time.Second
	)
	const (
		factor   float64 = 2
		attempts float64 = 0
	)

	if opts.Context != nil {
		ctx = opts.Context
	}
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
		ctx:          ctx,
		factor:       factor,
		baseInterval: baseInterval,
		maxInterval:  maxInterval,
		maxAttempts:  float64(maxAttempts),
		attempts:     attempts,
	}
}
