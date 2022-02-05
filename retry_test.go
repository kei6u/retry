package retry

import (
	"context"
	"testing"
	"time"
)

func TestConstant(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		algorithm            algorithm
		exactAttempts        int
		leastAttempts        int
		durationForOverwrite time.Duration
	}{
		{
			name: "timeout",
			algorithm: Constant{
				Context:  timeoutCtx(5 * time.Millisecond),
				Interval: time.Millisecond,
			},
			leastAttempts: 4,
		},
		{
			name: "max attempts",
			algorithm: Constant{
				Interval:    time.Millisecond,
				MaxAttempts: 5,
			},
			exactAttempts: 5,
		},
		{
			name:                 "default",
			algorithm:            Constant{},
			leastAttempts:        3,
			durationForOverwrite: 3 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.durationForOverwrite != 0 {
				overwrite_defaltTimeoutDuration(t, tt.durationForOverwrite)
			}
			r := New(tt.algorithm)
			t.Logf("algorithm: %#v", r)
			attempts := 0
			start := time.Now()
			for r.Next() {
				d := time.Since(start)
				t.Logf("attempt %d, %s elapsed", attempts, d)
				start = time.Now()
				attempts++
			}
			if tt.leastAttempts == 0 && attempts != tt.exactAttempts {
				t.Fatalf("expected to reach %d attempts, actual: %d", tt.exactAttempts, attempts)
			}
			if tt.exactAttempts == 0 && attempts < tt.leastAttempts {
				t.Fatalf("expected to reach %d attempts at least, but actual: %d", tt.leastAttempts, attempts)
			}
		})
	}
}

func TestJitter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		algorithm            algorithm
		exactAttempts        int
		leastAttempts        int
		durationLimit        time.Duration
		durationForOverwrite time.Duration
	}{
		{
			name: "timeout",
			algorithm: Jitter{
				Context: timeoutCtx(10 * time.Millisecond),
				Base:    time.Millisecond,
			},
			leastAttempts: 3,
		},
		{
			name: "max attempts",
			algorithm: Jitter{
				Base:        time.Millisecond,
				MaxAttempts: 5,
			},
			exactAttempts: 5,
		},
		{
			name: "max duration",
			algorithm: Jitter{
				Base:        time.Millisecond,
				Max:         time.Millisecond,
				MaxAttempts: 10,
			},
			durationLimit: 2 * time.Millisecond,
			exactAttempts: 10,
		},
		{
			name:                 "default",
			algorithm:            Jitter{},
			leastAttempts:        2,
			durationForOverwrite: 3 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.durationForOverwrite != 0 {
				overwrite_defaltTimeoutDuration(t, tt.durationForOverwrite)
			}
			r := New(tt.algorithm)
			t.Logf("algorithm: %#v", r)
			attempts := 0
			start := time.Now()
			for r.Next() {
				d := time.Since(start)
				t.Logf("attempt %d, %s elapsed", attempts, d)
				if tt.durationLimit != 0 && d > tt.durationLimit {
					t.Fatalf("expected to limit duration to %s, actual %d", tt.durationLimit, d)
				}
				start = time.Now()
				attempts++
			}
			if tt.leastAttempts == 0 && attempts != tt.exactAttempts {
				t.Fatalf("expected to reach %d attempts, actual: %d", tt.exactAttempts, attempts)
			}
			if tt.exactAttempts == 0 && attempts < tt.leastAttempts {
				t.Fatalf("expected to reach %d attempts at least, but actual: %d", tt.leastAttempts, attempts)
			}
		})
	}
}

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		algorithm            algorithm
		exactAttempts        int
		leastAttempts        int
		durationLimit        time.Duration
		durationForOverwrite time.Duration
	}{
		{
			name: "timeout",
			algorithm: ExponentialBackoff{
				Context: timeoutCtx(10 * time.Millisecond),
				Base:    time.Millisecond,
			},
			leastAttempts: 3,
		},
		{
			name: "max attempts",
			algorithm: ExponentialBackoff{
				Base:        time.Millisecond,
				MaxAttempts: 5,
			},
			exactAttempts: 5,
		},
		{
			name: "max duration",
			algorithm: ExponentialBackoff{
				Base:        time.Millisecond,
				Max:         time.Millisecond,
				MaxAttempts: 10,
			},
			durationLimit: 2 * time.Millisecond,
			exactAttempts: 10,
		},
		{
			name:                 "default",
			algorithm:            ExponentialBackoff{},
			leastAttempts:        2,
			durationForOverwrite: 3 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.durationForOverwrite != 0 {
				overwrite_defaltTimeoutDuration(t, tt.durationForOverwrite)
			}
			r := New(tt.algorithm)
			t.Logf("algorithm: %#v", r)
			attempts := 0
			start := time.Now()
			for r.Next() {
				d := time.Since(start)
				t.Logf("attempt %d, %s elapsed", attempts, d)
				if tt.durationLimit != 0 && d > tt.durationLimit {
					t.Fatalf("expected to limit duration to %s, actual %d", tt.durationLimit, d)
				}
				start = time.Now()
				attempts++
			}
			if tt.leastAttempts == 0 && attempts != tt.exactAttempts {
				t.Fatalf("expected to reach %d attempts, actual: %d", tt.exactAttempts, attempts)
			}
			if tt.exactAttempts == 0 && attempts < tt.leastAttempts {
				t.Fatalf("expected to reach %d attempts at least, but actual: %d", tt.leastAttempts, attempts)
			}
		})
	}
}

func overwrite_defaltTimeoutDuration(t *testing.T, d time.Duration) {
	defaultTimeoutDuration = d
	t.Cleanup(func() {
		// reset to default not to effect other tests.
		defaultTimeoutDuration = time.Minute
	})
}

func timeoutCtx(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}
