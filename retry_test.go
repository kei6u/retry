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
			name: "prefer max attempts to default timeout",
			algorithm: Constant{
				Interval:    time.Millisecond,
				MaxAttempts: 5,
			},
			exactAttempts:        5,
			durationForOverwrite: time.Millisecond,
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
		mostDuration         time.Duration
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
			name: "prefer max attempts to default timeout",
			algorithm: Jitter{
				Base:        time.Millisecond,
				MaxAttempts: 5,
			},
			exactAttempts:        5,
			durationForOverwrite: time.Millisecond,
		},
		{
			name: "max duration",
			algorithm: Jitter{
				Base:        time.Millisecond,
				Max:         time.Millisecond,
				MaxAttempts: 10,
			},
			mostDuration:  2 * time.Millisecond,
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
				if tt.mostDuration != 0 && d > tt.mostDuration {
					t.Fatalf("expected not to exceed %s at most, actual %d", tt.mostDuration, d)
				}
				t.Logf("attempt %d, %s elapsed", attempts, d)
				attempts++
				start = time.Now()
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
		mostDuration         time.Duration
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
			name: "prefer max attempts to default timeout",
			algorithm: ExponentialBackoff{
				Base:        time.Millisecond,
				MaxAttempts: 5,
			},
			exactAttempts:        5,
			durationForOverwrite: time.Millisecond,
		},
		{
			name: "max duration",
			algorithm: ExponentialBackoff{
				Base:        time.Millisecond,
				Max:         time.Millisecond,
				MaxAttempts: 10,
			},
			mostDuration:  2 * time.Millisecond,
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
				if tt.mostDuration != 0 && d > tt.mostDuration {
					t.Fatalf("expected not to exceed %s at most, actual %d", tt.mostDuration, d)
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
	ctx, cancel := context.WithTimeout(context.Background(), d)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return ctx
}

func TestJitter_calc(t *testing.T) {
	t.Parallel()
	j := Jitter{
		Base: time.Millisecond,
		Max:  time.Hour,
	}
	prev := time.Millisecond
	isJitter := false
	for i := 0; i < 10; i++ {
		d := j.calc()
		t.Logf("calc %d, %s", i, d)
		if d < prev {
			isJitter = true
		}
		prev = d
	}
	if !isJitter {
		t.FailNow()
	}
}

func TestExponentialBackoff_calc(t *testing.T) {
	t.Parallel()
	b := ExponentialBackoff{
		Base: time.Millisecond,
		Max:  time.Hour,
	}
	prev := time.Millisecond
	for i := 0; i < 10; i++ {
		d := b.calc()
		t.Logf("calc %d, %s", i, d)
		if d < prev {
			t.Fatalf("calculated duration must be greater than previous one")
		}
		prev = d
	}
}
