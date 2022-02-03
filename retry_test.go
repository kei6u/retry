package retry

import (
	"context"
	"testing"
	"time"
)

func TestConstant(t *testing.T) {
	t.Run("use DefaultConstant of global variable", func(t *testing.T) {
		overwrite_defaltTimeoutDuration(t, 10*time.Millisecond)
		start := time.Now()
		for DefaultConstant.Next() {
		}
		if time.Since(start) < 10*time.Millisecond {
			t.Fatalf("expected to timeout after 100ms")
		}
	})
	t.Run("set max attempts only", func(t *testing.T) {
		retryCount := 0
		constant := Constant(ConstantOptions{
			Interval:    time.Millisecond,
			MaxAttempts: 10,
		})
		for constant.Next() {
			retryCount++
		}
		if retryCount != 10 {
			t.Fatalf("expected 10 retries, but %d", retryCount)
		}
	})
	t.Run("set timeout only", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		constant := Constant(ConstantOptions{
			Context:  ctx,
			Interval: time.Second,
		})
		start := time.Now()
		for constant.Next() {
		}
		if time.Since(start) < 10*time.Millisecond {
			t.Fatal("expected timeout after 10ms")
		}
	})
	t.Run("prioritize max attempts", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		constant := Constant(ConstantOptions{
			Context:     ctx,
			Interval:    time.Millisecond,
			MaxAttempts: 10,
		})
		start := time.Now()
		retryCount := 0
		for constant.Next() {
			retryCount++
		}
		if time.Second < time.Since(start) && retryCount < 10 {
			t.Fatal("expected to reach 10 retries before 1 second elapsed")
		}
	})
	t.Run("prioritize timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		constant := Constant(ConstantOptions{
			Context:     ctx,
			Interval:    2 * time.Millisecond,
			MaxAttempts: 10,
		})
		start := time.Now()
		retryCount := 0
		for constant.Next() {
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
	t.Run("use DefaultExponentialBackoff of global variable", func(t *testing.T) {
		overwrite_defaltTimeoutDuration(t, 10*time.Millisecond)
		start := time.Now()
		for DefaultExponentialBackoff.Next() {
		}
		if time.Since(start) < 10*time.Millisecond {
			t.Fatalf("expected to timeout after 10ms")
		}
	})
	t.Run("set max attempts only", func(t *testing.T) {
		backoff := ExponentialBackoff(ExponentialBackoffOptions{
			BaseInterval: time.Millisecond,
			MaxAttempts:  10,
		})
		start := time.Now()
		retryCount := 0
		for backoff.Next() {
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
		backoff := ExponentialBackoff(ExponentialBackoffOptions{
			Context:      ctx,
			BaseInterval: time.Second,
		})
		start := time.Now()
		for backoff.Next() {
		}
		if time.Since(start) < 10*time.Millisecond {
			t.Fatalf("expected to timeout 10ms, but %s elapsed", time.Since(start))
		}
	})
	t.Run("max attempts is prioritized", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		backoff := ExponentialBackoff(ExponentialBackoffOptions{
			Context:      ctx,
			BaseInterval: time.Millisecond,
			MaxAttempts:  10,
		})
		start := time.Now()
		retryCount := 0
		for backoff.Next() {
			retryCount++
		}
		if time.Second < time.Since(start) && retryCount < 10 {
			t.Fatal("expected to reach 10 retries before 1 second elapsed")
		}
	})
	t.Run("timeout is prioritized", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		backoff := ExponentialBackoff(ExponentialBackoffOptions{
			Context:      ctx,
			BaseInterval: time.Millisecond,
			MaxAttempts:  10,
		})
		start := time.Now()
		retryCount := 0
		for backoff.Next() {
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
