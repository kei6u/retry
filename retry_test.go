package retry

import (
	"context"
	"testing"
	"time"
)

func TestConstant(t *testing.T) {
	t.Run("set empty Constant struct", func(t *testing.T) {
		overwrite_defaltTimeoutDuration(t, 10*time.Millisecond)
		start := time.Now()
		retrier := New(Constant{})
		for retrier.Next() {
		}
		if time.Since(start) < 10*time.Millisecond {
			t.Fatalf("expected to timeout after 100ms")
		}
	})
	t.Run("set max attempts only", func(t *testing.T) {
		retryCount := 0
		retrier := New(Constant{
			Interval:    time.Millisecond,
			MaxAttempts: 10,
		})
		for retrier.Next() {
			retryCount++
		}
		if retryCount != 10 {
			t.Fatalf("expected 10 retries, but %d", retryCount)
		}
	})
	t.Run("set timeout only", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		retrier := New(Constant{
			Context:  ctx,
			Interval: time.Second,
		})
		start := time.Now()
		for retrier.Next() {
		}
		if time.Since(start) < 10*time.Millisecond {
			t.Fatal("expected timeout after 10ms")
		}
	})
	t.Run("prioritize max attempts", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		retrier := New(Constant{
			Context:     ctx,
			Interval:    time.Millisecond,
			MaxAttempts: 10,
		})
		start := time.Now()
		retryCount := 0
		for retrier.Next() {
			retryCount++
		}
		if time.Second < time.Since(start) && retryCount < 10 {
			t.Fatal("expected to reach 10 retries before 1 second elapsed")
		}
	})
	t.Run("prioritize timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		retrier := New(Constant{
			Context:     ctx,
			Interval:    2 * time.Millisecond,
			MaxAttempts: 10,
		})
		start := time.Now()
		retryCount := 0
		for retrier.Next() {
			t.Logf("retry %d, %s elapsed", retryCount, time.Since(start))
			retryCount++
		}
		if 10 <= retryCount && 10*time.Millisecond < time.Since(start) {
			t.Logf("%d retries took %s", retryCount, time.Since(start))
			t.Fatal("expected to timeout before 10 retries")
		}
	})
}

func TestExponentialBackoff(t *testing.T) {
	t.Run("set empty ExponentialBackoff struct", func(t *testing.T) {
		overwrite_defaltTimeoutDuration(t, 10*time.Millisecond)
		start := time.Now()
		retrier := New(ExponentialBackoff{})
		for retrier.Next() {
		}
		if time.Since(start) < 10*time.Millisecond {
			t.Fatalf("expected to timeout after 10ms")
		}
	})
	t.Run("set max attempts only", func(t *testing.T) {
		retrier := New(ExponentialBackoff{
			BaseInterval: time.Millisecond,
			MaxAttempts:  10,
		})
		start := time.Now()
		retryCount := 0
		for retrier.Next() {
			t.Logf("retry %d, %s elapsed", retryCount, time.Since(start))
			retryCount++
		}
		if retryCount != 10 && time.Since(start) < 600*time.Millisecond {
			t.Fatalf(
				"expected 10 retries, elapse at least 600ms, actual: %d retries, %s elapsed",
				retryCount,
				time.Since(start),
			)
		}
	})
	t.Run("set timeout only", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		retrier := New(ExponentialBackoff{
			Context:      ctx,
			BaseInterval: time.Second,
		})
		start := time.Now()
		for retrier.Next() {
		}
		if time.Since(start) < 10*time.Millisecond {
			t.Fatalf("expected to timeout 10ms, but %s elapsed", time.Since(start))
		}
	})
	t.Run("max attempts is prioritized", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		retrier := New(ExponentialBackoff{
			Context:      ctx,
			BaseInterval: time.Millisecond,
			MaxAttempts:  10,
		})
		start := time.Now()
		retryCount := 0
		for retrier.Next() {
			retryCount++
		}
		if time.Second < time.Since(start) && retryCount < 10 {
			t.Fatal("expected to reach 10 retries before 1 second elapsed")
		}
	})
	t.Run("timeout is prioritized", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		retrier := New(ExponentialBackoff{
			Context:      ctx,
			BaseInterval: time.Millisecond,
			MaxAttempts:  10,
		})
		start := time.Now()
		retryCount := 0
		for retrier.Next() {
			retryCount++
		}
		if 10 <= retryCount && 10*time.Millisecond < time.Since(start) {
			t.Fatal("expected to timeout before 10 retries")
		}
	})
}

func overwrite_defaltTimeoutDuration(t *testing.T, d time.Duration) {
	defaultTimeoutDuration = d
	t.Cleanup(func() {
		// reset to default not to effect other tests.
		defaultTimeoutDuration = time.Minute
	})
}
