package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// retrier provides retry functionalities.
type retrier struct {
	calculator
	ctx         context.Context
	maxAttempts float64
	attempts    float64
}

// calculator calculates duration to wait for next retry.
type calculator interface {
	calc() time.Duration
}

// Next returns true if the next retry should be performed
// and waits for the interval before the next retry.
func (r *retrier) Next() bool {
	defer func() {
		r.attempts++
	}()
	if r.ctx == nil {
		if r.maxAttempts == 0 {
			// Set timeout to prevent infinite loop.
			ctx, cancel := context.WithTimeout(
				context.Background(),
				defaultTimeoutDuration,
			)
			r.ctx = ctx
			go func() {
				<-ctx.Done()
				cancel()
			}()
		} else {
			// Prefer max attempts over timeout.
			r.ctx = context.Background()
		}
	}
	if r.attempts == 0 {
		return true
	}
	if r.attempts == r.maxAttempts {
		return false
	}
	select {
	case <-r.ctx.Done():
		return false
	case <-time.After(time.Duration(r.calc())):
		return true
	}
}

type algorithm interface {
	new() retrier
}

// New creates a new Retrier.
func New(a algorithm) retrier {
	return a.new()
}

// defining this as a global variable for testing.
var defaultTimeoutDuration = time.Minute

// randomBetween returns a random float64 number between min and max.
func randomBetween(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

// Jitter provides options for jitter intervals.
// You can set empty for any fields, it will use default values.
//
// An interval can be computed by this expression.
//
// interval = min(max, randomBetween(base, interval * 3))
//
// Example: Given 1 second for Base, the sequence 10 retries will be:
//
// Retry #1:  2.20932s
// Retry #2:  6.293147s
// Retry #3:  12.881962s
// Retry #4:  17.478132s
// Retry #5:  22.84098s
// Retry #6:  47.376313s
// Retry #7:  10.263282s
// Retry #8:  5.662684s
// Retry #9:  2.550353s
// Retry #10: 6.793149s
type Jitter struct {
	// Context is for timeout or canceling retry loop. Default is 1 minute timeout.
	Context context.Context
	// Base is the base wait duration to retry. Default is 1 second.
	Base time.Duration
	// Max is the maximum wait duration to retry. Default is 15 seconds.
	Max time.Duration
	// MaxAttempts is the maximum number of retries. Default is 0.
	// If set 0, it will prioritize timeout.
	MaxAttempts float64

	interval time.Duration
}

func (j *Jitter) calc() time.Duration {
	if j.interval == 0 {
		j.interval = j.Base
	}
	d := time.Duration(math.Min(
		float64(j.Max),
		randomBetween(float64(j.Base), float64(j.interval)*3),
	))
	j.interval = d
	return d
}

func (j Jitter) new() retrier {
	if j.Base == 0 {
		j.Base = time.Second
	}
	if j.Max == 0 {
		j.Max = time.Minute
	}
	return retrier{
		calculator:  &j,
		ctx:         j.Context,
		maxAttempts: j.MaxAttempts,
	}
}

// Constant provides options for constant intervals.
// You can set empty for any fields, it will use default values.
type Constant struct {
	// Context is for timeout or canceling retry loop. Default is 1 minute timeout.
	Context context.Context
	// Interval is the interval between retries. Default is 1 second.
	Interval time.Duration
	// MaxAttempts is the maximum number of retries. Default is 0.
	// If set 0, it will prioritize timeout.
	MaxAttempts float64
}

func (c Constant) calc() time.Duration {
	return c.Interval
}

func (c Constant) new() retrier {
	if c.Interval == 0 {
		c.Interval = time.Second
	}
	return retrier{
		calculator:  c,
		ctx:         c.Context,
		maxAttempts: c.MaxAttempts,
	}
}

// ExponentialBackoff provides options for the exponential backoff algorithm.
// You can set empty for any fields, it will use default values.
//
// An interval can be computed by this expression.
//
// temp = base * (2 ^ attempts)
// interval = min(max, randomBetween(temp / 2, temp))
//
// Example: Given 1 second for Base and 2 minute for Max, the sequence 10 retries will be:
//
// Retry #1:  1.60466s
// Retry #2:  3.881018s
// Retry #3:  6.65824s
// Retry #4:  11.501713s
// Retry #5:  22.84098s
// Retry #6:  53.978338s
// Retry #7:  68.200769s
// Retry #8:  2m
// Retry #9:  2m
// Retry #10: 2m
type ExponentialBackoff struct {
	// Context is for timeout or canceling retry loop. Default is 1 minute timeout.
	Context context.Context
	// Base controls the rate of exponential backoff interval growth.
	// Default is 1 second.
	Base time.Duration
	// Max is the maximum wait duration to retry. Default is 15 seconds.
	Max time.Duration
	// MaxAttempts is the maximum number of retries. Default is 0.
	// If set 0, it will prioritize timeout.
	MaxAttempts float64

	attempt float64
}

func (b *ExponentialBackoff) calc() time.Duration {
	b.attempt++
	temp := float64(b.Base) * math.Pow(2, b.attempt)
	return time.Duration(math.Min(
		float64(b.Max),
		randomBetween(temp/2, temp),
	))
}

func (b ExponentialBackoff) new() retrier {
	if b.Base == 0 {
		b.Base = time.Second
	}
	if b.Max == 0 {
		b.Max = 15 * time.Second
	}
	return retrier{
		calculator:  &b,
		ctx:         b.Context,
		maxAttempts: b.MaxAttempts,
	}
}
